#!/bin/sh
set -e

docker compose -p gohlay-kafka-properties down --remove-orphans
