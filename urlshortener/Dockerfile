FROM --platform=$BUILDPLATFORM golang AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# TARGETARCH comes from the requested platform, so one Dockerfile builds for
# amd64 and arm64 without cross-compiling by hand.
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
# The binary was in /root/, which alpine ships as 0750 root:root — unreadable
# to any other user, so the image could only run as root. Kubernetes runs it
# with runAsNonRoot, so it lives at / instead.
COPY --from=builder /app/main /main
WORKDIR /
CMD ["/main"]
