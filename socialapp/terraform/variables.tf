variable "state_bucket" {
  type        = string
  description = "GCS bucket holding the infrastructure layers' state. Same bucket as this root's own backend; named here because a backend block cannot be referenced from configuration."
  default     = "microservices-341219"
  nullable    = false
}

# No default on purpose. This repo is public, and the domain is kept out of
# it: set it in terraform.tfvars (gitignored) for a local run, and as
# TF_VAR_domain in the Semaphore environment.
variable "domain" {
  type        = string
  description = "Base domain every hostname here is built from: the Vault and Argo CD APIs (vault. / argocd.homelab.internal.<domain>), the Kafka brokers (brokerN.internal.<domain>), and the public socialapp hostnames passed to the chart. Matches EXPOSED_HOST in the remote docker-compose .env."
  nullable    = false
}

variable "application_repo_url" {
  type        = string
  description = "HTTPS URL of this repo, which both Argo CD Applications render from. Public, so Argo CD needs no repository credentials for it — an SSH URL would demand a key even to read."
  default     = "https://github.com/igomez10/microservices.git"
  nullable    = false
}
