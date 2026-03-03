// =====================================================
// File: backend/internal/httpapi/handlers.go
// Purpose: HTTP Handlers
// =====================================================

package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

// -----------------------------------------------------
// Dashboard Mock
// -----------------------------------------------------

func HandleDashboard(w http.ResponseWriter, r *http.Request) {

	resp := map[string]interface{}{
		"docker": map[string]float64{
			"cpu":    12.5,
			"memory": 37.8,
		},
		"dns": map[string]int{
			"queries": 2345,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// -----------------------------------------------------
// Service Action Handler
// Path: /api/service/{name}/{action}
// -----------------------------------------------------

func handleServiceAction(w http.ResponseWriter, r *http.Request) {

	// Pfad zerlegen
	// Beispiel: /api/service/adguard/restart
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 5 {
		http.Error(w, "Invalid service path", http.StatusBadRequest)
		return
	}

	name := parts[3]
	action := parts[4]

	response := map[string]string{
		"service": name,
		"action":  action,
		"status":  "accepted",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
