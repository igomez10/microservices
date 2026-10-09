# urlshortener — Terraform

urlshortener's own Terraform root: the objects that describe *how urlshortener
deploys*, in the repo that holds its code. Same arrangement as
[`../../socialapp/terraform`](../../socialapp/terraform/README.md), which has
the longer explanation.

It was `layers/40-apps/modules/urlshortener/` in
`igomez10/microservices-infrastructure` until it moved here.

## What is in here

| File | Resources |
|---|---|
| `vault.tf` | `vault_policy.urlshortener_databases` — read on `kv/urlshortener/*`; `vault_kubernetes_auth_backend_role.urlshortener` for the `urlshortener` namespace |
| `argocd.tf` | `argocd_project.homelab-urlshortener`, and the `urlshortener-data` and `urlshortener-app` Applications |
| `grafana.tf` | On the **homelab** Grafana: the `urlshortener` folder, three dashboards (`dashboards/overview.json`, `resources.json`, `logs.json`), 14 alert rules, and the `urlshortener-slack` contact point they route to directly (no shared notification policy) |
| `remote-state.tf` | The one read of the infrastructure repo's state: layer 30's two Vault paths |

State: `terraform/state/urlshortener` in the `microservices-341219` bucket.

## What the Applications render

Everything, from this repo:

- `../deploy/manifests/` — the namespace and the CloudNativePG `Cluster`.
  `urlshortener-data` syncs it, by hand.
- `../helm/urlshortener/` with `values.homelab.yaml` — the service, migration
  Job and smoke test, with the image digests the image workflow pins on each
  build. `urlshortener-app` syncs it, automatically.

## Credentials and variables

- `VAULT_TOKEN`, with read on `kv/semaphore/terraform` (the Argo CD API token
  comes from there — `providers.tf`).
- Two Vault values for Grafana: the service account token at
  `kv/semaphore/terraform` → `grafana_homelab_auth`, read ephemerally (the
  same one socialapp's root uses); and the Slack webhook at
  `kv/alerting/urlshortener` → `slack_webhook_url`, read by a data source and
  therefore stored in this root's state (the contact point's `url` is not
  write-only). Store it without it passing through shell history, e.g.
  `pbpaste | vault kv put kv/alerting/urlshortener slack_webhook_url=-`.
- Google application-default credentials, for this root's state and the
  layer 30 remote-state read.
- Three host variables, no defaults, because this repo is public and the
  domain is kept out of it. Locally, in `terraform.tfvars` (gitignored):
  ```hcl
  vault_address      = "https://vault.<...>"           # scheme included
  argocd_server_addr = "argocd.<...>:443"              # host:port, no scheme
  grafana_url        = "https://grafana.homelab.<...>" # the cluster's Grafana, not the OCI VM's
  ```
  In Semaphore they are `TF_VAR_<name>` on the `microservices` project's
  environment, set in the infrastructure repo.

## Running it

```bash
terraform init && terraform plan && terraform apply
```

In Semaphore: project **microservices**, template **urlshortener**.
