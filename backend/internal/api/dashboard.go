// =====================================================
// File: backend/internal/api/dashboard.go
// Purpose: Dashboard API endpoint
// =====================================================

package api

import (
	"encoding/json"
	"net/http"

	"github.com/foxly-it/rootguard-webapp/backend/internal/docker"
)

// -----------------------------------------------------
// HandleDashboard
// HTTP endpoint returning dashboard statistics
// -----------------------------------------------------

func HandleDashboard(w http.ResponseWriter, r *http.Request) {

	stats, err := docker.GetDashboardStats()
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stats)
}
