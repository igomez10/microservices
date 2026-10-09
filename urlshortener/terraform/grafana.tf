# ── Grafana ──────────────────────────────────────────────────────────────────
# urlshortener's dashboards and alerts, on the homelab cluster's Grafana
# (var.grafana_url). Same shape as ../../socialapp/terraform/grafana.tf, which
# has the longer explanation of the choices below; this file only notes where
# urlshortener differs.
#
# ROUTING WITHOUT THE POLICY TREE. Every rule names its contact point directly
# (notification_settings), so this root never touches the org's shared
# notification policy tree — socialapp's root does the same, and two roots
# managing that tree would overwrite each other on every apply.
#
# WHAT THE ALERTS CAN SEE. urlshortener sends traces over OTLP but no metrics,
# and its Prometheus /metrics endpoint (the meta port) is not scraped, so there
# are no request metrics in Prometheus. The rules use:
#   - kube-state-metrics and cAdvisor — availability, restarts, OOM, resources
#   - kubelet volume stats — the Postgres volumes
#   - Loki — urlshortener's own "finished request" log line, which carries the
#     status of every request (no proxy sits in front of it, unlike socialapp's
#     nginx), and its ERROR lines
#
# Datasource UIDs are the ones kube-prometheus-stack provisions, fixed by the
# chart values in the infrastructure repo's layer 20.

locals {
  grafana_prometheus_uid = "prometheus"
  grafana_loki_uid       = "loki"

  # urlshortener's request log: one line per request, with status_code,
  # method and x-pattern (the chi route, so aliases do not become labels).
  # The line filter first, so Loki only JSON-parses the lines it needs.
  urlshortener_requests = <<-EOT
    {k8s_namespace_name="urlshortener", k8s_container_name="urlshortener"} |= "finished request" | json | msg="finished request"
  EOT

  # One entry per alert; see socialapp's grafana.tf for the fields.
  urlshortener_alerts = {
    # ── Availability ─────────────────────────────────────────────────────────
    api_down = {
      title       = "urlshortener API down"
      datasource  = "prometheus"
      expr        = "kube_deployment_status_replicas_available{namespace=\"urlshortener\", deployment=\"urlshortener\"}"
      op          = "lt"
      threshold   = 1
      for         = "2m"
      severity    = "critical"
      no_data     = "NoData"
      summary     = "No urlshortener replica is available."
      description = "Deployment urlshortener has had 0 available replicas for 2 minutes. Every request fails, including socialapp's /api/v1/urls."
    }
    postgres_down = {
      title       = "urlshortener Postgres down"
      datasource  = "prometheus"
      expr        = "sum(kube_pod_status_ready{namespace=\"urlshortener\", pod=~\"urlshortener-postgres-[0-9]+\", condition=\"true\"}) or vector(0)"
      op          = "lt"
      threshold   = 1
      for         = "2m"
      severity    = "critical"
      no_data     = "NoData"
      summary     = "No urlshortener Postgres instance is ready."
      description = "The CloudNativePG cluster urlshortener-postgres has had no ready instance for 2 minutes."
    }

    # ── Degraded ─────────────────────────────────────────────────────────────
    deployment_degraded = {
      title       = "urlshortener deployment degraded"
      datasource  = "prometheus"
      expr        = "kube_deployment_status_replicas_available{namespace=\"urlshortener\"} / kube_deployment_spec_replicas{namespace=\"urlshortener\"}"
      op          = "lt"
      threshold   = 1
      for         = "10m"
      severity    = "warning"
      no_data     = "NoData"
      summary     = "{{ $labels.deployment }} is running below its desired replicas."
      description = "Deployment {{ $labels.deployment }} has had fewer available replicas than desired for 10 minutes. A rollout is stuck or pods cannot schedule."
    }
    # 3 is the instance count in ../deploy/manifests/postgres.yaml.
    postgres_degraded = {
      title       = "urlshortener Postgres degraded"
      datasource  = "prometheus"
      expr        = "sum(kube_pod_status_ready{namespace=\"urlshortener\", pod=~\"urlshortener-postgres-[0-9]+\", condition=\"true\"}) or vector(0)"
      op          = "lt"
      threshold   = 3
      for         = "10m"
      severity    = "warning"
      no_data     = "NoData"
      summary     = "urlshortener Postgres is running with fewer than 3 ready instances."
      description = "The CloudNativePG cluster has had fewer than 3 ready instances for 10 minutes. It still serves, with less redundancy."
    }
    crashlooping = {
      title       = "urlshortener container crash-looping"
      datasource  = "prometheus"
      expr        = "max by (pod, container) (kube_pod_container_status_waiting_reason{namespace=\"urlshortener\", reason=\"CrashLoopBackOff\"})"
      op          = "gt"
      threshold   = 0
      for         = "5m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} is in CrashLoopBackOff."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} has been crash-looping for 5 minutes."
    }
    frequent_restarts = {
      title       = "urlshortener container restarting"
      datasource  = "prometheus"
      expr        = "sum by (pod, container) (increase(kube_pod_container_status_restarts_total{namespace=\"urlshortener\"}[30m]))"
      op          = "gt"
      threshold   = 3
      for         = "0m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} restarted more than 3 times in 30 minutes."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} keeps restarting. Check its logs and last termination reason."
    }
    oom_killed = {
      title       = "urlshortener container OOM-killed"
      datasource  = "prometheus"
      expr        = "max by (pod, container) ((increase(kube_pod_container_status_restarts_total{namespace=\"urlshortener\"}[10m]) > 0) * on (pod, container) group_left kube_pod_container_status_last_terminated_reason{namespace=\"urlshortener\", reason=\"OOMKilled\"})"
      op          = "gt"
      threshold   = 0
      for         = "0m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} was OOM-killed."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} restarted after running out of memory in the last 10 minutes. Its limit is too low, or it is leaking."
    }

    # ── Resources ────────────────────────────────────────────────────────────
    memory_near_limit = {
      title       = "urlshortener memory near limit"
      datasource  = "prometheus"
      expr        = "max by (pod, container) (container_memory_working_set_bytes{namespace=\"urlshortener\", container!=\"\"} / on (namespace, pod, container) kube_pod_container_resource_limits{namespace=\"urlshortener\", resource=\"memory\"})"
      op          = "gt"
      threshold   = 0.9
      for         = "15m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} is above 90% of its memory limit."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} has used more than 90% of its memory limit for 15 minutes. The next step is an OOM kill."
    }
    cpu_throttled = {
      title       = "urlshortener CPU throttled"
      datasource  = "prometheus"
      expr        = "max by (pod, container) (rate(container_cpu_cfs_throttled_periods_total{namespace=\"urlshortener\", container!=\"\"}[5m]) / rate(container_cpu_cfs_periods_total{namespace=\"urlshortener\", container!=\"\"}[5m]))"
      op          = "gt"
      threshold   = 0.5
      for         = "15m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} is CPU-throttled over half the time."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} has been throttled in more than 50% of CFS periods for 15 minutes. Latency suffers before anything fails."
    }
    volume_filling = {
      title       = "urlshortener volume filling up"
      datasource  = "prometheus"
      expr        = "max by (persistentvolumeclaim) (kubelet_volume_stats_used_bytes{namespace=\"urlshortener\"} / kubelet_volume_stats_capacity_bytes{namespace=\"urlshortener\"})"
      op          = "gt"
      threshold   = 0.85
      for         = "30m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "Volume {{ $labels.persistentvolumeclaim }} is over 85% full."
      description = "PVC {{ $labels.persistentvolumeclaim }} has been more than 85% full for 30 minutes. Postgres stops accepting writes when it fills."
    }

    # ── Deploys ──────────────────────────────────────────────────────────────
    # On the Job's terminal Failed condition, not kube_job_status_failed: that
    # counts failed pod attempts, and the migration Job retries
    # (backoffLimit: 3), so one failed attempt followed by a successful retry
    # would leave it nonzero and this alert firing on a sync that went fine.
    # kube_job_failed{condition="true"} exists only once the Job has given up.
    migration_failed = {
      title       = "urlshortener migration failed"
      datasource  = "prometheus"
      expr        = "max by (job_name) (kube_job_failed{namespace=\"urlshortener\", job_name=~\"urlshortener-migrate.*\", condition=\"true\"})"
      op          = "gt"
      threshold   = 0
      for         = "0m"
      severity    = "critical"
      no_data     = "OK"
      summary     = "The urlshortener schema migration Job failed."
      description = "Job {{ $labels.job_name }} failed. urlshortener-app's sync stopped at the migration hook, so the new version did not roll out."
    }
    smoke_test_failed = {
      title       = "urlshortener smoke test failed"
      datasource  = "prometheus"
      expr        = "max by (job_name) (kube_job_failed{namespace=\"urlshortener\", job_name=~\"urlshortener-smoke-test.*\", condition=\"true\"})"
      op          = "gt"
      threshold   = 0
      for         = "0m"
      severity    = "critical"
      no_data     = "OK"
      summary     = "The urlshortener post-deploy smoke test failed."
      description = "Job {{ $labels.job_name }} failed after a sync: the version that just rolled out does not serve."
    }

    # ── Requests and logs (Loki) ─────────────────────────────────────────────
    # 5xx only. 4xx is normal traffic here: 404 when the integration test
    # reads an alias after deleting it, 409 when an alias already exists.
    api_5xx_rate = {
      title       = "urlshortener API 5xx rate high"
      datasource  = "loki"
      expr        = "(sum(count_over_time(${trimspace(local.urlshortener_requests)} | status_code=~\"5..\" [10m])) or vector(0)) / sum(count_over_time(${trimspace(local.urlshortener_requests)} [10m]))"
      op          = "gt"
      threshold   = 0.05
      for         = "5m"
      severity    = "critical"
      no_data     = "OK"
      summary     = "More than 5% of urlshortener requests are failing with 5xx."
      description = "Over the last 10 minutes, more than 5% of urlshortener requests returned a 5xx (from its own request log)."
    }
    api_error_logs = {
      title       = "urlshortener API logging errors"
      datasource  = "loki"
      expr        = "sum(count_over_time({k8s_namespace_name=\"urlshortener\", k8s_container_name=\"urlshortener\"} | json | level=\"ERROR\" [10m])) or vector(0)"
      op          = "gt"
      threshold   = 10
      for         = "0m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "The urlshortener API logged more than 10 errors in 10 minutes."
      description = "The API logged more than 10 ERROR lines in the last 10 minutes. Check Loki: {k8s_namespace_name=\"urlshortener\", k8s_container_name=\"urlshortener\"} | json | level=\"ERROR\"."
    }
  }
}

resource "grafana_folder" "urlshortener" {
  uid   = "urlshortener"
  title = "urlshortener"
}

# Its own Vault path and its own contact point, as socialapp has: the slack
# block's `url` is not write-only, so whatever this reads lands in state, and
# kv/alerting/urlshortener holds one webhook and nothing else. See socialapp's
# grafana.tf for why this is a data source rather than ephemeral.
data "vault_kv_secret_v2" "alerting" {
  mount = local.vault_kv_mount_path
  name  = "alerting/urlshortener"
}

resource "grafana_contact_point" "urlshortener_slack" {
  name = "urlshortener-slack"

  slack {
    # trimspace: a value stored with `vault kv put key=@file` keeps the file's
    # trailing newline, and Grafana rejects the URL with it.
    url   = trimspace(data.vault_kv_secret_v2.alerting.data["slack_webhook_url"])
    title = "{{ .Status | toUpper }}: {{ .CommonLabels.alertname }}"
    text  = <<-EOT
      {{ range .Alerts }}*{{ .Labels.severity | toUpper }}* — {{ .Annotations.summary }}
      {{ .Annotations.description }}
      {{ end }}
    EOT
  }
}

# overview, resources and logs, linked to each other by the `urlshortener` tag.
# No service-observability dashboard, unlike socialapp: it is built on the
# app's own HTTP metrics, which urlshortener does not export to Prometheus.
#
# No dashboard may use the UID `urlshortener`: that is the folder's, and a
# dashboard whose UID equals its folder's breaks Grafana's folder browser (see
# socialapp's grafana.tf). Hence urlshortener-<name>.
resource "grafana_dashboard" "urlshortener" {
  for_each = toset(["overview", "resources", "logs"])

  folder      = grafana_folder.urlshortener.uid
  config_json = file("${path.module}/dashboards/${each.key}.json")
  overwrite   = true
}

resource "grafana_rule_group" "urlshortener" {
  name             = "urlshortener"
  folder_uid       = grafana_folder.urlshortener.uid
  interval_seconds = 60

  dynamic "rule" {
    for_each = local.urlshortener_alerts

    content {
      name           = rule.value.title
      condition      = "B"
      for            = rule.value.for
      no_data_state  = rule.value.no_data
      exec_err_state = "Error"

      labels = {
        app      = "urlshortener"
        severity = rule.value.severity
      }

      annotations = {
        summary       = rule.value.summary
        description   = rule.value.description
        dashboard_uid = jsondecode(grafana_dashboard.urlshortener["overview"].config_json)["uid"]
      }

      notification_settings {
        contact_point   = grafana_contact_point.urlshortener_slack.name
        group_by        = ["alertname"]
        repeat_interval = "4h"
      }

      # A: the query, instant, one value per series.
      data {
        ref_id         = "A"
        datasource_uid = rule.value.datasource == "loki" ? local.grafana_loki_uid : local.grafana_prometheus_uid
        # Grafana derives this from the Loki model's queryType and stores it;
        # declaring it keeps the plan from trying to clear it every time.
        query_type = rule.value.datasource == "loki" ? "instant" : null
        relative_time_range {
          from = 600
          to   = 0
        }
        # Two whole jsonencode() calls, not one merge() over a conditional —
        # see socialapp's grafana.tf: a unified map(string) sends Prometheus
        # "range": "false" as a string and fails every evaluation.
        model = rule.value.datasource == "loki" ? jsonencode({
          refId     = "A"
          expr      = rule.value.expr
          instant   = true
          queryType = "instant"
          }) : jsonencode({
          refId   = "A"
          expr    = rule.value.expr
          instant = true
          range   = false
        })
      }

      # B: fires per series where A crosses the threshold.
      data {
        ref_id         = "B"
        datasource_uid = "__expr__"
        relative_time_range {
          from = 0
          to   = 0
        }
        model = jsonencode({
          refId      = "B"
          type       = "threshold"
          expression = "A"
          conditions = [{
            evaluator = {
              type   = rule.value.op
              params = [rule.value.threshold]
            }
          }]
        })
      }
    }
  }
}
