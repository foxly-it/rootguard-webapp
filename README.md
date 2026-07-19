# 🛡 RootGuard WebApp

![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)
![Backend](https://img.shields.io/badge/backend-Go-00ADD8?logo=go)
![Frontend](https://img.shields.io/badge/frontend-React-61DAFB?logo=react)
![Docker](https://img.shields.io/badge/docker-ready-2496ED?logo=docker)
![Status](https://img.shields.io/badge/status-active--development-orange)

---

## 📌 Overview

RootGuard WebApp is the authenticated HTTP and UI layer of the RootGuard
ecosystem. It contains no direct Docker control code and communicates with
RootGuard Core over an internal token-protected API.

It exposes infrastructure orchestration capabilities provided by **RootGuard Core** and delivers a modern web-based control plane for DNS stack management.

The WebApp is intentionally separated from the engine layer to maintain:

- Clear responsibility boundaries  
- Replaceable transport layers  
- Modular evolution  
- Security isolation  

---

## 🏗 System Architecture

`````
Browser
   ↓
RootGuard WebApp
   ├── REST API (Go)
   └── UI (React + Vite)
   ↓
RootGuard Core (Engine)
   ↓
Docker / System Services
`````

RootGuard WebApp acts purely as a presentation and API transport layer.

---

## 📂 Repository Structure

`````
rootguard-webapp/
├── backend/
│   ├── cmd/
│   │   └── rootguard-webapp/
│   ├── internal/
│   │   ├── httpapi/
│   │   └── unboundctl/
│   └── go.mod
│
├── frontend/
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── vite.config.ts
│
└── Dockerfile
`````

---

## 🚀 Local Development

### Backend

````bash
cd backend
ROOTGUARD_CORE_URL=http://localhost:8081 \
ROOTGUARD_API_TOKEN=development-token \
ROOTGUARD_ADMIN_PASSWORD=change-me \
go run ./cmd/rootguard-webapp
````

### Frontend

````bash
cd frontend
npm install
npm run dev
````

---

## 🐳 Containerized Deployment

### Build Image

````bash
docker build -t rootguard-webapp:dev .
````

### Run Container

````bash
docker run -p 8080:8080 \
  -e ROOTGUARD_CORE_URL=http://rootguard-core:8081 \
  -e ROOTGUARD_API_TOKEN=replace-me \
  -e ROOTGUARD_ADMIN_PASSWORD=replace-me \
  rootguard-webapp:dev
````

### Test Endpoints

````bash
curl http://localhost:8080/health
curl -u admin:replace-me http://localhost:8080/api/version
````

---

## 🔐 Design Principles

RootGuard WebApp follows strict design constraints:

- No orchestration logic duplication  
- No direct system-level manipulation  
- API-only exposure of infrastructure actions  
- Minimal runtime surface (distroless container)  
- Deterministic build process  

---

## 🛣 Roadmap

Planned development stages:

- Core integration via an internal HTTP API
- DNS configuration UI  
- Stack deployment interface  
- Live container monitoring  
- Session-based authentication and role management
- Multi-architecture image builds  
- GitHub Actions CI/CD  
- GHCR image distribution  

---

## 📜 License

RootGuard WebApp is licensed under the Apache License 2.0.

See the `LICENSE` file for full details.

---

## ⚠ Development Status

This project is currently under active development.

Breaking changes may occur prior to v1.0.0.
