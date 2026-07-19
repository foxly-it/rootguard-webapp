package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/foxly-it/rootguard-webapp/backend/internal/coreclient"
)

func HandleGetUnboundSettings(w http.ResponseWriter, r *http.Request, core *coreclient.Client) {
	settings, err := core.UnboundSettings(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func HandlePutUnboundSettings(w http.ResponseWriter, r *http.Request, core *coreclient.Client) {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	var settings coreclient.UnboundSettings
	if err := decoder.Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := core.UpdateUnboundSettings(r.Context(), settings)
	if err != nil {
		status := http.StatusBadGateway
		var apiError *coreclient.APIError
		if errors.As(err, &apiError) && apiError.StatusCode >= 400 && apiError.StatusCode < 500 {
			status = apiError.StatusCode
		}
		http.Error(w, err.Error(), status)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
