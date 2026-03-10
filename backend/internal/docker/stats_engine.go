// =====================================================
// File: backend/internal/docker/stats_engine.go
// Project: RootGuard WebApp
// Purpose: Docker statistics calculation helpers
//
// Responsibilities:
//
// - calculate CPU usage from Docker stats snapshots
// - normalize Docker stats into RootGuard metrics
//
// Notes:
//
// Docker CPU usage cannot be calculated from a single
// snapshot. It requires two consecutive stats samples.
// The formula implemented here matches the Docker CLI.
//
// =====================================================

package docker

import (
	"github.com/docker/docker/api/types"
)

// =====================================================
// calculateCPU()
// =====================================================
//
// Calculates container CPU usage using two Docker stats
// snapshots.
//
// Formula (same as docker stats CLI):
//
// cpuDelta = current.TotalUsage - previous.TotalUsage
// systemDelta = current.SystemUsage - previous.SystemUsage
//
// cpu% = (cpuDelta / systemDelta) * number_of_CPUs * 100
//
// =====================================================

func calculateCPU(previous types.StatsJSON, current types.StatsJSON) float64 {

	cpuDelta := float64(current.CPUStats.CPUUsage.TotalUsage) -
		float64(previous.CPUStats.CPUUsage.TotalUsage)

	systemDelta := float64(current.CPUStats.SystemUsage) -
		float64(previous.CPUStats.SystemUsage)

	if systemDelta <= 0 || cpuDelta <= 0 {
		return 0
	}

	cpuCount := float64(len(current.CPUStats.CPUUsage.PercpuUsage))

	if cpuCount == 0 {
		cpuCount = 1
	}

	cpuPercent := (cpuDelta / systemDelta) * cpuCount * 100

	return cpuPercent
}

// =====================================================
// calculateMemoryMB()
// =====================================================
//
// Converts Docker memory stats into megabytes.
//
// Docker reports memory usage in bytes.
// RootGuard dashboard uses MB for readability.
//
// =====================================================

func calculateMemoryMB(stats types.StatsJSON) float64 {

	return float64(stats.MemoryStats.Usage) / (1024 * 1024)
}
