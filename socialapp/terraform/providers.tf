provider "vault" {
  address = var.vault_address

  # Authentication is intentionally supplied through VAULT_TOKEN. Never put a
  # Vault token in Terraform configuration or a committed variable file.
}

provider "kafka" {
  # The remote stack's PLAINTEXT_HOST listeners.
  bootstrap_servers = var.kafka_bootstrap_servers
  tls_enabled       = false
}

# The Argo CD and Grafana API tokens, read ephemerally so neither lands in this
# root's state.
# It shares kv/semaphore/terraform with the infrastructure layers rather than
# getting a path of its own: it is the same Argo CD local account, and a second
# copy of the same token is a second thing to rotate.
#
# The practical consequence is that this root's VAULT_TOKEN needs read on
# kv/semaphore/terraform — which the Semaphore Terraform environment's token
# already has, so the socialapp pipeline works with the same credential set as
# the infrastructure pipeline. A local run needs a token with that read too.
ephemeral "vault_kv_secret_v2" "semaphore_terraform" {
  mount = local.vault_kv_mount_path
  name  = "semaphore/terraform"
}

# The homelab cluster's Grafana — see var.grafana_url, and grafana.tf for what
# this root puts there.
#
# Its token comes from the same ephemeral read as Argo CD's, so it never lands
# in state: a service account token under kv/semaphore/terraform →
# grafana_homelab_auth, the same key layer 30 of the infrastructure repo
# expects for this instance (TF_VAR_grafana_homelab_auth there). Editor role is
# enough for folders, dashboards, alert rules and contact points.
provider "grafana" {
  url  = var.grafana_url
  auth = ephemeral.vault_kv_secret_v2.semaphore_terraform.data["grafana_homelab_auth"]
}

# grpc_web because the gateway sits in front of argocd-server: plain gRPC
# through a proxy is where this provider's connection errors usually come from.
provider "argocd" {
  server_addr = var.argocd_server_addr
  auth_token  = ephemeral.vault_kv_secret_v2.semaphore_terraform.data["argocd_auth_token"]
  grpc_web    = true
}
