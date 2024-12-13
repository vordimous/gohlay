# Gohlay with Redpanda

This example will demonstrate how producers can send delayed messages using the Gohlay header. The first message will be deliver immediately and the second message will deliver after 30 seconds.

## run

You can download and run this example with this script:

```bash
wget -qO- https://github.com/vordimous/gohlay/releases/latest/download/startup.sh | sh -s -- redpanda
```

Or run setup manually:

```bash
./setup.sh
```

Redpanda is setup and the Gohlay service has executed once and did not find any gohlayed messages.

- Two gohlayed messages have been produced onto the `gohlay` topic. One message was scheduled to be sent immediately and one message was set to be delivered in the future. Delivered messages will replace the Gohlay header with a delivery header.
- Check the [Redpanda Console UI](http://localhost:8080/topics/gohlay?p=-1&s=50&o=-1#messages) for your waiting messages. will see the the `now` message has a second message and the `future` message only has one.
- The gohlay service is auto restarting to run new checks for delayed messages.
- Wait for the `future` messages delay time to pass and you will see the second delayed message after refreshing the [Redpanda Console UI](http://localhost:8080/topics/gohlay?p=-1&s=50&o=-1#messages).
- You can Generate a new delayed message and see it get delivered after 30 seconds.

```bash
echo '{"id":200000,"message":"Hello, Future Gohlay"}' | docker compose exec -T kafkacat \
  kafkacat -P \
    -b redpanda:29092 \
    -k "wait1min" \
    -t gohlay \
    -H GOHLAY="$(date -v +1M)"
```

- Teardown the compose stack.

```bash
./teardown.sh
```
