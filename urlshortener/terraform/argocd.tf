# ── Argo CD ──────────────────────────────────────────────────────────────────
# Argo CD's view of urlshortener: one AppProject and two Applications — the
# namespace and database (`urlshortener-data`) and the stateless service,
# migration and smoke test (`urlshortener-app`). Same split as socialapp.
#
# Both Applications render only from this repo: `urlshortener-data` from
# ../deploy/manifests, `urlshortener-app` from ../helm/urlshortener with its
# values.homelab.yaml.

# Scoped narrowly. Infrastructure repo doc 06: "Scope projects to only their
# repositories, namespaces, and required resource kinds."
#
# The app renders a Deployment, a Service, Jobs and one namespaced
# CloudNativePG Cluster — nothing cluster-scoped but its own namespace. The
# CRDs it depends on are installed by the operator's Helm release, deliberately
# outside this project's reach.
resource "argocd_project" "urlshortener" {
  metadata {
    name      = "homelab-urlshortener"
    namespace = "argocd"
  }

  spec {
    # At most 255 bytes: Kubernetes rejects a longer AppProject description.
    description = "urlshortener's own resources: urlshortener-data (namespace, PostgreSQL) and urlshortener-app (service, migration, smoke test), both rendered from github.com/igomez10/microservices."

    # This repo only. Both Applications render from it.
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
      namespace = "urlshortener"
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

# Tracks ../deploy/manifests on mainline as a plain Kustomize directory: the
# namespace and the CloudNativePG Cluster.
#
# A separate Application from `urlshortener_app`, and the reason is the sync
# policy below: that one runs automated with prune and self_heal, and folding
# the database into it would put the Cluster under automatic pruning.
#
# Sync is MANUAL: this is a database, and an automated sync of a bad render with
# no one at the keyboard is worse than a stale OutOfSync badge.
#
# Requires, in this order, before the first sync:
#   1. the CloudNativePG operator (infrastructure layer 20)
#   2. the urlshortener-postgres-app Secret — see the infrastructure repo's
#      homelab/apps/urlshortener/README.md
resource "argocd_application" "urlshortener_data" {
  # Explicit, so removing this resource from Terraform deletes the Application
  # object and nothing else. A cascading delete would take every resource the
  # Application manages with it — including the database.
  cascade = false

  metadata {
    name      = "urlshortener-data"
    namespace = "argocd"
  }

  spec {
    project = argocd_project.urlshortener.metadata[0].name

    source {
      repo_url        = var.application_repo_url
      target_revision = "mainline"
      path            = "urlshortener/deploy/manifests"
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "urlshortener"
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

# The service, its schema migration Job and a post-sync smoke test, rendered
# from ../helm/urlshortener with values.homelab.yaml from the chart directory.
# That file holds everything cluster-specific, including the image digests.
#
# Sync is AUTOMATED, unlike `urlshortener-data`. A push to mainline builds
# images, the image workflow commits their digests to values.homelab.yaml on
# mainline, and Argo CD deploys that commit. That commit cannot start another
# build: urlshortener/helm/** is excluded from the image workflow's paths, and
# a GITHUB_TOKEN push never starts workflow runs.
#
# Every such sync runs the migration Job against the production database
# unattended; that is acceptable because it is a checksum-guarded no-op unless
# db/setup/schema.sql changed, and the PostSync smoke test fails the sync when
# the result does not serve.
resource "argocd_application" "urlshortener_app" {
  # Explicit, so removing this resource from Terraform deletes the Application
  # object and nothing else.
  cascade = false

  metadata {
    name      = "urlshortener-app"
    namespace = "argocd"
  }

  spec {
    project = argocd_project.urlshortener.metadata[0].name

    source {
      repo_url        = var.application_repo_url
      target_revision = "mainline"
      path            = "urlshortener/helm/urlshortener"

      helm {
        # The chart's object names are fixed (urlshortener,
        # urlshortener-migrate, ...) rather than derived from this, so it is
        # cosmetic — but it is the Helm release name the templates see.
        release_name = "urlshortener"
        # Relative to the chart directory, and applied over its values.yaml.
        value_files = ["values.homelab.yaml"]
      }
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "urlshortener"
    }

    sync_policy {
      # `prune` because removing something from the chart should actually
      # remove it; nothing here owns a PersistentVolume. `self_heal` so a
      # hand edit to the Deployment is reverted to what Git says.
      automated {
        prune     = true
        self_heal = true
      }

      sync_options = [
        # The namespace belongs to the `urlshortener-data` Application.
        "CreateNamespace=false",
        "ServerSideApply=true",
      ]
    }
  }
}
