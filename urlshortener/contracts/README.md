# Pact Contracts

This folder stores the Pact artifacts produced when verifying the URL Shortener service.

- `pacts/` currently holds `socialapp-urlshortener.json`, which Socialapp emits when it calls `/v1/urls/{alias}/data`.
- `logs/` captures Pact mock server logs and is ignored by version control.

Run `make contract-test` from the `urlshortener/` directory after installing the Pact CLI to re-verify the provider against Socialapp's expectations.
