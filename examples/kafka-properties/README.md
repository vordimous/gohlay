# Kafka Properties

This example uses custom kafka authentication options to demonstrate how producers can send delayed messages using the Gohlay header. The first message will be deliver immediately and the second message will deliver after 1 minute.

## run

- start all of the components

```bash
./setup.sh
```

Kafka is setup and the Gohlay service has executed once and did not find any gohlayed messages.

- Two gohlayed messages have been producted onto the `gohlay` topic. One message was scheduled to be sent immediately and one message was set to be delivered in the future. Delivered messages will replace the Gohlay header with a delivery header
- Check the [Kafka UI](http://localhost:8080/ui/clusters/local/all-topics/gohlay/messages) for your waiting messages. will see the the `now` message has a second message and the `future` message only has one.
- The gohlay service is auto restarting to run new checks for delayed messages.
- Wait for the `future` messages delay time to pass and you will see the second delayed message after refreshing the [Kafka UI](http://localhost:8080/ui/clusters/local/all-topics/gohlay/messages).
- You can Generate a new delayed message and see it get delivered after 1 minute

```bash
echo '{"id":200000,"message":"Hello, Future Gohlay"}' | docker compose exec -T kafkacat \
  kafkacat -P \
    -b kafka:29092 \
    -X security.protocol=SASL_PLAINTEXT \
    -X sasl.username=user \
    -X sasl.password=bitnami \
    -X sasl.mechanism=PLAIN \
    -k "wait1min" \
    -t gohlay \
    -H GOHLAY="$(date -v +1M)"
```

- Teardown the compose stack.

```bash
./teardown.sh
```
