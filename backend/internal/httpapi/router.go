// =====================================================
// File: backend/internal/httpapi/router.go
// Purpose: Central HTTP router configuration
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
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/health", HealthHandler)
	mux.HandleFunc("/api/version", VersionHandler)

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

		// SPA Fallback → index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("web", "index.html"))
			return
		}

		fileServer.ServeHTTP(w, r)
	})

	return mux
}
