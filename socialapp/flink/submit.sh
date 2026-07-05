#!/usr/bin/env bash
# Waits for the Flink JobManager, points the Flink client at it, then submits
# the Apache Flink CDC Postgres -> Kafka pipeline. One-shot: exits once the job
# is accepted by the cluster (the job itself keeps running on the cluster).
set -euo pipefail

JOBMANAGER_HOST="${JOBMANAGER_HOST:-flink-jobmanager}"
JOBMANAGER_REST_PORT="${JOBMANAGER_REST_PORT:-8081}"
REST="http://${JOBMANAGER_HOST}:${JOBMANAGER_REST_PORT}"

# This container is only a Flink client/submitter: point its REST address at the
# remote JobManager. The stock image ships rest.address as 0.0.0.0, so we rewrite
# it explicitly. Handles standardized config.yaml (Flink 1.19+) and legacy
# flink-conf.yaml.
CFG_YAML="${FLINK_HOME}/conf/config.yaml"
CFG_LEGACY="${FLINK_HOME}/conf/flink-conf.yaml"
if [ -f "$CFG_YAML" ]; then
  # Nested YAML: the rest.address line is the only 2-space-indented `address:`
  # (jobmanager.rpc.address is nested deeper at 4 spaces).
  sed -i "s|^  address: .*|  address: ${JOBMANAGER_HOST}|" "$CFG_YAML"
fi
if [ -f "$CFG_LEGACY" ]; then
  if grep -q "^rest.address:" "$CFG_LEGACY"; then
    sed -i "s|^rest.address:.*|rest.address: ${JOBMANAGER_HOST}|" "$CFG_LEGACY"
  else
    echo "rest.address: ${JOBMANAGER_HOST}" >> "$CFG_LEGACY"
  fi
fi

echo "[flink-cdc-job] Waiting for JobManager REST at ${REST} ..."
until curl -sf "${REST}/overview" >/dev/null 2>&1; do
  echo "[flink-cdc-job]   ... not ready yet, retrying in 5s"
  sleep 5
done
echo "[flink-cdc-job] JobManager is up."

echo "[flink-cdc-job] Submitting Postgres -> Kafka CDC pipeline"
"${FLINK_CDC_HOME}/bin/flink-cdc.sh" \
  --flink-home "${FLINK_HOME}" \
  "${FLINK_CDC_HOME}/pipeline.yaml"

echo "[flink-cdc-job] Pipeline submitted. It now runs on the Flink cluster; see the Flink UI at ${REST}."
