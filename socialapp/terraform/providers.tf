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

# Argo CD's API token, read ephemerally so it never lands in this root's state.
# It shares kv/semaphore/terraform with the infrastructure layers rather than
# getting a path of its own: it is the same Argo CD local account, and a second
# copy of the same token is a second thing to rotate.
#
# The practical consequence is that this root's VAULT_TOKEN needs read on
# kv/semaphore/terraform — which the Semaphore Terraform environment's token
# already has, so the socialapp pipeline works with the same credential set as
# the infrastructure pipeline. A local run needs a token with that read too.
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
