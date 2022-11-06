FROM golang:1.18-alpine as builder
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /app
COPY ./DigiCertGlobalRootCA.crt.pem ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN go build -o consumer cmd/main.go

FROM alpine:3
RUN apk --no-cache add tzdata
WORKDIR /app
COPY ./DigiCertGlobalRootCA.crt.pem ./
COPY go.mod ./
COPY go.sum ./
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY --from=builder /app/consumer ./bin/consumer
ENTRYPOINT ["/app/bin/consumer"]