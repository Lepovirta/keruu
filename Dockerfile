FROM golang:1.13-alpine AS builder

# User
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home /none \
    --shell /sbin/nologin \
    --no-create-home \
    --uid 10101 \
    appuser

# Git and certificates
RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Build
ENV GOARCH=amd64 GOOS=linux CGO_ENABLED=0
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o keruu

# Target image
FROM scratch
WORKDIR /workspace
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /build/keruu /keruu
USER appuser:appuser
ENTRYPOINT ["/keruu"]
