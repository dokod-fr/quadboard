# Build
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /quadboard ./cmd/quadboard

# Runtime
FROM alpine:3.22

RUN apk add --no-cache ca-certificates

LABEL org.opencontainers.image.title="QuadBoard"
LABEL org.opencontainers.image.description="Dashboard for Quadlet resources"
LABEL org.opencontainers.image.source="https://github.com/dokod-fr/quadboard"
LABEL org.opencontainers.image.licenses="GPLv3"

WORKDIR /app

RUN adduser -D -u 1001 quadboard

COPY --from=builder /quadboard /app/quadboard

USER quadboard

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/app/quadboard"]