# Layer 30 of the infrastructure repo owns the Vault KV mount and the
# kubernetes auth backend; this root owns one policy and one role underneath
# them. Reading the two paths out of its state is the whole coupling between
# the two repositories — no resource here is declared in terms of a layer 30
# resource, only in terms of two strings it publishes.
#
# This is the one place the app repo depends on the infrastructure repo's
# state. Keep it that way: each additional terraform_remote_state turns two
# independently deployable things back into one.
#
# Read-only, and it fails loudly rather than silently: if layer 30 has never
# been applied, this data source errors at plan time instead of producing an
# empty path that would create policies against a mount that does not exist.
data "terraform_remote_state" "platform_config" {
  backend = "gcs"

  config = {
    bucket = var.state_bucket
    prefix = "terraform/state/30-platform-config"
  }
}

locals {
  vault_kv_mount_path        = data.terraform_remote_state.platform_config.outputs.vault_kv_mount_path
  vault_kubernetes_auth_path = data.terraform_remote_state.platform_config.outputs.vault_kubernetes_auth_path
}
