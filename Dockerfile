# syntax=docker/dockerfile:1
# stage 1: Builder
FROM golang:1.26-alpine AS builder

# Install dependencies to build (Task, Templ, Node.js pour Vite/Svelte)
RUN apk add --no-cache nodejs npm \
    && go install github.com/go-task/task/v3/cmd/task@latest 
    
WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Cache Node modules pour le frontend Svelte
COPY web/package.json web/package-lock.json* ./web/
RUN cd web && npm install

# Copy the rest of the source code
COPY . .

# Variable definition from GitHub Actions
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

# Build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    VERSION=${VERSION} COMMIT=${COMMIT} DATE=${DATE} task release

# stage 2: Final runner image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Binary in "./bin/quadboard"
COPY --from=builder /app/bin/quadboard /usr/local/bin/quadboard

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/quadboard"]
CMD ["serve"]