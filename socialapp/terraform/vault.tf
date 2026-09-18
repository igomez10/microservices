# ── Vault ────────────────────────────────────────────────────────────────────
# The policy scoping socialapp's slice of the shared KV mount, and the
# kubernetes auth role that lets the Vault Secrets Operator use it. The mount
# and the auth backend themselves belong to layer 30 of the infrastructure
# repo; this root only hangs a policy and a role off the paths it publishes
# (remote-state.tf).

# Read-only access to socialapp's credentials, for whatever ends up consuming
# them — the Vault Secrets Operator, an injected sidecar, or a human. Covers
# the whole kv/socialapp/* prefix, which today is postgresql and redis; a
# third datastore lands under the same policy without editing it.
#
# The values themselves are deliberately not Terraform-managed. A
# vault_kv_secret_v2 resource would write every password into the state file
# in the GCS backend in cleartext — a second copy, less protected than Vault
# itself, which defeats the point of putting them there. They are written with
# `vault kv put`; see homelab/apps/socialapp/README.md in the infrastructure
# repo.
resource "vault_policy" "socialapp_databases" {
  name = "socialapp-databases"

  policy = <<-EOT
    path "${local.vault_kv_mount_path}/data/socialapp/*" {
      capabilities = ["read"]
    }

    path "${local.vault_kv_mount_path}/metadata/socialapp/*" {
      capabilities = ["read", "list"]
    }
  EOT
}

# Lets the Vault Secrets Operator read this app's credentials on behalf of
# workloads in the socialapp namespace, and nothing else: the role is bound to
# that namespace and grants only the read policy above.
resource "vault_kubernetes_auth_backend_role" "socialapp" {
  backend   = local.vault_kubernetes_auth_path
  role_name = "socialapp"

  bound_service_account_names      = ["default"]
  bound_service_account_namespaces = ["socialapp"]

  token_policies = [vault_policy.socialapp_databases.name]
  token_ttl      = 3600
}
