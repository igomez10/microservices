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
- `domain` — the base domain every hostname is built from. It has no
  default, because this repo is public and the domain is kept out of it. For
  a local run put it in `terraform.tfvars` (gitignored):
  ```hcl
  domain = "<your domain>"
  ```
  In Semaphore, set `TF_VAR_domain` on the environment.
- Network reachability: the Vault and Argo CD APIs are on the tailnet, at
  `vault.homelab.internal.<domain>` and `argocd.homelab.internal.<domain>`.

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
- `argocd_project.socialapp` — the infrastructure repo leaves `source_repos`,
  and the `description` is reworded.

The project updates first, because both Applications reference it. For the
seconds between that and the Application updates, Argo CD may report both
Applications as "repository not permitted". That is a comparison error: it
does not sync or prune anything, and it clears when the apply finishes.

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
