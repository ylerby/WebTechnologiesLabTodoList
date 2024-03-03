#!/bin/bash
set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASSWORD pg_isready -h "$host" -U "postgres"; do
  >&2 echo "Postgres ожидает запуска"
  sleep 1
done

>&2 echo "Postgres запущен"
exec $cmd