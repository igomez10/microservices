FROM golang AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 go build -o app ./integration_tests


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app ./app
CMD ["./app"]

