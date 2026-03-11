// =====================================================
// File: backend/internal/services/registry.go
// Project: RootGuard WebApp
//
// Purpose:
//
// Central service registry for RootGuard.
//
// The registry defines which logical services exist
// and how they map to Docker containers.
//
// This allows the API to work with stable service
// identifiers instead of Docker container names.
//
// Example:
//
//   /api/service/unbound/restart
//
// The registry maps:
//
//   service "unbound"
//        -> container "rootguard-unbound"
//
// =====================================================

package services

// =====================================================
// Registry
//
// List of all services managed by RootGuard.
//
// This acts as the single source of truth for:
//
//   - API layer
//   - Service control
//   - Dashboard UI
//
// =====================================================

var Registry = []Service{

	{
		Name:        "unbound",
		DisplayName: "Unbound DNS",
		Description: "Recursive DNS resolver with DNSSEC validation",
		Icon:        "shield",

		Containers: []string{
			"rootguard-unbound",
		},
	},

	{
		Name:        "adguard",
		DisplayName: "AdGuard Home",
		Description: "DNS filtering and query logging",
		Icon:        "filter",

		Containers: []string{
			"rootguard-adguard",
		},
	},
}

// =====================================================
// GetService
//
// Returns a service by its name.
// =====================================================

func GetService(name string) *Service {

	for _, service := range Registry {

		if service.Name == name {
			return &service
		}

	}

	return nil
}

// =====================================================
// ListServices
//
// Returns the complete service registry.
//
// Used by:
//
//   GET /api/services
//
// =====================================================

func ListServices() []Service {

	return Registry

}
