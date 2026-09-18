# ── Grafana ──────────────────────────────────────────────────────────────────
# socialapp's dashboard and alerts, on the homelab cluster's Grafana
# (var.grafana_url). Everything here belongs to socialapp alone: its folder,
# its dashboard, its alert rules and the Slack contact point they notify.
#
# ROUTING WITHOUT THE POLICY TREE. Grafana keeps one notification policy tree
# per org, and grafana_notification_policy manages that whole tree as a single
# object — two roots declaring it would overwrite each other on every apply.
# So this root does not touch it. Each rule names its contact point directly
# (notification_settings, "simplified routing", Grafana 11+; this instance is
# 13), which lets every app own its alerting end to end without a shared file.
#
# WHAT THE ALERTS CAN SEE. The cluster's OpenTelemetry collector sends metrics
# to `debug` only (infrastructure repo, homelab/apps/otel-collector), so
# socialapp's own HTTP/auth/cache metrics never reach this Prometheus. The
# rules are built on what does:
#   - kube-state-metrics and cAdvisor — availability, restarts, OOM, resources
#   - kubelet volume stats — the Postgres volumes
#   - Loki — the frontend's nginx access log, which carries the status of every
#     /api request (nginx proxies /api/ to the API), and the API's own ERRORs
# HTTP latency panels and SLO-style alerts wait for app metrics to be routed.
#
# Datasource UIDs are the ones kube-prometheus-stack provisions, fixed by the
# chart values in the infrastructure repo's layer 20 — not looked up, because
# a rename there should fail loudly here rather than silently repoint.

locals {
  grafana_prometheus_uid = "prometheus"
  grafana_loki_uid       = "loki"

  # The frontend's nginx access log in combined format, parsed into fields.
  # Every /api request passes through it, so it is the app's request log.
  nginx_api_requests = <<-EOT
    {k8s_namespace_name="socialapp", k8s_container_name="nginx"} | pattern `<_> - <_> [<_>] "<method> <path> <_>" <status> <_> "<_>" "<agent>" <_>` | path=~"/api/.*"
  EOT

  # One entry per alert. `expr` returns the value compared against `threshold`
  # with `op` (gt / lt). no_data_state is OK for alerts whose series only
  # exist while something is wrong (CrashLoop, OOM), and NoData for the
  # availability ones, where a vanished series is itself worth hearing about.
  socialapp_alerts = {
    # ── Availability ─────────────────────────────────────────────────────────
    api_down = {
      title       = "socialapp API down"
      datasource  = "prometheus"
      expr        = "kube_deployment_status_replicas_available{namespace=\"socialapp\", deployment=\"socialapp\"}"
      op          = "lt"
      threshold   = 1
      for         = "2m"
      severity    = "critical"
      no_data     = "NoData"
      summary     = "No socialapp API replica is available."
      description = "Deployment socialapp has had 0 available replicas for 2 minutes. Every /api request is failing."
    }
    frontend_down = {
      title       = "socialapp frontend down"
      datasource  = "prometheus"
      expr        = "kube_deployment_status_replicas_available{namespace=\"socialapp\", deployment=\"socialapp-frontend\"}"
      op          = "lt"
      threshold   = 1
      for         = "2m"
      severity    = "critical"
      no_data     = "NoData"
      summary     = "No socialapp frontend replica is available."
      description = "Deployment socialapp-frontend has had 0 available replicas for 2 minutes. The site, and /api behind it, is unreachable."
    }
    postgres_down = {
      title       = "socialapp Postgres down"
      datasource  = "prometheus"
      expr        = "sum(kube_pod_status_ready{namespace=\"socialapp\", pod=~\"socialapp-postgres-[0-9]+\", condition=\"true\"}) or vector(0)"
      op          = "lt"
      threshold   = 1
      for         = "2m"
      severity    = "critical"
      no_data     = "NoData"
      summary     = "No socialapp Postgres instance is ready."
      description = "The CloudNativePG cluster socialapp-postgres has had no ready instance for 2 minutes."
    }
    redis_down = {
      title       = "socialapp Redis down"
      datasource  = "prometheus"
      expr        = "kube_statefulset_status_replicas_ready{namespace=\"socialapp\", statefulset=\"socialapp-redis\"}"
      op          = "lt"
      threshold   = 1
      for         = "2m"
      severity    = "critical"
      no_data     = "NoData"
      summary     = "socialapp Redis is not ready."
      description = "StatefulSet socialapp-redis has had 0 ready replicas for 2 minutes."
    }

    # ── Degraded ─────────────────────────────────────────────────────────────
    deployment_degraded = {
      title       = "socialapp deployment degraded"
      datasource  = "prometheus"
      expr        = "kube_deployment_status_replicas_available{namespace=\"socialapp\"} / kube_deployment_spec_replicas{namespace=\"socialapp\"}"
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
      title       = "socialapp Postgres degraded"
      datasource  = "prometheus"
      expr        = "sum(kube_pod_status_ready{namespace=\"socialapp\", pod=~\"socialapp-postgres-[0-9]+\", condition=\"true\"}) or vector(0)"
      op          = "lt"
      threshold   = 3
      for         = "10m"
      severity    = "warning"
      no_data     = "NoData"
      summary     = "socialapp Postgres is running with fewer than 3 ready instances."
      description = "The CloudNativePG cluster has had fewer than 3 ready instances for 10 minutes. It still serves, with less redundancy."
    }
    crashlooping = {
      title       = "socialapp container crash-looping"
      datasource  = "prometheus"
      expr        = "max by (pod, container) (kube_pod_container_status_waiting_reason{namespace=\"socialapp\", reason=\"CrashLoopBackOff\"})"
      op          = "gt"
      threshold   = 0
      for         = "5m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} is in CrashLoopBackOff."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} has been crash-looping for 5 minutes."
    }
    # A rollout replaces pods rather than restarting them, so deploys do not
    # count here.
    frequent_restarts = {
      title       = "socialapp container restarting"
      datasource  = "prometheus"
      expr        = "sum by (pod, container) (increase(kube_pod_container_status_restarts_total{namespace=\"socialapp\"}[30m]))"
      op          = "gt"
      threshold   = 3
      for         = "0m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} restarted more than 3 times in 30 minutes."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} keeps restarting. Check its logs and last termination reason."
    }
    oom_killed = {
      title       = "socialapp container OOM-killed"
      datasource  = "prometheus"
      expr        = "max by (pod, container) ((increase(kube_pod_container_status_restarts_total{namespace=\"socialapp\"}[10m]) > 0) * on (pod, container) group_left kube_pod_container_status_last_terminated_reason{namespace=\"socialapp\", reason=\"OOMKilled\"})"
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
      title       = "socialapp memory near limit"
      datasource  = "prometheus"
      expr        = "max by (pod, container) (container_memory_working_set_bytes{namespace=\"socialapp\", container!=\"\"} / on (namespace, pod, container) kube_pod_container_resource_limits{namespace=\"socialapp\", resource=\"memory\"})"
      op          = "gt"
      threshold   = 0.9
      for         = "15m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} is above 90% of its memory limit."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} has used more than 90% of its memory limit for 15 minutes. The next step is an OOM kill."
    }
    cpu_throttled = {
      title       = "socialapp CPU throttled"
      datasource  = "prometheus"
      expr        = "max by (pod, container) (rate(container_cpu_cfs_throttled_periods_total{namespace=\"socialapp\", container!=\"\"}[5m]) / rate(container_cpu_cfs_periods_total{namespace=\"socialapp\", container!=\"\"}[5m]))"
      op          = "gt"
      threshold   = 0.5
      for         = "15m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "{{ $labels.pod }}/{{ $labels.container }} is CPU-throttled over half the time."
      description = "Container {{ $labels.container }} in {{ $labels.pod }} has been throttled in more than 50% of CFS periods for 15 minutes. Latency suffers before anything fails."
    }
    # The Postgres volumes sat at ~59% when this was written.
    volume_filling = {
      title       = "socialapp volume filling up"
      datasource  = "prometheus"
      expr        = "max by (persistentvolumeclaim) (kubelet_volume_stats_used_bytes{namespace=\"socialapp\"} / kubelet_volume_stats_capacity_bytes{namespace=\"socialapp\"})"
      op          = "gt"
      threshold   = 0.85
      for         = "30m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "Volume {{ $labels.persistentvolumeclaim }} is over 85% full."
      description = "PVC {{ $labels.persistentvolumeclaim }} has been more than 85% full for 30 minutes. Postgres stops accepting writes when it fills."
    }

    # ── Deploys ──────────────────────────────────────────────────────────────
    # Both are Argo CD hooks with BeforeHookCreation, so a failed Job stays
    # until the next sync replaces it — the alert holds until then.
    migration_failed = {
      title       = "socialapp migration failed"
      datasource  = "prometheus"
      expr        = "max by (job_name) (kube_job_status_failed{namespace=\"socialapp\", job_name=~\"socialapp-migrate.*\"})"
      op          = "gt"
      threshold   = 0
      for         = "0m"
      severity    = "critical"
      no_data     = "OK"
      summary     = "The socialapp schema migration Job failed."
      description = "Job {{ $labels.job_name }} failed. socialapp-app's sync stopped at the migration hook, so the new version did not roll out."
    }
    smoke_test_failed = {
      title       = "socialapp smoke test failed"
      datasource  = "prometheus"
      expr        = "max by (job_name) (kube_job_status_failed{namespace=\"socialapp\", job_name=~\"socialapp-smoke-test.*\"})"
      op          = "gt"
      threshold   = 0
      for         = "0m"
      severity    = "critical"
      no_data     = "OK"
      summary     = "The socialapp post-deploy smoke test failed."
      description = "Job {{ $labels.job_name }} failed after a sync: the version that just rolled out does not serve."
    }

    # ── Requests and logs (Loki) ─────────────────────────────────────────────
    # 5xx only. 4xx is dominated by clients calling paths that do not exist —
    # when this was written, one client behind the Cloudflare tunnel was
    # sending ~36 POST /api/v1/users a minute, all 404.
    api_5xx_rate = {
      title       = "socialapp API 5xx rate high"
      datasource  = "loki"
      expr        = "(sum(count_over_time(${trimspace(local.nginx_api_requests)} | status=~\"5..\" [10m])) or vector(0)) / sum(count_over_time(${trimspace(local.nginx_api_requests)} [10m]))"
      op          = "gt"
      threshold   = 0.05
      for         = "5m"
      severity    = "critical"
      no_data     = "OK"
      summary     = "More than 5% of /api requests are failing with 5xx."
      description = "Over the last 10 minutes, more than 5% of /api requests through the frontend returned a 5xx (from nginx's access log)."
    }
    # `failed to shutdown` is excluded: every rollout logs it once per replaced
    # pod, and it is not an outage.
    api_error_logs = {
      title       = "socialapp API logging errors"
      datasource  = "loki"
      expr        = "sum(count_over_time({k8s_namespace_name=\"socialapp\", k8s_container_name=\"socialapp\"} | json | level=\"ERROR\" | msg!=\"failed to shutdown\" [10m])) or vector(0)"
      op          = "gt"
      threshold   = 10
      for         = "0m"
      severity    = "warning"
      no_data     = "OK"
      summary     = "The socialapp API logged more than 10 errors in 10 minutes."
      description = "The API logged more than 10 ERROR lines in the last 10 minutes (excluding shutdown noise from rollouts). Check Loki: {k8s_namespace_name=\"socialapp\", k8s_container_name=\"socialapp\"} | json | level=\"ERROR\"."
    }
  }
}

resource "grafana_folder" "socialapp" {
  uid   = "socialapp"
  title = "socialapp"
}

# The webhook comes from its own Vault path, not kv/semaphore/terraform: the
# slack block's `url` is sensitive but not write-only, so whatever this data
# source reads lands in state. Reading kv/semaphore/terraform here would put
# every cloud credential in it. This way state holds one Slack webhook, which
# can post to one channel and nothing else.
#
# `terraform validate` warns that this data source is deprecated in favour of
# the ephemeral vault_kv_secret_v2. Ignore it: an ephemeral value can only feed
# a write-only argument, and the slack block's `url` is not one, so the
# ephemeral form does not work here.
data "vault_kv_secret_v2" "alerting" {
  mount = local.vault_kv_mount_path
  name  = "alerting/socialapp"
}

resource "grafana_contact_point" "socialapp_slack" {
  name = "socialapp-slack"

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

# Small dashboards rather than one large one, linked to each other by the
# `socialapp` tag: overview (is it up, are requests succeeding), resources,
# logs, and service-observability. One per file in dashboards/.
#
# service-observability is the dashboard that used to be provisioned on the OCI
# VM's Grafana (../grafana/provisioning/dashboards/dashboard-socialapp-improved.json),
# repointed at this Prometheus. It is built on socialapp's own OTel metrics,
# which the cluster's collector does not export yet — its panels stay empty
# until that metrics pipeline is routed to Prometheus.
resource "grafana_dashboard" "socialapp" {
  for_each = toset(["overview", "resources", "logs", "service-observability"])

  folder      = grafana_folder.socialapp.uid
  config_json = file("${path.module}/dashboards/${each.key}.json")
  overwrite   = true
}

# The single dashboard this used to be is the overview now.
#
# No dashboard may use the UID `socialapp`: that is the folder's. Grafana's
# folder browser keys its tree by UID, and a dashboard whose UID equals its
# folder's looks like its own parent — the folder page then recurses until the
# browser throws "Maximum call stack size exceeded". Hence socialapp-<name>.
moved {
  from = grafana_dashboard.socialapp
  to   = grafana_dashboard.socialapp["overview"]
}

resource "grafana_rule_group" "socialapp" {
  name             = "socialapp"
  folder_uid       = grafana_folder.socialapp.uid
  interval_seconds = 60

  dynamic "rule" {
    for_each = local.socialapp_alerts

    content {
      name           = rule.value.title
      condition      = "B"
      for            = rule.value.for
      no_data_state  = rule.value.no_data
      exec_err_state = "Error"

      labels = {
        app      = "socialapp"
        severity = rule.value.severity
      }

      annotations = {
        summary       = rule.value.summary
        description   = rule.value.description
        dashboard_uid = jsondecode(grafana_dashboard.socialapp["overview"].config_json)["uid"]
      }

      # Straight to socialapp's contact point — see the ROUTING note above.
      notification_settings {
        contact_point   = grafana_contact_point.socialapp_slack.name
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
        # Two whole jsonencode() calls, not one merge() over a conditional:
        # the branches of `cond ? {queryType = "instant"} : {range = false}`
        # get unified to map(string), which sends Prometheus "range": "false"
        # as a string and fails every evaluation.
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
