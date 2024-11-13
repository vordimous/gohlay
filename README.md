- start

```bash
./setup.sh
```

- Generate a message to be delayed

```bash
echo '{"id":200000,"username":"fehguy","status":"approved","address":[{"street":"437 Lytton","city":"Palo Alto","state":"CA","zip":"94301"}]}' | docker compose exec -T kafkacat \              
  kafkacat -P \
    -b kafka:29092 \
    -k "c234d09b-2fdf-4538-9d31-27c8e2912d4e" \
    -t myTopic \
    -H GOHLAY="$(date)"
```

- stop

```bash
./teardown.sh
```
