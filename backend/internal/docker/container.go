// =====================================================
// File: backend/internal/docker/container.go
// Project: RootGuard WebApp
// Purpose: Docker container discovery helpers
//
// Responsibilities:
//
// - find containers by name
// - list containers
//
// This file contains NO lifecycle actions.
// =====================================================

package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// -----------------------------------------------------
// FindContainer
//
// Returns container ID by name.
//
// Example:
//
// name: "adguard"
// matches container:
// "/adguard"
// -----------------------------------------------------

func FindContainer(name string) (string, error) {

	cli, err := NewClient()
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})

	if err != nil {
		return "", err
	}

	for _, c := range containers {

		for _, n := range c.Names {

			trim := strings.TrimPrefix(n, "/")

			if trim == name {
				return c.ID, nil
			}
		}
	}

	return "", nil
}
