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

func HandlePreviewUnboundSettings(w http.ResponseWriter, r *http.Request, core *coreclient.Client) {
	settings, ok := decodeUnboundSettings(w, r)
	if !ok {
		return
	}
	preview, err := core.PreviewUnboundSettings(r.Context(), settings)
	if err != nil {
		writeCoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, preview)
}

func HandleUnboundHistory(w http.ResponseWriter, r *http.Request, core *coreclient.Client) {
	history, err := core.UnboundHistory(r.Context())
	if err != nil {
		writeCoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, history)
}

func HandleRestoreUnboundVersion(w http.ResponseWriter, r *http.Request, core *coreclient.Client) {
	settings, err := core.RestoreUnboundVersion(r.Context(), r.PathValue("id"))
	if err != nil {
		writeCoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func HandleUnboundDiagnostics(w http.ResponseWriter, r *http.Request, core *coreclient.Client) {
	report, err := core.UnboundDiagnostics(r.Context())
	if err != nil {
		writeCoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func decodeUnboundSettings(w http.ResponseWriter, r *http.Request) (coreclient.UnboundSettings, bool) {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	var settings coreclient.UnboundSettings
	if err := decoder.Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return coreclient.UnboundSettings{}, false
	}
	return settings, true
}

func writeCoreError(w http.ResponseWriter, err error) {
	status := http.StatusBadGateway
	var apiError *coreclient.APIError
	if errors.As(err, &apiError) && apiError.StatusCode >= 400 && apiError.StatusCode < 500 {
		status = apiError.StatusCode
	}
	http.Error(w, err.Error(), status)
}
