FROM golang AS builder
WORKDIR /go/src/github.com/igomez10/microservices
COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...
RUN go install -v ./...
COPY ./integration_tests ./integration_tests
COPY ./client ./client
RUN go build -o /go/src/github.com/igomez10/microservices/app ./integration_tests
CMD ["/go/src/github.com/igomez10/microservices/app"]
