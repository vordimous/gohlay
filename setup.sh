#!/bin/sh
set -e

# Start or restart Golay
if [ -z "$(docker compose ps -q golay)" ]; then
    docker compose up -d
else
    docker compose up -d --force-recreate --no-deps golay
fi
