# ── ZITADEL ──────────────────────────────────────────────────────────────────
# socialapp's own identity: a machine user it authenticates as when it calls
# other services, with the client-credentials-style JWT profile grant (it signs
# an assertion with its private key and trades it for an access token).
#
# The ZITADEL objects that let *people* log in to the homelab's services live in
# layer 30 of the infrastructure repo. This one is socialapp's alone, so it sits
# beside socialapp's Vault policy.

# The organization every homelab user, project and grant belongs to — the same
# id as local.zitadel_org_id in the infrastructure repo's layer 30. Not a
# secret; named explicitly so the user cannot land in whichever org the
# provider's token happens to belong to.
locals {
  zitadel_org_id = "388381702403194926"
}

# JWT access tokens rather than opaque ones, so a service socialapp calls can
# verify the token against ZITADEL's JWKS without an introspection round trip.
resource "zitadel_machine_user" "socialapp_api" {
  org_id            = local.zitadel_org_id
  user_name         = "socialapp-api"
  name              = "socialapp API"
  description       = "Identity of the socialapp API when it calls other services. Managed by socialapp/terraform in igomez10/microservices."
  access_token_type = "ACCESS_TOKEN_TYPE_JWT"
}

# The key is registered by its PUBLIC half only. The private key is generated
# outside Terraform and written straight to Vault, so ZITADEL never sees it and
# neither does this root's state — a with_secret client secret, a PAT or a
# ZITADEL-generated key would all have been stored in state in cleartext. See
# "ZITADEL key" in README.md for how the pair is made.
#
# The two halves sit at two Vault paths, and that split is load-bearing:
#   kv/socialapp/zitadel         → private_key, user_id, key_id  (the app's)
#   kv/socialapp/zitadel-public  → public_key                     (this root's)
# public_key is not a write-only argument, so it cannot take an ephemeral
# value, and this data source writes the WHOLE secret it reads into state.
# Pointed at kv/socialapp/zitadel it would put the private key there. Same
# deprecation warning as data.vault_kv_secret_v2.alerting, same reason to
# ignore it.
data "vault_kv_secret_v2" "zitadel_public" {
  mount = local.vault_kv_mount_path
  name  = "socialapp/zitadel-public"
}

# Rotating is: generate a new pair, overwrite both Vault paths, apply. The
# provider replaces the key (public_key is ForceNew, so a new key id), so patch
# key_id beside the private key straight after.
resource "zitadel_machine_key" "socialapp_api" {
  org_id     = local.zitadel_org_id
  user_id    = zitadel_machine_user.socialapp_api.id
  key_type   = "KEY_TYPE_JSON"
  public_key = data.vault_kv_secret_v2.zitadel_public.data["public_key"]
}

# Grants on the services socialapp calls go here, as zitadel_user_grant
# resources against those services' projects — none yet.

# Neither is a secret: together with the private key in Vault they make up the
# JWT profile assertion (`iss`/`sub` = user id, `kid` = key id).
output "zitadel_user_id" {
  value       = zitadel_machine_user.socialapp_api.id
  description = "ZITADEL user id of the socialapp-api machine user."
}

output "zitadel_key_id" {
  value       = zitadel_machine_key.socialapp_api.id
  description = "Id of socialapp-api's machine key — the `kid` of the assertions it signs."
}
