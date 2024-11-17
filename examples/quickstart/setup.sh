#!/bin/sh
set -e

# Start or restart Gohlay
if [ -z "$(docker compose ps -q gohlay)" ]; then
  docker pull ghcr.io/vordimous/gohlay
  docker compose up -d
else
  docker compose up -d --force-recreate --no-deps gohlay
fi
