# Quickstart

This quickstart will demonstrate how producers can send delayed messages using the Gohlay header. The first message will be deliver immediately and the second message will deliver after 1 minute

## run

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
    -k "1min" \
    -t gohlay \
    -H GOHLAY="$(date -v +1M)"
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
