// =====================================================
// File: backend/internal/api/service.go
// Project: RootGuard WebApp
// Purpose: Handle service control actions
//
// Endpoint:
//
// POST /api/service/{name}/{action}
//
// Examples:
//
// POST /api/service/adguard/start
// POST /api/service/unbound/restart
//
// =====================================================

package api

import (
	"net/http"
	"strings"

	"github.com/foxly-it/rootguard-webapp/backend/internal/docker"
)

// -----------------------------------------------------
// HandleServiceAction
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

	var err error

	switch action {

	case "start":

		err = docker.StartContainer(serviceName)

	case "stop":

		err = docker.StopContainer(serviceName)

	case "restart":

		err = docker.RestartContainer(serviceName)

	default:

		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
