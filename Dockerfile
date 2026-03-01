# =====================================================
# File: Dockerfile
# Repository: foxly-it/rootguard-webapp
# Purpose: Full-Stack Multi-Arch Image (Frontend + Backend)
# =====================================================

# ==========================================
# Stage 1 – Frontend Builder (React + Vite)
# ==========================================
FROM node:22-bookworm AS frontend-builder

WORKDIR /frontend

# Package Definition zuerst kopieren (Layer Cache)
COPY frontend/package.json frontend/package-lock.json ./

# Install Production Dependencies
RUN npm ci

# Restliches Frontend kopieren
COPY frontend/ ./

# Production Build (output: /frontend/dist)
RUN npm run build


# ==========================================
# Stage 2 – Backend Builder (Go)
# ==========================================
FROM golang:1.26-bookworm AS backend-builder

ARG VERSION=dev
ARG COMMIT=unknown
ARG TARGETARCH

WORKDIR /app

# Backend Source kopieren
COPY backend/ ./

# Frontend Build integrieren
COPY --from=frontend-builder /frontend/dist ./web

# Multi-Arch Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build \
    -ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT" \
    -o rootguard-webapp ./cmd/rootguard-webapp


# ==========================================
# Stage 3 – Minimal Runtime (Distroless)
# ==========================================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Binary + Web Assets übernehmen
COPY --from=backend-builder /app/rootguard-webapp .
COPY --from=backend-builder /app/web ./web

EXPOSE 8080

# Security: non-root runtime
USER nonroot:nonroot

ENTRYPOINT ["/app/rootguard-webapp"]