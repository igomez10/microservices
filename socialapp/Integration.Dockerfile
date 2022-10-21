FROM golang AS builder
WORKDIR /go/src/github.com/igomez10/microservices
COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...
RUN go install -v ./...
COPY . .
CMD ["go", "test", "-v", "./integration_tests/..."]
