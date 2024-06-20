# Dockerfile for task-service

FROM golang:1.22.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -mod=readonly -o /go/bin/app ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/bin/app .

ENTRYPOINT ["./app"]
