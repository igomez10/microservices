# Docker Compose Logging Standards

## Overview

All services in the socialapp docker-compose.yml use Grafana Loki as the logging driver with standardized configuration options.

## Standard Configuration

Every service should use the following logging configuration:

```yaml
logging:
  driver: loki
  options:
    loki-url: "http://localhost:3100/loki/api/v1/push"
    loki-batch-size: "5000"
    loki-retries: "3"
    loki-max-backoff: "1000ms"
    loki-timeout: "1000ms"
    loki-batch-wait: "5s"
    labels: "service_name"
    max-size: "10m"
    max-file: "3"
```

## Configuration Options Explained

| Option | Value | Description |
|--------|-------|-------------|
| `loki-url` | `http://localhost:3100/loki/api/v1/push` | Loki push API endpoint |
| `loki-batch-size` | `5000` | Maximum batch size (in bytes) before sending to Loki |
| `loki-retries` | `3` | Number of retry attempts for failed log pushes |
| `loki-max-backoff` | `1000ms` | Maximum backoff time between retries |
| `loki-timeout` | `1000ms` | Timeout for sending logs to Loki |
| `loki-batch-wait` | `5s` | Maximum time to wait before sending a batch |
| `labels` | `service_name` | Labels to attach to log entries |
| `max-size` | `10m` | Maximum size of log file before rotation |
| `max-file` | `3` | Number of rotated log files to keep |

## Notes

- The `loki` service itself does not have a logging configuration (it is the log aggregator)

