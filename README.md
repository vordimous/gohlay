# Gohlay

The Kafka delayed message delivery tool.

![gopher](gohlay_gopher.png)

Gohlay is a simple CLI tool that can produce messages on Kafka topics after a deadline that is set by a header on the message. This allows the original message producer to communicate a desired execution time for a consumer. Gohlay will republish the message exactly as it was sent removing the deadline header. A consumer will get the updated message and act on it.

Gohlay gives power to both the Producer and Consumer because it doesn't prevent either from acting in a normal manner. The Producer can send any messages without a deadline header and the Consumer can choose to ignore the deadline header.

## Try it out

Run the [Quickstart](./examples/quickstart/) compose example to see Gohlay in action. The following script will download and run the quickstart with the lates Gohlay version or you can copy and run the [compose.yaml](./examples/quickstart/compose.yaml) yourself.

```bash
wget -qO- https://github.com/vordimous/gohlay/releases/latest/download/startup.sh | sh -
```

## Install

Download the binary for your OS from the [latest Gohlay release](https://github.com/vordimous/gohlay/releases/latest).

Run Gohlay in a container:

```bash
docker run --rm ghcr.io/vordimous/gohlay --help
```

```text
Gohlay is a delayed delivery tool for producing messages onto
Kafka topics on a schedule set by a Kafka message header.

Usage:
  gohlay [flags]
  gohlay [command]

Available Commands:
  check       Check for gohlayed messages that are past the deadline
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Check for gohlayed messages and deliver them.

Flags:
  -b, --bootstrap-servers stringArray                          Sets the "bootstrap.servers" property in the kafka.ConfigMap (default [localhost:9092])
      --config-dir string                                      config file directory. (default ".")
  -d, --deadline int                                           Sets the delivery deadline. (Format: Unix Timestamp) (default 1733013640995)
      --debug                                                  Display debugging output in the console.
  -h, --help                                                   help for gohlay
      --json                                                   Display output in the console as JSON.
  -p, --kafka-properties property=value                        Sets the standard librdkafka configuration properties property=value documented in: https://github.com/confluentinc/librdkafka/tree/master/CONFIGURATION.md
  -o, --override-headers default_header_name=new_header_name   Maps the name of default headers to a custom header default_header_name=new_header_name. ex: `GOHLAY=DELAY_TIME`
      --silent                                                 Don't display output in the console.
  -t, --topics stringArray                                     Sets the kafka topics to use (default [gohlay])
  -v, --verbose 
```

## Roadmap

- [X] Config with flags, yaml, and Environment vars
- [X] configurable header names
- [ ] configurable group id
- [X] Kafka SASL auth
- [ ] Deliver to multiple topics
- [ ] Use a configurable offset to start from
- [ ] Deliver all message
- [ ] Native CRON trigger
- [ ] Tested Kafka 2.x and 3.x support
