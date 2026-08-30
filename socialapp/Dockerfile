FROM --platform=$BUILDPLATFORM golang AS builder
WORKDIR /go/src/github.com/igomez10/microservices
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# TARGETARCH is supplied by the builder from the requested platform, so this
# image builds correctly for amd64 and arm64 alike. It was hardcoded to arm64,
# which produced a binary that cannot exec on an amd64 host.
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o /main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
COPY --from=builder /main /main
# cmd/main.go reads "openapi.yaml" relative to the working directory and exits
# 1 if it is missing. docker-compose bind-mounts it, so the image never needed
# to carry it; anything else running this image — Kubernetes included — does.
COPY openapi.yaml /openapi.yaml
WORKDIR /
CMD ["/main"]
