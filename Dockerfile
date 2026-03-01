# ==========================================
# Stage 1 – Backend Builder
# ==========================================
FROM golang:1.26-bookworm AS backend-builder

WORKDIR /app

# Build arguments (werden von GitHub Actions gesetzt)
ARG VERSION=dev
ARG COMMIT=unknown
ARG TARGETARCH

# Module Files zuerst kopieren (Layer Cache)
COPY backend/go.mod backend/go.sum* ./

RUN go mod download

# Backend Source
COPY backend/ ./

# Multi-Arch Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build \
    -ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT" \
    -o rootguard-webapp ./cmd/rootguard-webapp


# ==========================================
# Stage 2 – Runtime (Distroless)
# ==========================================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Binary übernehmen
COPY --from=backend-builder /app/rootguard-webapp .

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/rootguard-webapp"]