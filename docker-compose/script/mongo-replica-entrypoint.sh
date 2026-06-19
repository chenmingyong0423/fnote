#!/bin/bash
set -e

keyfile=/tmp/mongo-replica-keyfile
printf '%s' "${MONGO_REPLICA_SET_KEY:-fnote-local-replica-set-key-change-me}" | base64 | tr -d '\n' > "${keyfile}"
chmod 400 "${keyfile}"
chown mongodb:mongodb "${keyfile}"

exec docker-entrypoint.sh "$@" --replSet "${MONGO_REPLICA_SET_NAME:-rs0}" --keyFile "${keyfile}" --bind_ip_all
