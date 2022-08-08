# syntax=docker/dockerfile:1
FROM public.ecr.aws/docker/library/golang:1.19 AS builder
WORKDIR /application
COPY . .
RUN GOPROXY=direct go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o app .


FROM alpine:3.15.4
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /application ./
CMD ["./app"]  

