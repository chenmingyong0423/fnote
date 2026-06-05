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

read_env_file_value() {
  local key="$1"
  local line value

  line="$(grep -E "^[[:space:]]*(export[[:space:]]+)?${key}=" "$ENV_FILE" | tail -n 1 || true)"
  if [ -z "$line" ]; then
    return 0
  fi

  value="${line#*=}"
  value="${value%$'\r'}"
  value="${value%\"}"
  value="${value#\"}"
  value="${value%\'}"
  value="${value#\'}"

  printf '%s' "$value"
}

get_env_value() {
  local key="$1"
  local value="${!key:-}"

  if [ -z "$value" ]; then
    value="$(read_env_file_value "$key")"
  fi

  printf '%s' "$value"
}

require_env() {
  local key="$1"
  local value

  value="$(get_env_value "$key")"
  if [ -z "$value" ]; then
    echo "缺少必填环境变量: $key"
    return 1
  fi
}

missing_env=0
for key in WEBSITE_BASE_HOST WEBSITE_ADMIN_HOST WEBSITE_SERVER_HOST SERVER_NAME SSL_CERT_FILE SSL_KEY_FILE; do
  require_env "$key" || missing_env=1
done

if [ "$missing_env" -ne 0 ]; then
  echo "请补全环境变量文件: $ENV_FILE"
  exit 1
fi

SSL_CERT_FILE_VALUE="$(get_env_value SSL_CERT_FILE)"
SSL_KEY_FILE_VALUE="$(get_env_value SSL_KEY_FILE)"

if [ ! -f "$COMPOSE_DIR/nginx/ssl/$SSL_CERT_FILE_VALUE" ]; then
  echo "未找到 SSL 证书文件: $COMPOSE_DIR/nginx/ssl/$SSL_CERT_FILE_VALUE"
  exit 1
fi

if [ ! -f "$COMPOSE_DIR/nginx/ssl/$SSL_KEY_FILE_VALUE" ]; then
  echo "未找到 SSL 私钥文件: $COMPOSE_DIR/nginx/ssl/$SSL_KEY_FILE_VALUE"
  exit 1
fi

cd "$COMPOSE_DIR"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build

echo "Nginx 部署已启动。"
echo "如需证书，请放到: $COMPOSE_DIR/nginx/ssl/"

