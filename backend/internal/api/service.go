// =====================================================
// File: backend/internal/api/service.go
// Purpose: Handle service control actions
// =====================================================

package api

import (
	"net/http"
	"strings"
)

// -----------------------------------------------------
// HandleServiceAction
// Example endpoint:
//
// POST /api/service/adguard/start
// POST /api/service/unbound/restart
// -----------------------------------------------------

func HandleServiceAction(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/api/service/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 {

		http.Error(w, "Invalid service path", http.StatusBadRequest)
		return
	}

	serviceName := parts[0]
	action := parts[1]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("service=" + serviceName + " action=" + action))
}
