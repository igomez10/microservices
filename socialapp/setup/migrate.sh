#!/usr/bin/env bash
# Applies db/setup/schema.sql to the database named by the connection string.
#
# Runs as an init container in Kubernetes (see the deployment in
# microservices-infrastructure). The compose stack keeps using init.sh, which
# also provisions Elasticsearch and assumes compose's hostnames — neither
# applies here.
set -euo pipefail

# socialapp reads DATABASE_URL, urlshortener reads DB_URL. Accept either, so
# one script serves both services.
DSN="${DATABASE_URL:-${DB_URL:-}}"
if [[ -z "$DSN" ]]; then
    echo "migrate: neither DATABASE_URL nor DB_URL is set" >&2
    exit 1
fi

echo "migrate: waiting for the database to accept connections"
until pg_isready -d "$DSN" >/dev/null 2>&1; do
    sleep 2
done

# Identifies the schema by content, so an edited file is treated as new work
# and an unchanged one is skipped. apply.sql does the deciding, under a lock.
checksum="$(sha256sum /schema.sql | cut -d' ' -f1)"
echo "migrate: schema checksum ${checksum:0:12}"

psql "$DSN" -v ON_ERROR_STOP=1 -v "checksum=$checksum" -q -f /apply.sql

echo "migrate: done"
