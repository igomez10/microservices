FROM golang:1.11-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /src
COPY ./ ./
RUN apk add --no-cache ca-certificates git
RUN go mod download
RUN go build -installsuffix 'static' -o /server .

FROM alpine
EXPOSE 4444
COPY --from=builder /server /server
ENTRYPOINT ["/server"]
