// =====================================================
// File: backend/internal/services/detector.go
// Purpose: Detect running services
// =====================================================

package services

// -----------------------------------------------------
// DetectServices
// Main detection entrypoint
// -----------------------------------------------------

func DetectServices() []Service {

	services := []Service{}

	// AdGuard
	if svc := DetectAdGuard(); svc != nil {
		services = append(services, *svc)
	}

	// Unbound
	if svc := DetectUnbound(); svc != nil {
		services = append(services, *svc)
	}

	// Docker Services
	docker := DetectDockerServices()
	services = append(services, docker...)

	return services
}
