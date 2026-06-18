#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/operation-timer.sh"
start_operation_timer "重启"

COMPOSE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="${ENV_FILE:-$COMPOSE_DIR/.env.nginx}"
ENV_EXAMPLE_FILE="$COMPOSE_DIR/.env.nginx.example"
COMPOSE_FILE="$COMPOSE_DIR/docker-compose.nginx.yaml"

compose_env_args=()
if [ -f "$ENV_FILE" ]; then
  compose_env_args+=(--env-file "$ENV_FILE")
elif [ -f "$ENV_EXAMPLE_FILE" ]; then
  compose_env_args+=(--env-file "$ENV_EXAMPLE_FILE")
fi

cd "$COMPOSE_DIR"
docker compose "${compose_env_args[@]}" -f "$COMPOSE_FILE" restart

echo "Nginx 部署已重启。"

