# Grafana Dashboards for Service Observability

This directory contains Grafana dashboards that provide comprehensive observability for the microservices in this project.

## Dashboards

### 1. URL Shortener - Service Observability
**File:** `dashboard-urlshortener.json`

Provides complete observability for the URL shortener service, tracking:

#### Traffic Metrics
- **Requests per Second by Endpoint**: Shows the rate of requests for each endpoint with method breakdown
- **Total Traffic**: Overall request rate across all endpoints
- **Requests by Status Code**: Stacked view of all status codes (2xx, 4xx, 5xx)

#### Status & Health Metrics
- **Success Rate (2xx)**: Percentage of successful requests
- **Client Error Rate (4xx)**: Percentage of client error responses
- **Server Error Rate (5xx)**: Percentage of server error responses
- **Client Errors by Endpoint**: Detailed breakdown of 4xx errors per endpoint
- **Server Errors by Endpoint**: Detailed breakdown of 5xx errors per endpoint

#### Latency Metrics
- **Overall Percentiles (p50, p90, p99)**: Response time distribution
- **p99 by Endpoint**: Tail latency for each endpoint

**Data Source:** Prometheus metrics from the `urlshortener` job
- Uses `http_request_duration_*` metrics exposed by the go-http-metrics library
- Metrics are scraped from the `/metrics` endpoint on port 8082

---

### 2. SocialApp - Service Observability
**File:** `dashboard-socialapp-improved.json`

Provides complete observability for the SocialApp service, tracking:

#### Traffic Metrics
- **Requests per Second by Endpoint**: Shows the rate of requests for each pattern/endpoint
- **Total Traffic**: Overall request rate across all endpoints
- **Requests by Status Code**: Stacked view of all status codes

#### Status & Health Metrics
- **Success Rate (2xx)**: Percentage of successful requests
- **Client Error Rate (4xx)**: Percentage of client error responses
- **Server Error Rate (5xx)**: Percentage of server error responses
- **Client Errors by Endpoint**: Detailed breakdown of 4xx errors per endpoint
- **Server Errors by Endpoint**: Detailed breakdown of 5xx errors per endpoint

#### Latency Metrics
- **P50/P90/P99 Latency**: Current latency percentiles
- **Overall Percentiles**: Response time distribution over time
- **p99 by Endpoint**: Tail latency for each endpoint

#### Authentication & Cache Metrics
- **Authentication Duration (p99) by Result**: Time spent authenticating, broken down by auth result (JWT, Basic, failed, etc.)
- **Gandalf Token Cache Operations**: Cache hit/miss rate for JWT tokens

**Data Source:** Prometheus metrics via OpenTelemetry Collector
- Uses `http_server_response_duration` histogram from the Beacon middleware
- Uses `authentication_duration_seconds` histogram from the Gandalf middleware
- Uses `gandalf_token_cache_total` counter for cache metrics
- Metrics are exported via OTLP to the otel-collector, which exposes them to Prometheus

---

## Metrics Architecture

### URL Shortener Service
```
URL Shortener App (Go)
  └─> go-http-metrics library
      └─> Prometheus metrics exposed on :8082/metrics
          └─> OpenTelemetry Collector (Prometheus receiver scrapes urlshortener:8082)
              └─> Prometheus Exporter (exposes on otel-collector:8889)
                  └─> Prometheus scrapes (job: otel-collector)
```

**Key Metrics:**
- `http_request_duration_seconds_count` - Request counter
- `http_request_duration_seconds_sum` - Total request duration
- `http_request_duration_seconds_bucket` - Histogram buckets for percentile calculation

**Labels:**
- `handler` - The endpoint pattern
- `method` - HTTP method (GET, POST, etc.)
- `code` - HTTP status code
- `service` - Service name (urlshortener)

---

### SocialApp Service
```
SocialApp (Go)
  └─> OpenTelemetry SDK
      └─> Beacon Middleware (http_server_response_duration)
      └─> Gandalf Middleware (authentication_duration_seconds, gandalf_token_cache)
      └─> OTLP Exporter
          └─> OpenTelemetry Collector
              └─> Prometheus Remote Write → Mimir
              └─> Prometheus Exporter → Prometheus (scrapes otel-collector:8889)
```

**Key Metrics:**
- `http_server_response_duration_count` - Request counter
- `http_server_response_duration_sum` - Total request duration
- `http_server_response_duration_bucket` - Histogram buckets for percentile calculation
- `authentication_duration_seconds_bucket` - Authentication latency histogram
- `gandalf_token_cache_total` - Cache operation counter

**Labels:**
- `pattern` - The request pattern/route
- `method` - HTTP method
- `status_code` - HTTP status code
- `auth_result` - Authentication result (passed_with_jwt, failed_jwt, etc.)
- `cache` - Cache name
- `status` - Cache operation status (hit, miss)

---

## Dashboard Features

### Color Coding
- **Green**: 2xx success responses
- **Yellow**: 4xx client errors
- **Red**: 5xx server errors

### Percentile Latencies
- **p50 (Median)**: Typical request latency
- **p90**: 90% of requests complete faster than this
- **p99**: Tail latency - important for SLAs

### Time Windows
- All dashboards use `$__rate_interval` for optimal query performance
- Default refresh: 30 seconds
- Default time range: Last 6 hours

### Stat Panels
Both dashboards include stat panels with threshold-based color coding:
- Success rate thresholds: >95% green, >99% yellow
- Error rate thresholds: <1% yellow, <5% red
- Latency thresholds: Based on expected performance

---

## Usage

### Accessing Dashboards
1. Navigate to Grafana (default: http://localhost:3000)
2. Go to Dashboards → Browse
3. Look for:
   - "URL Shortener - Service Observability"
   - "SocialApp - Service Observability"

### Customization
- Dashboards are provisioned automatically from this directory
- To modify: Edit the JSON files and restart Grafana
- Variables can be added in the `templating.list` section

### Alerts
Consider setting up alerts based on:
- Error rate exceeding threshold (e.g., >5% 5xx errors)
- Latency exceeding SLA (e.g., p99 > 1s)
- Traffic drops (potential outage indicator)
- Authentication failures spike

---

## Troubleshooting

### No Data in URL Shortener Dashboard
1. Check that the URL shortener service is running
2. Verify OpenTelemetry Collector is scraping the urlshortener:
   - The urlshortener is scraped by the OTel collector via Prometheus receiver
   - Check OTel collector logs for scraping errors
   - The metrics are then exposed by OTel collector to Prometheus
3. Verify Prometheus is scraping the otel-collector target:
   - Check Prometheus targets: http://localhost:9090/targets
   - Look for `job="otel-collector"` with state "UP"
4. Verify metrics endpoint is accessible:
   - `curl http://urlshortener:8082/metrics`
   - Look for metrics starting with `http_request_duration_seconds`
5. Test the query in Prometheus:
   - Go to http://localhost:9090
   - Run query: `http_request_duration_seconds_count{service="urlshortener"}`
   - You should see metrics with labels like `handler`, `method`, and `code`

### No Data in SocialApp Dashboard
1. Check that the SocialApp service is running
2. Verify OpenTelemetry Collector is running and receiving metrics
3. Check that Prometheus is scraping the otel-collector:
   - Check Prometheus targets: http://localhost:9090/targets
   - Look for `job="otel-collector"` with state "UP"
4. Verify OTLP endpoint is accessible from SocialApp

### Incorrect Latency Values
- Histograms require rate() and histogram_quantile()
- Ensure `$__rate_interval` is appropriate for your traffic volume
- Higher traffic = more accurate percentile calculations

---

## Metrics Best Practices

1. **Use Rate for Counters**: Always use `rate()` with counter metrics
2. **Histogram Buckets**: Review bucket boundaries for your latency patterns
3. **Label Cardinality**: Be careful not to create too many unique label combinations
4. **Aggregation**: Use `sum by()` to aggregate across instances/replicas

---

## Future Enhancements

Potential improvements for these dashboards:

1. **Variables/Filters**
   - Add instance selector to filter by specific service instances
   - Add time range selector
   - Add status code filter

2. **Additional Panels**
   - Request size distribution
   - Response size distribution
   - Active connections
   - Error rate trends

3. **Correlations**
   - Link to Tempo for distributed tracing
   - Link to Loki for log correlation
   - Link to Pyroscope for profiling data

4. **SLO Tracking**
   - Define and track Service Level Objectives
   - SLI (Service Level Indicator) calculations
   - Error budget tracking

---

## Related Documentation

- [Prometheus Configuration](../../prometheus/prometheus.yml)
- [OpenTelemetry Collector Config](../../otel-collector/otel-config.yaml)
- [Beacon Middleware](../../internal/middlewares/beacon/)
- [Gandalf Middleware](../../internal/middlewares/gandalf/)

