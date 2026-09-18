variable "state_bucket" {
  type        = string
  description = "GCS bucket holding the infrastructure layers' state. Same bucket as this root's own backend; named here because a backend block cannot be referenced from configuration."
  default     = "microservices-341219"
  nullable    = false
}

# None of the variables in this file that name a host has a default, on
# purpose. This repo is public and the domain is kept out of it: set them in
# terraform.tfvars (gitignored) for a local run, and as TF_VAR_<name> on the
# Semaphore environment (layers/40-apps/modules/semaphore/project-socialapp.tf
# in the infrastructure repo).

variable "domain" {
  type        = string
  description = "Base domain of the public socialapp hostnames passed to the chart (socialapp.<domain>, properties.<domain>, ...). Matches EXPOSED_HOST in the remote docker-compose .env."
  nullable    = false
}

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

variable "kafka_bootstrap_servers" {
  type        = list(string)
  description = "Kafka brokers as host:port — the remote stack's PLAINTEXT_HOST listeners. Only contacted once kafka.tf's topics are uncommented. As an environment variable, HCL list syntax: TF_VAR_kafka_bootstrap_servers='[\"a:29092\",\"b:29093\"]'."
  nullable    = false
}

variable "grafana_url" {
  type        = string
  description = "URL of the homelab cluster's Grafana (kube-prometheus-stack, infrastructure layer 20), scheme included. NOT the OCI VM's compose-stack Grafana: socialapp's dashboards and alerts live on the instance that sees the cluster it runs in. On the tailnet."
  nullable    = false
}

variable "application_repo_url" {
  type        = string
  description = "HTTPS URL of this repo, which both Argo CD Applications render from. Public, so Argo CD needs no repository credentials for it — an SSH URL would demand a key even to read."
  default     = "https://github.com/igomez10/microservices.git"
  nullable    = false
}
