# Flink CDC: Postgres → Kafka

Replaces the previous Debezium Kafka Connect setup with [Apache Flink CDC](https://nightlies.apache.org/flink/flink-cdc-docs-stable/).
Change data from the `socialapp` Postgres database is streamed into Kafka using
the Flink CDC **pipeline** framework, which infers table schemas automatically.

## Components (docker-compose)

| Service            | Role                                                                 |
| ------------------ | ------------------------------------------------------------------- |
| `flink-jobmanager` | Flink 1.20 JobManager. Web UI on <http://localhost:8081>.           |
| `flink-taskmanager`| Flink 1.20 TaskManager (worker). Runs the pipeline tasks.           |
| `flink-cdc-job`    | One-shot submitter. Ships the Postgres→Kafka pipeline, then exits.  |

The submitter image (`flink/Dockerfile`) layers the Flink CDC `3.6.0-1.20`
distribution plus the Postgres source and Kafka sink pipeline connectors on top
of the stock Flink image.

## Pipeline

Defined in [`pipeline.yaml`](./pipeline.yaml):

- **Source**: `postgres` — captures every table in `socialapp.public.*` via
  logical replication (slot `flink_cdc_socialapp`, `pgoutput` plugin).
- **Sink**: `kafka` — one topic per table named `<namespace>.<schema>.<table>`
  (e.g. `socialapp.public.users`), encoded as `debezium-json`.

Requires the source database to run with `wal_level=logical` (already set on the
`database` service).

## Usage

```bash
docker compose up -d --build flink-jobmanager flink-taskmanager flink-cdc-job
```

Then:

- Flink UI: <http://localhost:8081> — the running `socialapp-postgres-to-kafka` job.
- Inspect topics in Kafka UI, or:
  ```bash
  docker compose exec broker1 kafka-topics --bootstrap-server broker1:9092 --list
  ```

To re-submit after changing `pipeline.yaml`:

```bash
docker compose up -d --build flink-cdc-job
```
