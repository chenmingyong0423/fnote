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

if [ "$#" -ne 1 ]; then
  echo "用法: bash script/rebuild-nginx-service.sh <server|admin|web>"
  exit 1
fi

SERVICE="$1"
case "$SERVICE" in
  server|admin|web)
    ;;
  *)
    echo "不支持的服务: $SERVICE"
    echo "可选服务: server admin web"
    exit 1
    ;;
esac

cd "$COMPOSE_DIR"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build --no-deps "$SERVICE"

echo "服务已重新构建并启动: $SERVICE"

