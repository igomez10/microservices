FROM golang AS builder
WORKDIR /go/src/github.com/igomez10/microservices
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
COPY --from=builder /main /main
CMD ["/main"]
