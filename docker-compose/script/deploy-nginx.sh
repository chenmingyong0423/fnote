#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="${ENV_FILE:-$COMPOSE_DIR/.env.nginx}"
COMPOSE_FILE="$COMPOSE_DIR/docker-compose.nginx.yaml"

if [ ! -f "$ENV_FILE" ]; then
  echo "未找到环境变量文件: $ENV_FILE"
  echo "你可以先执行: cp $COMPOSE_DIR/.env.nginx.example $COMPOSE_DIR/.env.nginx"
  exit 1
fi

cd "$COMPOSE_DIR"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build

echo "Nginx 部署已启动。"
echo "如需证书，请放到: $COMPOSE_DIR/nginx/ssl/"

