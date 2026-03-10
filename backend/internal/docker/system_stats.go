// =====================================================
// File: backend/internal/docker/system_stats.go
// Project: RootGuard WebApp
// Purpose: Collect Docker system overview metrics
//
// Responsibilities:
//
// - fetch Docker version
// - count containers
// - count running containers
// - aggregate CPU usage
// - aggregate memory usage
// - normalize values for dashboard display
//
// =====================================================

package docker

import (
	"context"
	"math"

	"github.com/docker/docker/api/types/container"

	"github.com/foxly-it/rootguard-webapp/backend/internal/models"
)

// =====================================================
// GetSystemStats
//
// Aggregates Docker metrics for the RootGuard dashboard
// overview cards.
//
// =====================================================

func GetSystemStats() (*models.SystemResponse, error) {

	ctx := context.Background()

	cli, err := NewClient()
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------
	// Docker Version
	// -------------------------------------------------

	versionInfo, err := cli.ServerVersion(ctx)
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------
	// Container List
	// -------------------------------------------------

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})

	if err != nil {
		return nil, err
	}

	total := len(containers)
	running := 0

	totalCPU := 0.0
	totalMemory := 0.0

	// -------------------------------------------------
	// Aggregate metrics from stats cache
	// -------------------------------------------------

	statsCache.RLock()

	for _, c := range containers {

		if c.State == "running" {
			running++
		}

		if metrics, ok := statsCache.data[c.ID]; ok {

			totalCPU += metrics.CPU
			totalMemory += metrics.Memory

		}
	}

	statsCache.RUnlock()

	// -------------------------------------------------
	// Normalize values for dashboard display
	// -------------------------------------------------

	totalCPU = math.Round(totalCPU*100) / 100
	totalMemory = math.Round(totalMemory*10) / 10

	return &models.SystemResponse{
		DockerVersion: versionInfo.Version,
		Containers:    total,
		Running:       running,
		CPU:           totalCPU,
		Memory:        totalMemory,
	}, nil
}
