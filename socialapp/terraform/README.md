# socialapp — Terraform

socialapp's own Terraform root. It holds the infrastructure objects that
describe *how socialapp deploys*, in the repo that holds socialapp's code, so
a change to a flag and a change to how that flag reaches the cluster are one
commit and one review.

It was `layers/40-apps/modules/socialapp/` in
[`igomez10/microservices-infrastructure`](https://github.com/igomez10/microservices-infrastructure)
until it moved here.

## What is in here

| File | Resources |
|---|---|
| `vault.tf` | `vault_policy.socialapp_databases` — read on `kv/socialapp/*`; `vault_kubernetes_auth_backend_role.socialapp` — lets the Vault Secrets Operator use it for the `socialapp` namespace |
| `argocd.tf` | `argocd_project.homelab-socialapp`, and the `socialapp-data` and `socialapp-app` Applications |
| `kafka.tf` | The CDC topics. Commented out; the provider is configured but never contacted |
| `grafana.tf` | On the **homelab** Grafana: the `socialapp` folder and three dashboards (`dashboards/overview.json`, `resources.json`, `logs.json`), 16 alert rules, and the `socialapp-slack` contact point they route to directly (no shared notification policy) |
| `zitadel.tf` | `zitadel_machine_user.socialapp_api` — the identity socialapp calls other services as — and its `zitadel_machine_key`, registered by public key only, read from `kv/socialapp/zitadel-public` (see **ZITADEL key** below) |
| `remote-state.tf` | The one read of the infrastructure repo's state: layer 30's two Vault paths |

Its state is `terraform/state/socialapp` in the `microservices-341219` GCS
bucket — the same bucket the infrastructure layers use, a prefix of its own.

## What the Applications render

Everything both Applications render is in this repo:

- `../deploy/manifests/` — the namespace, the CloudNativePG `Cluster`, Redis.
  What `socialapp-data` syncs, by hand.
- `../helm/socialapp/` with `values.homelab.yaml` — the API, frontend,
  migration Job and smoke test, with the image digests the image workflow pins
  on each build. What `socialapp-app` syncs, automatically.

## What is *not* in here

Layers 20 and 30 of the infrastructure repo still own everything underneath:
the Vault KV mount and kubernetes auth backend this root hangs a policy off,
the CloudNativePG operator, Argo CD itself.

## Credentials

- `VAULT_TOKEN` — for the `vault` provider, and it needs read on
  `kv/semaphore/terraform`, because that is where the Argo CD API token comes
  from (`providers.tf`).
- Google application-default credentials — for this root's own state and for
  the layer 30 remote-state read.
- Two Vault values for Grafana: a service account token (role Editor) at
  `kv/semaphore/terraform` → `grafana_homelab_auth`, read ephemerally; and the
  Slack webhook at `kv/alerting/socialapp` → `slack_webhook_url`, read by a data
  source and therefore stored in this root's state (the contact point's `url`
  is not write-only). Store a token without it passing through your clipboard
  history or shell history, e.g. `pbpaste | vault kv patch kv/semaphore/terraform grafana_homelab_auth=-`.
- A ZITADEL personal access token at `kv/semaphore/terraform` →
  `zitadel_access_token`, read ephemerally, for `zitadel.tf`. It is the same
  IAM admin PAT layer 30 of the infrastructure repo uses (the `iam-admin-pat`
  Secret in the `zitadel` namespace), so it is instance-wide.
- Six host variables, none with a default, because this repo is public and
  the domain is kept out of it. For a local run put them in `terraform.tfvars`
  (gitignored):
  ```hcl
  domain                  = "<domain>"                 # public hostnames passed to the chart
  vault_address           = "https://vault.<...>"      # scheme included
  argocd_server_addr      = "argocd.<...>:443"         # host:port, no scheme
  kafka_bootstrap_servers = ["broker1.<...>:29092", "broker2.<...>:29093", "broker3.<...>:29094"]
  grafana_url             = "https://grafana.homelab.<...>" # the cluster's Grafana, not the OCI VM's
  zitadel_domain          = "zitadel.<...>"            # host only, no scheme
  ```
  In Semaphore they are `TF_VAR_<name>` on the socialapp project's
  environment, set in `layers/40-apps/modules/semaphore/project-socialapp.tf`
  in the infrastructure repo (private).
- Network reachability: the Vault and Argo CD APIs are on the tailnet.

## Running it

```bash
terraform init
terraform plan
terraform apply
```

In Semaphore: the **socialapp** project, template **socialapp**. That project
is defined in the infrastructure repo, at
`layers/40-apps/modules/semaphore/project-socialapp.tf` — every `semaphoreui`
resource lives in the one module so a single place holds the Semaphore API
token. Its *Terraform provider credentials* environment is created empty; fill
in `VAULT_TOKEN` and the Google credentials through the Semaphore UI so they
are not written to layer 40's state.

## ZITADEL key

`socialapp-api` authenticates with an RSA key pair that is generated outside
Terraform. Its two halves live at two Vault paths:

| Path | Keys | Read by |
|---|---|---|
| `kv/socialapp/zitadel` | `private_key`, `user_id`, `key_id` | socialapp, through the Vault Secrets Operator |
| `kv/socialapp/zitadel-public` | `public_key` | this root, which registers it with ZITADEL |

Keep them apart: this root reads the public path with a data source, which
writes the whole secret into state — the private key must never be at that
path. This root's `VAULT_TOKEN` needs read on `kv/socialapp/zitadel-public`.

Creating or rotating the pair (the private key never touches disk):

```bash
key=$(openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048)
printf '%s\n' "$key" | vault kv put kv/socialapp/zitadel private_key=-
printf '%s\n' "$key" | openssl pkey -pubout | vault kv put kv/socialapp/zitadel-public public_key=-
unset key

terraform apply
vault kv patch kv/socialapp/zitadel \
  user_id="$(terraform output -raw zitadel_user_id)" \
  key_id="$(terraform output -raw zitadel_key_id)"
```

A rotation replaces the ZITADEL key, so the old private key stops working as
soon as the apply finishes. Patch `key_id` straight after, and restart
socialapp once the Vault Secrets Operator has synced the new values.

## One-time migration from layer 40

Do this **once**, and do the state move before applying either side. The five
objects already exist in the cluster and in Vault; nothing here should be
created or destroyed, only re-homed.

The infrastructure repo's `layers/40-apps/removed.tf` makes the order
forgiving — if layer 40 is applied first, the module is dropped from its state
without destroying anything — but the sequence below avoids having to import.

```bash
# 1. Lift the five resources out of layer 40's state into a local file,
#    renaming them from module.socialapp.X to X on the way out.
cd ~/microservices-infrastructure/layers/40-apps
terraform init
for r in vault_policy.socialapp_databases \
         vault_kubernetes_auth_backend_role.socialapp \
         argocd_project.socialapp \
         argocd_application.socialapp_data \
         argocd_application.socialapp_app; do
  terraform state mv -state-out=/tmp/socialapp.tfstate "module.socialapp.$r" "$r"
done

# 2. Push that file into this root's (empty) state.
cd ~/microservices/socialapp/terraform
terraform init
terraform state push -force /tmp/socialapp.tfstate
rm /tmp/socialapp.tfstate

# 3. Plan both sides.
terraform plan
terraform -chdir=~/microservices-infrastructure/layers/40-apps plan
```

`-force` on the push is because the pushed state carries layer 40's lineage
and this root's freshly initialised state has its own; there is nothing in
this state to lose.

Step 3 is the check that matters. Expect exactly **three in-place updates**
(`0 to add, 3 to change, 0 to destroy`):

- `argocd_application.socialapp_data` — its source moved from the
  infrastructure repo's `homelab/apps/socialapp/manifests` to this repo's
  `socialapp/deploy/manifests` (`repo_url`, `path`, `target_revision`);
- `argocd_application.socialapp_app` — the second `source` (the
  infrastructure repo as `$repo`) is removed, `value_files` becomes
  `values.homelab.yaml`, and `values` gains the six public hostnames built
  from `var.domain`. The rendered manifests are byte-identical to before;
- `argocd_project.socialapp` — `description` only.

The infrastructure repo also has to leave the project's `source_repos`, and
Terraform cannot do that: its Argo CD account is not allowed to add or remove
SSH-registered repos (see the comment on `source_repos` in `argocd.tf`). An
apply that tries fails with `permission denied: repositories, update` before
changing anything. Remove it as an admin, AFTER the apply has moved both
Applications off it:

```bash
KUBECONFIG=~/.kube/homelab kubectl -n argocd patch appproject homelab-socialapp --type=json \
  -p '[{"op":"test","path":"/spec/sourceRepos/0","value":"git@github.com:igomez10/microservices-infrastructure.git"},{"op":"remove","path":"/spec/sourceRepos/0"}]'
```

Done before the apply instead, both Applications report "repository not
permitted" until the apply lands; nothing is synced or pruned meanwhile.

Anything else is wrong:

- a *create* of any of the five means the state move did not land, and
  applying would fail against objects that already exist;
- a *replace* (`-/+`) of `socialapp_data` means the provider treats the source
  change as immutable. `cascade = false` keeps the database safe even then,
  but stop and look before applying.

`socialapp/deploy/manifests` must already be pushed to `mainline` before this
apply, and the infrastructure repo's copy must not be deleted until after it —
otherwise `socialapp-data` points at a path that does not exist.

Afterwards, `layers/40-apps/removed.tf` in the infrastructure repo is a no-op
and can be deleted.
