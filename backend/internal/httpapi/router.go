// =====================================================
// File: backend/internal/httpapi/router.go
// Purpose: API + Static React SPA (panic-safe version)
// =====================================================

package httpapi

import (
	"net/http"
	"strings"
)

func NewRouter() http.Handler {

	mux := http.NewServeMux()

	// ==========================
	// API Routes
	// ==========================

	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/health", HealthHandler)

	// ==========================
	// Static React App
	// ==========================

	fileServer := http.FileServer(http.Dir("./web"))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// API-Pfade nicht vom Static Server behandeln
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		// Statische Dateien ausliefern (index.html wird automatisch bei "/" geladen)
		fileServer.ServeHTTP(w, r)
	}))

	return mux
}
