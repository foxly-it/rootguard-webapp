package httpapi

import (
	"encoding/json"
	"net/http"
)

var (
	Version   = "dev"
	GitCommit = "none"
)

type VersionResponse struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {

	response := VersionResponse{
		Version:   Version,
		GitCommit: GitCommit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
