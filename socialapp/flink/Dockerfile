# Flink CDC job submitter image.
#
# Built on top of the stock Flink image (which provides FLINK_HOME + the Flink
# client that `flink-cdc.sh` uses to talk to the cluster's REST endpoint) and
# layers the Apache Flink CDC pipeline distribution + the Postgres source and
# Kafka sink pipeline connectors on top.
#
# This container is a one-shot job submitter: it waits for the JobManager, then
# submits the Postgres -> Kafka CDC pipeline and exits.
FROM flink:1.20.1-scala_2.12-java11

# Keep these in sync: the Flink CDC dist and connectors must share a version and
# target the same base Flink minor (1.20).
ARG FLINK_CDC_VERSION=3.6.0
ARG FLINK_CDC_FLAVOR=3.6.0-1.20

ENV FLINK_CDC_HOME=/opt/flink-cdc

RUN set -eux; \
    mkdir -p "${FLINK_CDC_HOME}"; \
    curl -fSL "https://archive.apache.org/dist/flink/flink-cdc-${FLINK_CDC_VERSION}/flink-cdc-${FLINK_CDC_FLAVOR}-bin.tar.gz" \
        -o /tmp/flink-cdc.tgz; \
    tar -xzf /tmp/flink-cdc.tgz -C "${FLINK_CDC_HOME}" --strip-components=1; \
    rm /tmp/flink-cdc.tgz; \
    # Pipeline connectors go into the Flink CDC lib dir so flink-cdc.sh ships
    # them to the cluster as part of the submitted pipeline.
    curl -fSL "https://repo1.maven.org/maven2/org/apache/flink/flink-cdc-pipeline-connector-postgres/${FLINK_CDC_FLAVOR}/flink-cdc-pipeline-connector-postgres-${FLINK_CDC_FLAVOR}.jar" \
        -o "${FLINK_CDC_HOME}/lib/flink-cdc-pipeline-connector-postgres-${FLINK_CDC_FLAVOR}.jar"; \
    curl -fSL "https://repo1.maven.org/maven2/org/apache/flink/flink-cdc-pipeline-connector-kafka/${FLINK_CDC_FLAVOR}/flink-cdc-pipeline-connector-kafka-${FLINK_CDC_FLAVOR}.jar" \
        -o "${FLINK_CDC_HOME}/lib/flink-cdc-pipeline-connector-kafka-${FLINK_CDC_FLAVOR}.jar"

COPY pipeline.yaml "${FLINK_CDC_HOME}/pipeline.yaml"
COPY submit.sh /opt/submit.sh
RUN chmod +x /opt/submit.sh

ENTRYPOINT ["/opt/submit.sh"]
