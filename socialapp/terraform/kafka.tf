# ── Kafka topics ─────────────────────────────────────────────────────────────
# Provisions the Kafka topics used by the socialapp Flink CDC pipeline
# (flink/pipeline.yaml sinks the outbox `events` table to Kafka). Replaces the
# earlier rpk-based topic creation in setup/init.sh so topic layout and config
# are declarative and version-controlled.
#
# Still commented out: the brokers the provider points at are the remote
# docker-compose stack's, not the cluster's, and nothing reaches them from a
# pipeline run yet. Uncommenting is what makes this root need Kafka
# connectivity at plan time — until then the kafka provider is configured but
# never contacted.

# resource "kafka_topic" "socialapp_events" {
#   name               = "socialapp.public.events"
#   partitions         = 3
#   replication_factor = 3

#   config = {
#     # Retain event data for up to 72 hours (72 * 3600 * 1000 ms).
#     "retention.ms" = "259200000"
#   }
# }
