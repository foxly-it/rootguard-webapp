// =====================================================
// File: backend/internal/api/system.go
// Project: RootGuard WebApp
// Purpose: System overview API handler
//
// Endpoint:
//
// GET /api/system
//
// =====================================================

package api

import (
	"encoding/json"
	"net/http"

	"github.com/foxly-it/rootguard-webapp/backend/internal/docker"
)

// =====================================================
// HandleSystem
//
// Returns aggregated Docker system stats.
// =====================================================

func HandleSystem(w http.ResponseWriter, r *http.Request) {

	stats, err := docker.GetSystemStats()
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stats)
}
