FROM golang:1.11-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /src
COPY ./ ./
RUN go mod download
RUN go build -installsuffix 'static' -o /server .

FROM alpine
EXPOSE 8080
COPY --from=builder /server /server
ENTRYPOINT ["/server"]
