FROM golang AS builder
WORKDIR /go/src/github.com/igomez10/microservices
COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...
RUN go install -v ./...
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /main cmd/main.go

FROM scratch
COPY --from=builder /main /main
CMD ["/main"]
