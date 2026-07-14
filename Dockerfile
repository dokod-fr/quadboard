# syntax=docker/dockerfile:1
# stage 1: Builder
FROM golang:1.26-alpine AS builder

# Install dependencies to build
RUN go install github.com/go-task/task/v3/cmd/task@latest \
    && go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# Variable definition from GitHub Actions
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

# On lance la compilation via Task en injectant nos variables d'environnement.
# BuildKit va utiliser le cache Go et le cache des modules.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    VERSION=${VERSION} COMMIT=${COMMIT} DATE=${DATE} task release

# stage 2: Final runner image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Le binaire a été généré par Task dans "./bin/quadboard"
COPY --from=builder /app/bin/quadboard /usr/local/bin/quadboard

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/quadboard"]
CMD ["serve"]