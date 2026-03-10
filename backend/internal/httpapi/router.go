// =====================================================
// File: backend/internal/httpapi/router.go
// Project: RootGuard WebApp
// Purpose: Central HTTP router configuration
//
// Responsibilities:
//
// - register API routes
// - expose health endpoints
// - expose dashboard endpoints
// - expose system overview endpoint
// - serve frontend SPA
//
// Uses:
//
// - Go standard library net/http ServeMux
//
// =====================================================

package httpapi

import (
	"net/http"
	"path/filepath"

	"github.com/foxly-it/rootguard-webapp/backend/internal/api"
)

// =====================================================
// NewRouter
//
// Registers:
//
// /health
// /api/health
// /api/version
// /api/dashboard
// /api/system
// /api/service/{name}/{action}
//
// Also serves the frontend SPA.
// =====================================================

func NewRouter() http.Handler {

	mux := http.NewServeMux()

	// ==================================================
	// Health Endpoints
	// ==================================================

	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/health", HealthHandler)

	// ==================================================
	// Version Endpoint
	// ==================================================

	mux.HandleFunc("/api/version", VersionHandler)

	// ==================================================
	// Dashboard Endpoint
	// ==================================================

	mux.HandleFunc("/api/dashboard", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		api.HandleDashboard(w, r)

	})

	// ==================================================
	// System Overview Endpoint
	//
	// Provides aggregated Docker metrics:
	//
	// - docker version
	// - container count
	// - running containers
	// - total cpu
	// - total memory
	//
	// ==================================================

	mux.HandleFunc("/api/system", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		api.HandleSystem(w, r)

	})

	// ==================================================
	// Service Control Endpoint
	//
	// Example:
	// POST /api/service/adguard/start
	//
	// ==================================================

	mux.HandleFunc("/api/service/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		api.HandleServiceAction(w, r)

	})

	// ==================================================
	// Static Frontend (SPA)
	// ==================================================

	fileServer := http.FileServer(http.Dir("./web"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// If request is API → ignore
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			http.NotFound(w, r)
			return
		}

		// Root → index.html
		if r.URL.Path == "/" {

			http.ServeFile(
				w,
				r,
				filepath.Join("web", "index.html"),
			)

			return
		}

		// Static assets
		fileServer.ServeHTTP(w, r)

	})

	return mux
}
