FROM golang:1.11-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /src
COPY ./ ./
RUN go mod download
RUN go build -installsuffix 'static' -o /server .

FROM alpine

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

EXPOSE 8080
COPY --from=builder /server /server
ENTRYPOINT ["/server"]
