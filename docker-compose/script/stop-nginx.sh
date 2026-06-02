#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
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
docker compose "${compose_env_args[@]}" -f "$COMPOSE_FILE" down

echo "Nginx 部署已停止。"
echo "数据目录不会被删除: $COMPOSE_DIR/data/mongo $COMPOSE_DIR/data/logs $COMPOSE_DIR/data/static"
