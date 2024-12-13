# High Load

This example cam demonstrate how a large number of messages effects the run time and accuracy of deliver for Gohlay.

## run

- start all of the components

```bash
./setup.sh
```

Kafka is setup and the Gohlay service has executed once and did not find any gohlayed messages.

- 1k gohlayed messages have been produced onto the `gohlay` topic. Each message was scheduled to be sent at different times within the next 100 seconds
- The gohlay service is auto restarting to run new checks for delayed messages.
- You can watch messages be delivered on the kafka topic using the `kafkacat` consumer below.

```bash
docker compose exec kafkacat \
  kafkacat -C -b kafka:29092 -t gohlay -J -u | jq .
```

- Teardown the compose stack.

```bash
./teardown.sh
```
