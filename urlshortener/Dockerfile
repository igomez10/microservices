FROM golang AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...
COPY . .
RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]

