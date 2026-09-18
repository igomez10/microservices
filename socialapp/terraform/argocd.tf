# ── Argo CD ──────────────────────────────────────────────────────────────────
# Argo CD's view of socialapp: one AppProject and two Applications — the
# namespace and datastores (`socialapp-data`) and the stateless API, frontend,
# migration and smoke test (`socialapp-app`), rendered from the Helm chart in
# this repo.
#
# These are API objects that describe how socialapp deploys, so they live with
# socialapp — and so does everything they render. Both Applications read only
# this repo: `socialapp-data` from ../deploy/manifests, `socialapp-app` from
# ../helm/socialapp with its values.homelab.yaml. Nothing here reaches into the
# infrastructure repo's Git; the one link left to it is layer 30's state
# (remote-state.tf).
#
# The argocd provider is configured in providers.tf.

# Scoped narrowly. Infrastructure repo doc 06: "Scope projects to only their
# repositories, namespaces, and required resource kinds."
#
# The app renders one namespaced CloudNativePG Cluster and nothing
# cluster-scoped but its own namespace — the CRDs it depends on are installed by
# the operator's Helm release, deliberately outside this project's reach.
resource "argocd_project" "socialapp" {
  metadata {
    name      = "homelab-socialapp"
    namespace = "argocd"
  }

  spec {
    # At most 255 bytes: Kubernetes rejects a longer AppProject description.
    description = "socialapp's own resources: socialapp-data (namespace, PostgreSQL, Redis) and socialapp-app (API, frontend, migration, smoke test), both rendered from github.com/igomez10/microservices."

    # This repo only. Both Applications render from it, and listing nothing
    # else means neither can be pointed at another repo without widening the
    # project first.
    #
    # Adding a repo here needs care: Argo CD enforces `repositories, update` on
    # every repo a project update adds or removes, and Terraform's account
    # (role:terraform, layers/20-cluster-platform/modules/argocd/values.yaml in
    # the infrastructure repo) holds it for https://github.com/igomez10/* only.
    # An SSH URL here fails the apply with a permission error — on purpose.
    source_repos = [
      var.application_repo_url,
    ]

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "socialapp"
    }

    # The one cluster-scoped kind this project may touch: its own namespace,
    # which lives with the app (manifests/namespace.yaml). Everything else
    # cluster-scoped — CRDs, RBAC, other namespaces — stays denied.
    cluster_resource_whitelist {
      group = ""
      kind  = "Namespace"
    }

    namespace_resource_whitelist {
      group = "*"
      kind  = "*"
    }
  }
}

# Tracks ../deploy/manifests in this repo, on mainline, as a plain Kustomize
# directory: the namespace, the CloudNativePG Cluster and Redis.
#
# A separate Application from `socialapp_app`, not a second source on it, and
# the reason is the sync policy below. `socialapp_app` runs automated with
# prune and self_heal; folding the database into it would put the Cluster and
# the Redis StatefulSet under automatic pruning and re-sync them on every
# image-digest commit. Two Applications is what keeps "the data" and "the
# release" on different triggers.
#
# Sync is MANUAL: this is a database, and an automated sync of a bad render with
# no one at the keyboard is worse than a stale OutOfSync badge.
#
# Requires, in this order, before the first sync:
#   1. the CloudNativePG operator (infrastructure layer 20, helm_release.cloudnative_pg)
#   2. the socialapp-postgres-app Secret — see the infrastructure repo's
#      homelab/apps/socialapp/README.md
resource "argocd_application" "socialapp_data" {
  # Explicit, so removing this resource from Terraform deletes the Application
  # object and nothing else. A cascading delete would take every resource the
  # Application manages with it — including the database.
  cascade = false

  # Argo CD writes argocd.argoproj.io/hydrate onto every Application itself.
  # It is Argo's bookkeeping, not configuration: without this, every plan
  # would propose stripping it, and every apply would, until Argo re-adds it.
  lifecycle {
    ignore_changes = [metadata[0].annotations["argocd.argoproj.io/hydrate"]]
  }

  metadata {
    name      = "socialapp-data"
    namespace = "argocd"
  }

  spec {
    project = argocd_project.socialapp.metadata[0].name

    source {
      repo_url        = var.application_repo_url
      target_revision = "mainline"
      path            = "socialapp/deploy/manifests"
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "socialapp"
    }

    # No `automated` block on purpose — see above.
    sync_policy {
      sync_options = [
        # The namespace is manifests/namespace.yaml, applied before the rest of
        # the render. A namespace Argo conjures itself carries no Pod Security
        # labels.
        "CreateNamespace=false",
        # A client-side apply merges map fields instead of replacing them, and
        # the operator writes to the Cluster resource it is given, so Argo and
        # the operator both touch it.
        "ServerSideApply=true",
      ]
    }

    # A StatefulSet's volumeClaimTemplates come back from the API server with
    # fields no manifest can contain: a defaulted volumeMode, the PVC's own
    # apiVersion/kind, and a status block. Argo's status-stripping does not
    # reach them because they are nested under spec, so socialapp-redis sits
    # permanently OutOfSync without this. Scoped to those four fields — ignoring
    # the list wholesale would also hide a real change to storage size or class.
    ignore_difference {
      group = "apps"
      kind  = "StatefulSet"
      jq_path_expressions = [
        ".spec.volumeClaimTemplates[].status",
        ".spec.volumeClaimTemplates[].spec.volumeMode",
        ".spec.volumeClaimTemplates[].apiVersion",
        ".spec.volumeClaimTemplates[].kind",
      ]
    }

    # CloudNativePG defaults unset fields on the Cluster it manages and writes
    # back things like the image it selected. That is the operator's business,
    # not drift from Git.
    ignore_difference {
      group                   = "postgresql.cnpg.io"
      kind                    = "Cluster"
      managed_fields_managers = ["cloudnative-pg"]
    }
  }
}

# The API, frontend, schema migration Job and post-sync smoke test, rendered
# from the Helm chart that lives with the code, in this repo.
#
# One source: the chart, rendered with values.homelab.yaml from the chart
# directory itself. That file holds everything cluster-specific, including the
# image digests, so a change to a flag, to how it is templated, and to which
# image runs are all commits to this repo.
#
# Sync is AUTOMATED, unlike `socialapp-data`. A push to mainline builds images,
# the image workflow commits their digests to values.homelab.yaml on mainline,
# and Argo CD deploys that commit with no one clicking Sync. Chart changes on
# mainline deploy the same way.
#
# That digest commit lands on the branch whose pushes trigger the image
# workflow. It does not loop, for two independent reasons: socialapp/helm/** is
# excluded from that workflow's paths, and pushes made with the default
# GITHUB_TOKEN never start workflow runs.
#
# The cost is that every one of those syncs runs the migration Job against the
# production database unattended. That is acceptable because the migration is
# a checksum-guarded no-op unless db/setup/schema.sql changed, and the PostSync
# smoke test fails the sync when the result does not serve.
resource "argocd_application" "socialapp_app" {
  # Explicit, so removing this resource from Terraform deletes the Application
  # object and nothing else.
  cascade = false

  # Argo CD writes argocd.argoproj.io/hydrate onto every Application itself.
  # It is Argo's bookkeeping, not configuration: without this, every plan
  # would propose stripping it, and every apply would, until Argo re-adds it.
  lifecycle {
    ignore_changes = [metadata[0].annotations["argocd.argoproj.io/hydrate"]]
  }

  metadata {
    name      = "socialapp-app"
    namespace = "argocd"
  }

  spec {
    project = argocd_project.socialapp.metadata[0].name

    source {
      repo_url        = var.application_repo_url
      target_revision = "mainline"
      path            = "socialapp/helm/socialapp"

      helm {
        # The chart's object names are fixed (socialapp, socialapp-frontend,
        # ...) rather than derived from this, so it is cosmetic — but it is
        # the Helm release name the templates see.
        release_name = "socialapp"
        # Relative to the chart directory, and applied over its values.yaml.
        value_files = ["values.homelab.yaml"]

        # The public hostnames, built from var.domain so the domain never
        # appears in this public repo. Inline values override value_files, and
        # land only in the Application object and this root's (private) state.
        # Changing one is a `terraform apply`, not a commit — they change about
        # never, which is what makes that trade acceptable.
        values = yamlencode({
          api = {
            config = {
              SOCIALAPP_SUBDOMAIN    = "socialapp.${var.domain}"
              LOCAL_SUBDOMAIN        = "socialapp.${var.domain}"
              PROPERTIES_SUBDOMAIN   = "properties.${var.domain}"
              URLSHORTENER_SUBDOMAIN = "urlshortener.${var.domain}"
              KIBANA_SUBDOMAIN       = "kibana.${var.domain}"
              PUTTYKNIFE_DOMAIN      = "puttyknife-server.${var.domain}"
            }
          }
        })
      }
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "socialapp"
    }

    sync_policy {
      # `prune` because removing something from the chart should actually
      # remove it; nothing here owns a PersistentVolume. `self_heal` so a
      # hand edit to a Deployment is reverted to what Git says.
      automated {
        prune     = true
        self_heal = true
      }

      sync_options = [
        # The namespace belongs to the `socialapp-data` Application.
        "CreateNamespace=false",
        "ServerSideApply=true",
      ]
    }
  }
}
