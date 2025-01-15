package otelcli

import (
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

func CobraWithLogging(cmd *cobra.Command) {
	ctx := cmd.Context()
	name := cmd.Use

	// Create the OTLP log exporter that sends logs to configured destination
	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		panic("failed to initialize otlploghttp exporter")
	}

	// Create the logger provider
	lp := log.NewLoggerProvider(
		log.WithProcessor(
			log.NewBatchProcessor(logExporter),
		),
	)

	// Set the logger provider globally
	global.SetLoggerProvider(lp)

	// Instantiate a new slog logger
	logger := otelslog.NewLogger(name)

	persistentPreRun := cmd.PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		logger.Info(fmt.Sprintf("command (%s %s) called with args %v", name, cmd.Use, args))
		if persistentPreRun != nil {
			persistentPreRun(cmd, args)
		}
	}
	persistentPostRun := cmd.PersistentPostRun
	cmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		// Ensure the logger is shutdown before exiting so all pending logs are exported
		defer lp.Shutdown(ctx)
		if persistentPostRun != nil {
			persistentPostRun(cmd, args)
		}
	}
}

func CobraWithTracing(cmd *cobra.Command) {
	ctx := cmd.Context()
	name := cmd.Use

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		panic("failed to initialize otlptracehttp exporter")
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
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
	tracer := tp.Tracer(cmd.Use)

	var span trace.Span
	persistentPreRun := cmd.PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		_, span = tracer.Start(ctx, fmt.Sprintf("%v %v", name, cmd.Use))
		if persistentPreRun != nil {
			persistentPreRun(cmd, args)
		}
	}
	persistentPostRun := cmd.PersistentPostRun
	cmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		// Handle shutdown properly so nothing leaks.
		defer func() { _ = tp.Shutdown(ctx) }()

		span.End()
		if persistentPostRun != nil {
			persistentPostRun(cmd, args)
		}
	}
}

func CobraWithMetrics(cmd *cobra.Command) {
	ctx := cmd.Context()
	name := cmd.Use

	exp, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		panic("failed to initialize otlpmetrichttp exporter")
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
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

	meter := otel.Meter(name)
	// Finally, set the tracer that can be used for this package.
	cliCounter, err := meter.Int64Counter(
		"cli.counter",
		metric.WithDescription("Number of CLI calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

	persistentPreRun := cmd.PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cliCounter.Add(ctx, 1)
		if persistentPreRun != nil {
			persistentPreRun(cmd, args)
		}
	}
	persistentPostRun := cmd.PersistentPostRun
	cmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		// Handle shutdown properly so nothing leaks.
		defer func() { _ = meterProvider.Shutdown(ctx) }()
		if persistentPostRun != nil {
			persistentPostRun(cmd, args)
		}
	}
}
