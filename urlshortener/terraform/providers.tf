provider "vault" {
  address = var.vault_address

  # Authentication is intentionally supplied through VAULT_TOKEN. Never put a
  # Vault token in Terraform configuration or a committed variable file.
}

# The Argo CD and Grafana API tokens, read ephemerally so they never land in
# this root's state. The same path the infrastructure layers and
# ../../socialapp/terraform use, so this root's VAULT_TOKEN needs read on
# kv/semaphore/terraform — which the Semaphore environment's token has.
ephemeral "vault_kv_secret_v2" "semaphore_terraform" {
  mount = local.vault_kv_mount_path
  name  = "semaphore/terraform"
}

# grpc_web because the gateway sits in front of argocd-server: plain gRPC
# through a proxy is where this provider's connection errors usually come from.
provider "argocd" {
  server_addr = var.argocd_server_addr
  auth_token  = ephemeral.vault_kv_secret_v2.semaphore_terraform.data["argocd_auth_token"]
  grpc_web    = true
}

# The homelab cluster's Grafana — see var.grafana_url, and grafana.tf for what
# this root puts there. Same service account token as socialapp's root
# (kv/semaphore/terraform → grafana_homelab_auth); Editor role is enough for
# folders, dashboards, alert rules and contact points.
provider "grafana" {
  url  = var.grafana_url
  auth = ephemeral.vault_kv_secret_v2.semaphore_terraform.data["grafana_homelab_auth"]
}
