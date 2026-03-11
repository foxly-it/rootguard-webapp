// =====================================================
// File: backend/internal/services/models.go
// Project: RootGuard WebApp
//
// Purpose:
//
// Defines the logical RootGuard service model.
// =====================================================

package services

// =====================================================
// Service
//
// Represents a detected or configured service.
// =====================================================

type Service struct {

	// -------------------------------------------------
	// Internal name or identifier
	// -------------------------------------------------
	Name string `json:"name"`

	// -------------------------------------------------
	// Optional UI display name
	// -------------------------------------------------
	DisplayName string `json:"displayName,omitempty"`

	// -------------------------------------------------
	// Optional description
	// -------------------------------------------------
	Description string `json:"description,omitempty"`

	// -------------------------------------------------
	// UI icon (lucide / frontend)
	// -------------------------------------------------
	Icon string `json:"icon,omitempty"`

	// -------------------------------------------------
	// Docker containers belonging to the service
	// -------------------------------------------------
	Containers []string `json:"containers,omitempty"`

	// -------------------------------------------------
	// Service type
	//
	// Examples:
	//   binary
	//   docker
	//   system
	// -------------------------------------------------
	Type string `json:"type,omitempty"`

	// -------------------------------------------------
	// Runtime status
	//
	// running / stopped / unknown
	// -------------------------------------------------
	Status string `json:"status,omitempty"`

	// -------------------------------------------------
	// Docker container name (legacy detection)
	// -------------------------------------------------
	Container string `json:"container,omitempty"`
}

// =====================================================
// ServiceStatus
//
// Used by API responses.
// =====================================================

type ServiceStatus struct {
	Name string `json:"name"`

	DisplayName string `json:"displayName"`

	Status string `json:"status"`
}
