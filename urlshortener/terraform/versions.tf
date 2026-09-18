# urlshortener's own Terraform root.
#
# This used to be layers/40-apps/modules/urlshortener/ in the
# microservices-infrastructure repo. It lives here so that a change to how
# urlshortener is deployed — its Argo CD Applications, its Vault policy — is
# one commit in the repo that holds the code. Same arrangement as
# ../../socialapp/terraform.
#
# Everything both Argo CD Applications render is in this repo too:
# ../deploy/manifests and ../helm/urlshortener with its values.homelab.yaml.
#
# Its own state, never shared with anything. The one thing it reads from the
# infrastructure repo is layer 30's two Vault paths (remote-state.tf).
terraform {
  required_version = ">= 1.11.0"

  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.0"
    }
    argocd = {
      source  = "argoproj-labs/argocd"
      version = "7.17.0"
    }
  }

  # Same bucket as the infrastructure layers, a prefix of its own.
  backend "gcs" {
    bucket = "microservices-341219"
    prefix = "terraform/state/urlshortener"
  }
}
