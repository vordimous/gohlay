# Gohlay

The Kafka delayed message delivery tool.

Gohlay is a simple CLI tool that can produce messages on Kafka topics after a deadline that is set by a header on the message. This allows the origional message producer to communicate a desired execution time for a consumer. Gohlay will republish the message exactly as it was sent removing the deadline header. A consumer will get the updated message and act on it.

Gohlay gives power to both the Producer and Consumer because it doesn't prevent either from acting in a normal manner. The Producer can send any messages without a deadline header and the Consumer can choose to ignore the deadline header.

## Try it out

Use the included [Compose stack](compose.yaml) to see Gohlay in action

