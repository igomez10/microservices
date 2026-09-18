variable "state_bucket" {
  type        = string
  description = "GCS bucket holding the infrastructure layers' state. Same bucket as this root's own backend; named here because a backend block cannot be referenced from configuration."
  default     = "microservices-341219"
  nullable    = false
}

# Neither host variable has a default, on purpose. This repo is public and the
# domain is kept out of it: set them in terraform.tfvars (gitignored) for a
# local run, and as TF_VAR_<name> on the Semaphore environment
# (layers/40-apps/modules/semaphore/project-microservices.tf in the
# infrastructure repo).

variable "vault_address" {
  type        = string
  description = "URL of the Vault API, scheme included (https://...). On the tailnet."
  nullable    = false
}

variable "argocd_server_addr" {
  type        = string
  description = "Argo CD API as host:port, no scheme — what the argocd provider's server_addr expects. On the tailnet, through the Istio gateway."
  nullable    = false
}

variable "application_repo_url" {
  type        = string
  description = "HTTPS URL of this repo, which both Argo CD Applications render from. Public, so Argo CD needs no repository credentials for it — an SSH URL would demand a key even to read."
  default     = "https://github.com/igomez10/microservices.git"
  nullable    = false
}
