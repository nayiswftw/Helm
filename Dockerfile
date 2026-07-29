# Build stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache dependency downloads
COPY helm-core/go.mod helm-core/go.sum ./helm-core/
WORKDIR /app/helm-core
RUN go mod download

# Copy source and build static Linux binary
WORKDIR /app
COPY helm-core/ ./helm-core/
WORKDIR /app/helm-core
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /helm ./cmd/helm

# Production stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /helm /usr/local/bin/helm

ENV HELM_PORT=":8080" \
    HELM_LOG_LEVEL="info" \
    HELM_PROC_PATH="/host/proc"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/helm"]
