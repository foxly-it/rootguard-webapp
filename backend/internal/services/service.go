// =====================================================
// File: backend/internal/services/service.go
// Project: RootGuard WebApp
//
// Purpose:
//
// Defines the logical RootGuard service model.
//
// A service represents a logical unit in the system,
// such as:
//
//   - Unbound DNS
//   - AdGuard Home
//   - DNS Stack
//
// A service can control one or multiple Docker containers.
//
// This abstraction decouples:
//
//   API layer  -> logical service
//   Docker     -> container names
//
// =====================================================

package services

// =====================================================
// Service
//
// Logical representation of a manageable RootGuard
// service.
// =====================================================

type Service struct {

	// -------------------------------------------------
	// Name
	//
	// Internal service identifier used by the API.
	//
	// Example:
	//
	//   /api/service/unbound/restart
	//
	// -------------------------------------------------
	Name string

	// -------------------------------------------------
	// DisplayName
	//
	// Human readable service name for UI.
	// -------------------------------------------------
	DisplayName string

	// -------------------------------------------------
	// Description
	//
	// Optional description shown in the UI.
	// -------------------------------------------------
	Description string

	// -------------------------------------------------
	// Icon
	//
	// Icon identifier used by the frontend.
	// Example (lucide-react):
	//
	//   shield
	//   server
	//   network
	//
	// -------------------------------------------------
	Icon string

	// -------------------------------------------------
	// Containers
	//
	// Docker containers belonging to this service.
	//
	// This allows services to control multiple
	// containers (e.g. a stack).
	//
	// Example:
	//
	//   []string{"rootguard-unbound"}
	//
	// -------------------------------------------------
	Containers []string
}

// =====================================================
// ServiceStatus
//
// Runtime status used by the API response.
// =====================================================

type ServiceStatus struct {

	// logical service name
	Name string `json:"name"`

	// UI display name
	DisplayName string `json:"displayName"`

	// running / stopped / unknown
	Status string `json:"status"`
}
