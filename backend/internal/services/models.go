// =====================================================
// File: backend/internal/services/models.go
// Purpose: Shared service models
// =====================================================

package services

// -----------------------------------------------------
// Service
// Represents a detected service on the host
// -----------------------------------------------------

type Service struct {
	Name string `json:"name"`

	Type string `json:"type"`
	// docker
	// systemd
	// binary

	Status string `json:"status"`
	// running
	// stopped
	// unknown

	Version string `json:"version,omitempty"`

	Container string `json:"container,omitempty"`
}
