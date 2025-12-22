#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

cd "$SCRIPT_DIR"

docker exec -it MindSyncr-db psql -U postgres -d MindSyncr