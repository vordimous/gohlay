package otelcli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Option applies an option to the Exporter.
type config struct {
	ctx  context.Context
	name string
	cmd  *cobra.Command
}

type Exporter interface {
	enable(config)
}

func Cobra(cmd *cobra.Command, exporters ...Exporter) {
	cfg := config{
		cmd: cmd,
		ctx: cmd.Context(),
		name: cmd.Use,
	}

	for _, ex := range exporters {
		ex.enable(cfg)
	}
}

type LoggingExporter struct {
	otelOptions []otlploghttp.Option
}

// enable implements Exporter.
func (l LoggingExporter) enable(cfg config) {
	// Create the OTLP log exporter that sends logs to configured destination
	if logExporter, err := otlploghttp.New(cfg.ctx, l.otelOptions...); err == nil {

		// Create the logger provider
		lp := log.NewLoggerProvider(
			log.WithProcessor(
				log.NewBatchProcessor(logExporter),
			),
		)

		// Set the logger provider globally
		global.SetLoggerProvider(lp)

		// Instantiate a new slog logger
		logger := otelslog.NewLogger(cfg.name)

		persistentPreRun := cfg.cmd.PersistentPreRun
		cfg.cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			logger.Info(fmt.Sprintf("command (%s %s) called with args %v", cfg.name, cmd.Use, args))
			if persistentPreRun != nil {
				persistentPreRun(cmd, args)
			}
		}
		persistentPostRun := cfg.cmd.PersistentPostRun
		cfg.cmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
			// Ensure the logger is shutdown before exiting so all pending logs are exported
			defer lp.Shutdown(cfg.ctx)
			if persistentPostRun != nil {
				persistentPostRun(cmd, args)
			}
		}
	} else {
		fmt.Errorf("failed to initialize otlploghttp exporter", err)
	}
}

func WithLogging(otelOptions ...otlploghttp.Option) Exporter {
	return LoggingExporter{otelOptions:otelOptions}
}

type TracingExporter struct {
	otelOptions []otlptracehttp.Option
}

// enable implements Exporter.
func (t TracingExporter) enable(cfg config) {

	if exp, err := otlptracehttp.New(cfg.ctx, t.otelOptions...); err == nil {

		// Create a new tracer provider with a batch span processor and the given exporter.
		// Ensure default SDK resources and the required service name are set.
		r, err := resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.name),
			),
		)

		if err != nil {
			panic(err)
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(r),
		)

		otel.SetTracerProvider(tp)

		// Finally, set the tracer that can be used for this package.
		tracer := tp.Tracer(cfg.cmd.Use)

		var span trace.Span
		persistentPreRun := cfg.cmd.PersistentPreRun
		cfg.cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			_, span = tracer.Start(cfg.ctx, fmt.Sprintf("%v %v", cfg.name, cmd.Use))
			if persistentPreRun != nil {
				persistentPreRun(cmd, args)
			}
		}
		persistentPostRun := cfg.cmd.PersistentPostRun
		cfg.cmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
			// Handle shutdown properly so nothing leaks.
			defer func() { _ = tp.Shutdown(cfg.ctx) }()

			span.End()
			if persistentPostRun != nil {
				persistentPostRun(cmd, args)
			}
		}

	} else {
		fmt.Errorf("failed to initialize otlptracehttp exporter", err)
	}
}

func WithTracing(otelOptions ...otlptracehttp.Option) Exporter {
	return TracingExporter{otelOptions:otelOptions}
}

type MetricsExporter struct {
	otelOptions []otlpmetrichttp.Option
}

func (m MetricsExporter) enable(cfg config) {
	if exp, err := otlpmetrichttp.New(cfg.ctx, m.otelOptions...); err == nil {
		// Create a new tracer provider with a batch span processor and the given exporter.
		// Ensure default SDK resources and the required service name are set.
		r, err := resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.name),
			),
		)

		if err != nil {
			panic(err)
		}
		meterProvider := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(r),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp,
				// Default is 1m. Set to 3s for demonstrative purposes.
				sdkmetric.WithInterval(3*time.Second))),
		)

		// Register as global meter provider so that it can be used via otel.Meter
		// and accessed using otel.GetMeterProvider.
		// Most instrumentation libraries use the global meter provider as default.
		// If the global meter provider is not set then a no-op implementation
		// is used, which fails to generate data.
		otel.SetMeterProvider(meterProvider)

		meter := otel.Meter(cfg.name)
		// Finally, set the tracer that can be used for this package.
		cliCounter, err := meter.Int64Counter(
			"cli.counter",
			metric.WithDescription("Number of CLI calls."),
			metric.WithUnit("{call}"),
		)
		if err != nil {
			panic(err)
		}

		persistentPreRun := cfg.cmd.PersistentPreRun
		cfg.cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			cliCounter.Add(cfg.ctx, 1)
			if persistentPreRun != nil {
				persistentPreRun(cmd, args)
			}
		}
		persistentPostRun := cfg.cmd.PersistentPostRun
		cfg.cmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
			// Handle shutdown properly so nothing leaks.
			defer func() { _ = meterProvider.Shutdown(cfg.ctx) }()
			if persistentPostRun != nil {
				persistentPostRun(cmd, args)
			}
		}
	} else {
		fmt.Errorf("failed to initialize otlpmetrichttp exporter", err)
	}
}

func WithMetrics(otelOptions ...otlpmetrichttp.Option) Exporter {
	return MetricsExporter{otelOptions: otelOptions}
}
