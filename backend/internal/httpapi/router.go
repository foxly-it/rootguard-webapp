// =====================================================
// File: backend/internal/httpapi/router.go
// Purpose: Central HTTP router configuration
// Uses net/http ServeMux (no external router dependency)
// =====================================================

package httpapi

import (
	"net/http"
	"path/filepath"
)

// =====================================================
// NewRouter()
// Registers API routes + static frontend
// =====================================================

func NewRouter() http.Handler {

	mux := http.NewServeMux()

	// ==========================
	// API Routes
	// ==========================

	// Health
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/health", HealthHandler)

	// Version
	mux.HandleFunc("/api/version", VersionHandler)

	// Dashboard (GET only)
	mux.HandleFunc("/api/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		HandleDashboard(w, r)
	})

	// Service Actions (POST only)
	mux.HandleFunc("/api/service/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		handleServiceAction(w, r)
	})

	// ==========================
	// Static Frontend (SPA)
	// ==========================

	fileServer := http.FileServer(http.Dir("./web"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Wenn API → direkt durchreichen
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			http.NotFound(w, r)
			return
		}

		// SPA Root → index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("web", "index.html"))
			return
		}

		// Static Assets
		fileServer.ServeHTTP(w, r)
	})

	return mux
}
