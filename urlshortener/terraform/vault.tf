# ── Vault ────────────────────────────────────────────────────────────────────
# The policy scoping urlshortener's slice of the shared KV mount, and the
# kubernetes auth role that lets the Vault Secrets Operator use it. The mount
# and the auth backend belong to layer 30 of the infrastructure repo; this root
# only hangs a policy and a role off the paths it publishes (remote-state.tf).

# Read-only access to urlshortener's credentials, for whatever ends up
# consuming them — the Vault Secrets Operator, an injected sidecar, or a human.
# Covers the whole kv/urlshortener/* prefix, which today is just postgresql.
#
# The values themselves are deliberately not Terraform-managed. A
# vault_kv_secret_v2 resource would write every password into the state file
# in the GCS backend in cleartext — a second copy, less protected than Vault
# itself. They are written with `vault kv put`; see the infrastructure repo's
# homelab/apps/urlshortener/README.md.
resource "vault_policy" "urlshortener_databases" {
  name = "urlshortener-databases"

  policy = <<-EOT
    path "${local.vault_kv_mount_path}/data/urlshortener/*" {
      capabilities = ["read"]
    }

    path "${local.vault_kv_mount_path}/metadata/urlshortener/*" {
      capabilities = ["read", "list"]
    }
  EOT
}

# Lets the Vault Secrets Operator read this app's credentials on behalf of
# workloads in the urlshortener namespace, and nothing else: the role is bound
# to that namespace and grants only the read policy above.
resource "vault_kubernetes_auth_backend_role" "urlshortener" {
  backend   = local.vault_kubernetes_auth_path
  role_name = "urlshortener"

  bound_service_account_names      = ["default"]
  bound_service_account_namespaces = ["urlshortener"]

  token_policies = [vault_policy.urlshortener_databases.name]
  token_ttl      = 3600
}
