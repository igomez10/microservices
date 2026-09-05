# The socialapp CLI as an image, for machines rather than laptops.
#
# It exists so a Kubernetes Job can create roles and scopes through the same
# code path a human uses, instead of a bespoke script issuing raw HTTP calls or
# — worse — writing directly into socialapp's database from another service.
# See homelab/apps/puttyknife/manifests/socialapp-bootstrap.yaml in
# microservices-infrastructure, which is the only consumer today.
#
# Separate from Dockerfile, which builds the SERVER. Same repo, same module,
# different binary: the server image has no reason to carry a CLI, and this
# image has no reason to carry the server.
FROM --platform=$BUILDPLATFORM golang AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off and an explicit target so the result is a static binary that runs on
# the scratch-like base below and can be cross-compiled without a C toolchain.
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /socialapp-cli ./cmd/cli

FROM alpine:latest

# ca-certificates for TLS to a real deployment; jq because the CLI speaks JSON
# and the bootstrap script needs to read ids out of it; bash for `set -o
# pipefail`, which ash does not support.
RUN apk add --no-cache ca-certificates jq bash

COPY --from=builder /socialapp-cli /usr/local/bin/socialapp-cli

# Runs as a Kubernetes Job with runAsNonRoot; this uid does not need to exist in
# /etc/passwd for a static binary.
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/socialapp-cli"]
