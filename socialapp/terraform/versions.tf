# socialapp's own Terraform root.
#
# This used to be layers/40-apps/modules/socialapp/ in the
# microservices-infrastructure repo. It lives here instead so that a change to
# how socialapp is deployed — its Argo CD Applications, its Vault policy, its
# Kafka topics — is one commit in the repo that holds the code, reviewed by
# whoever is changing the app.
#
# What did NOT move: homelab/apps/socialapp/ in the infrastructure repo still
# holds the namespace, the CloudNativePG Cluster, Redis, and the Helm values
# file the image workflow pins digests into. The Argo CD Applications below
# point at that repo by URL. The boundary is "cluster-shaped YAML lives with
# the cluster; the API objects that deploy socialapp live with socialapp".
#
# Its own state, like every other root here — never share a state file with the
# infrastructure layers. The one thing it reads from them is layer 30's two
# Vault paths (remote-state.tf).
terraform {
  required_version = ">= 1.11.0"

  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.0"
    }
    kafka = {
      source  = "Mongey/kafka"
      version = "0.13.1"
    }
    argocd = {
      source  = "argoproj-labs/argocd"
      version = "7.17.0"
    }
    # Same major as layer 30 of the infrastructure repo, which configures the
    # same Grafana instance's login.
    grafana = {
      source  = "grafana/grafana"
      version = "~> 3.0"
    }
  }

  # Same bucket as the infrastructure layers, a prefix of its own. Sharing the
  # bucket is deliberate: one place to find every state, one set of bucket
  # permissions, and the layer 30 remote-state read below needs access to it
  # anyway.
  backend "gcs" {
    bucket = "microservices-341219"
    prefix = "terraform/state/socialapp"
  }
}
