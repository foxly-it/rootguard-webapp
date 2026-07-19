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
// - expose service detection endpoint
// - expose service control endpoint
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
	"github.com/foxly-it/rootguard-webapp/backend/internal/coreclient"
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
// /api/services
// /api/service/{name}/{action}
//
// Also serves the frontend SPA.
// =====================================================

func NewRouter(core *coreclient.Client) http.Handler {

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

		api.HandleDashboard(w, r, core)

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

		api.HandleSystem(w, r, core)

	})

	// ==================================================
	// Service Detection Endpoint
	//
	// GET /api/services
	//
	// Returns detected DNS related services
	// detected by the RootGuard service engine.
	//
	// Example response:
	//
	// [
	//   {
	//     "name": "AdGuard Home",
	//     "type": "binary",
	//     "status": "running"
	//   },
	//   {
	//     "name": "Unbound",
	//     "type": "binary",
	//     "status": "running"
	//   }
	// ]
	//
	// ==================================================

	mux.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		api.HandleServices(w, r, core)

	})

	// ==================================================
	// Service Control Endpoint
	//
	// Example:
	// POST /api/service/adguard/start
	// POST /api/service/unbound/restart
	//
	// ==================================================

	mux.HandleFunc("/api/service/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		api.HandleServiceAction(w, r, core)

	})

	mux.HandleFunc("/api/unbound/settings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetUnboundSettings(w, r, core)
		case http.MethodPut:
			api.HandlePutUnboundSettings(w, r, core)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/adguard/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		api.HandleGetAdGuardStatus(w, r, core)
	})

	mux.HandleFunc("/api/adguard/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		api.HandleBootstrapAdGuard(w, r, core)
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

		assetPath := filepath.Join("web", filepath.Clean(r.URL.Path))
		if info, err := filepath.Glob(assetPath); err == nil && len(info) > 0 {
			fileServer.ServeHTTP(w, r)
			return
		}

		// BrowserRouter fallback for client-side routes.
		http.ServeFile(w, r, filepath.Join("web", "index.html"))

	})

	return mux
}
