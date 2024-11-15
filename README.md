# Gohlay

The Kafka delayed message delivery tool.

Gohlay is a simple CLI tool that can produce messages on Kafka topics after a deadline that is set by a header on the message. This allows the origional message producer to communicate a desired execution time for a consumer. Gohlay will republish the message exactly as it was sent removing the deadline header. A consumer will get the updated message and act on it.

Gohlay gives power to both the Producer and Consumer because it doesn't prevent either from acting in a normal manner. The Producer can send any messages without a deadline header and the Consumer can choose to ignore the deadline header.

## Try it out

Use the included [Compose stack](compose.yaml) to see Gohlay in action

- start all of the components

```bash
./setup.sh
```

Kafka is setup and the Gohlay service has executed once and did not find any gohlayed messages.

- Generate a message to be delayed

```bash
echo '{"id":200000,"message":"Hello, Gohlay"}' | docker compose exec -T kafkacat \
  kafkacat -P \
    -b kafka:29092 \
    -k "now" \
    -t gohlay \
    -H GOHLAY="$(date)"
```

```bash
echo '{"id":200000,"message":"Hello, Future Gohlay"}' | docker compose exec -T kafkacat \
  kafkacat -P \
    -b kafka:29092 \
    -k "5min" \
    -t gohlay \
    -H GOHLAY="$(date -v +5M)"
```

Check the [Kafka UI](http://localhost:8080/ui/clusters/local/all-topics/gohlay/messages) for your waiting messages

- restart the gohlay service re-execute the `gohlay run` command

```bash
docker compose restart gohlay
```

- In the [Kafka UI]() you will see the the `now` message has a second message.

- After you wait 5 minutes and restart the `gohlay` service again you will see the second delayed message

- the compose stack

```bash
./teardown.sh
```
