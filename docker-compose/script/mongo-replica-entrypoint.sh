#!/bin/bash
set -e

keyfile=$(mktemp /tmp/mongo-replica-keyfile.XXXXXX)
printf '%s' "${MONGO_REPLICA_SET_KEY:-fnote-local-replica-set-key-change-me}" | base64 | tr -d '\n' > "${keyfile}"
chmod 400 "${keyfile}"
if [ "$(id -u)" = "0" ]; then
  chown mongodb:mongodb "${keyfile}"
fi

exec docker-entrypoint.sh "$@" --replSet "${MONGO_REPLICA_SET_NAME:-rs0}" --keyFile "${keyfile}" --bind_ip_all
