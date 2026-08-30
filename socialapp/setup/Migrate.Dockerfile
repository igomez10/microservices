# The database migration image: Ubuntu, bash, and psql, carrying this app's
# schema. Run as an init container by the Kubernetes deployment in
# microservices-infrastructure.
#
# Built with the APP DIRECTORY as context (not setup/, the way compose builds
# setup/Dockerfile), because it needs db/setup/schema.sql as well as this
# folder. setup/Dockerfile and init.sh are untouched, so compose keeps working.
FROM ubuntu:24.04

# Dependencies are installed here rather than at container start: an init
# container runs on every pod start and every restart, and apt-get at that
# point would add a network dependency to the startup path of the app — a
# failed mirror would stop the service from starting at all.
#
# postgresql-client is Ubuntu's 16.x; the server is 18.x. psql speaks to newer
# servers fine, and nothing here uses version-specific syntax.
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y \
        postgresql-client \
        ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# CREATE DATABASE is stripped at build time rather than at run time: the
# database already exists (CloudNativePG creates it from bootstrap.initdb) and
# CREATE DATABASE cannot run inside the transaction that apply.sql wraps
# everything in. Doing it here also keeps the runtime filesystem read-only.
COPY db/setup/schema.sql /schema.sql.orig
RUN grep -vi '^[[:space:]]*CREATE DATABASE' /schema.sql.orig > /schema.sql \
 && rm /schema.sql.orig

COPY setup/apply.sql /apply.sql
COPY setup/migrate.sh /migrate.sh

# Readable and executable by any uid: Kubernetes runs this with runAsNonRoot
# and a uid that does not exist in /etc/passwd.
RUN chmod 0755 /migrate.sh && chmod 0444 /apply.sql /schema.sql

USER 65532:65532
ENTRYPOINT ["/migrate.sh"]
