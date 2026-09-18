provider "vault" {
  address = var.vault_address

  # Authentication is intentionally supplied through VAULT_TOKEN. Never put a
  # Vault token in Terraform configuration or a committed variable file.
}

# Argo CD's API token, read ephemerally so it never lands in this root's state.
# The same local account and token path the infrastructure layers and
# ../../socialapp/terraform use, so this root's VAULT_TOKEN needs read on
# kv/semaphore/terraform — which the Semaphore environment's token has.
ephemeral "vault_kv_secret_v2" "argocd" {
  mount = local.vault_kv_mount_path
  name  = "semaphore/terraform"
}

# grpc_web because the gateway sits in front of argocd-server: plain gRPC
# through a proxy is where this provider's connection errors usually come from.
provider "argocd" {
  server_addr = var.argocd_server_addr
  auth_token  = ephemeral.vault_kv_secret_v2.argocd.data["argocd_auth_token"]
  grpc_web    = true
}
