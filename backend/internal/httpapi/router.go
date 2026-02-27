package httpapi

import (
	"net/http"
)

func NewRouter() http.Handler {

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", HealthHandler)

	// API
	mux.HandleFunc("/api/version", VersionHandler)

	// CORS wrapper (for Vite dev)
	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
