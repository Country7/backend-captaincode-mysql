#!/bin/sh

set -e

echo "run db migration from start.sh"
# - в контейнере UBUNTU используется /bin/sh, который не поддерживает команду source
# source /app/app.env
. /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"