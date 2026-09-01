# ---------- Stage 1: build ----------
FROM golang:1.26-alpine AS builder

# git needed if any deps pull from git; ca-certs for HTTPS module fetches
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go.mod/go.sum FIRST and download deps — this layer caches unless deps change
COPY go.mod go.sum ./
RUN go mod download

# Install swag CLI so we can generate docs inside the build
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Now copy the source
COPY . .

# Generate swagger docs (the docs package your code imports)
RUN swag init -g cmd/server/main.go -o infra/http/swagger/docs

# Build a static binary: no CGO, stripped, targeting linux
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server/main.go

# ---------- Stage 2: runtime ----------
FROM alpine:3.20

# ca-certificates so the app can make outbound HTTPS calls (DB TLS, APIs)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy ONLY the compiled binary from the build stage
COPY --from=builder /app/server /app/server

# Copy migrations so they can run in the container if needed
COPY --from=builder /build/migrations /app/migrations

# Run as non-root — never run a service as root
RUN adduser -D -u 10001 appuser
USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/server"]
