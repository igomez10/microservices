FROM golang AS builder
WORKDIR /go/src/github.com/igomez10/microservices
COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...
RUN go install -v ./...
COPY ./integration_tests ./integration_tests
COPY ./client ./client
CMD ["go", "test" , "-count", "1",  "-v", "./integration_tests/..."]
