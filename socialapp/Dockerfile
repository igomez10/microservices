FROM golang AS builder
WORKDIR /go/src/github.com/igomez10/microservices
COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...
RUN go install -v ./...
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/igomez10/microservices/main .
CMD ["./main"]

