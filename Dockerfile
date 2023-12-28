# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go run github.com/steebchen/prisma-client-go prefetch
RUN go run github.com/steebchen/prisma-client-go db push
RUN go run github.com/steebchen/prisma-client-go generate

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /build/apiserver .

# Final stage
FROM alpine:latest

COPY --from=builder /build/apiserver /apiserver

# Expose the port the application runs on
EXPOSE 8080

# Command to run when starting the container
ENTRYPOINT ["/apiserver"]
