// =====================================================
// File: backend/internal/api/services.go
// Project: RootGuard WebApp
// Purpose: Return detected services as JSON
// =====================================================

package api

import (
	"encoding/json"
	"net/http"

	"github.com/foxly-it/rootguard-webapp/backend/internal/services"
)

// -----------------------------------------------------
// HandleServices
//
// GET /api/services
//
// Returns all detected services from the service
// detection engine.
// -----------------------------------------------------

func HandleServices(w http.ResponseWriter, r *http.Request) {

	// Only allow GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Run detection engine
	serviceList := services.DetectServices()

	// JSON response
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(serviceList)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
