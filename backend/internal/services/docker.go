// =====================================================
// File: backend/internal/services/docker.go
// Purpose: Detect Docker based services
// =====================================================

package services

import (
	"os/exec"
	"strings"
)

// -----------------------------------------------------
// DetectDockerServices
//
// Detects relevant Docker containers and maps them
// to RootGuard service models.
//
// Strategy:
//
// 1. Query running containers via docker ps
// 2. Filter infrastructure containers
// 3. Detect known DNS services
//
// -----------------------------------------------------

func DetectDockerServices() []Service {

	var services []Service

	cmd := exec.Command(
		"docker",
		"ps",
		"--format",
		"{{.Names}}",
	)

	out, err := cmd.Output()
	if err != nil {
		return services
	}

	lines := strings.Split(string(out), "\n")

	for _, name := range lines {

		name = strings.TrimSpace(name)

		if name == "" {
			continue
		}

		lower := strings.ToLower(name)

		// -------------------------------------------------
		// Ignore infrastructure containers
		// -------------------------------------------------

		if strings.Contains(lower, "buildkit") ||
			strings.Contains(lower, "builder") ||
			strings.Contains(lower, "watchtower") ||
			strings.Contains(lower, "portainer") ||
			strings.Contains(lower, "traefik") ||
			strings.Contains(lower, "zoraxy") {

			continue
		}

		// -------------------------------------------------
		// Detect DNS services
		// -------------------------------------------------

		if strings.Contains(lower, "adguard") {

			services = append(services, Service{
				Name:      "AdGuard Home",
				Type:      "docker",
				Status:    "running",
				Container: name,
			})

			continue
		}

		if strings.Contains(lower, "unbound") {

			services = append(services, Service{
				Name:      "Unbound",
				Type:      "docker",
				Status:    "running",
				Container: name,
			})

			continue
		}

		if strings.Contains(lower, "pihole") {

			services = append(services, Service{
				Name:      "Pi-hole",
				Type:      "docker",
				Status:    "running",
				Container: name,
			})

			continue
		}

		if strings.Contains(lower, "coredns") {

			services = append(services, Service{
				Name:      "CoreDNS",
				Type:      "docker",
				Status:    "running",
				Container: name,
			})

			continue
		}

		if strings.Contains(lower, "bind") {

			services = append(services, Service{
				Name:      "Bind9",
				Type:      "docker",
				Status:    "running",
				Container: name,
			})

			continue
		}
	}

	return services
}
