#!/bin/sh
set -e

export COMPOSE_PROFILES="${COMPOSE_PROFILES:-kafka}"

# Start or restart Gohlay
if [ -z "$(docker compose ps -q gohlay)" ]; then
  docker compose pull
  docker compose up -d
else
  docker compose up -d --force-recreate --no-deps gohlay
fi
