// =====================================================
// File: backend/internal/docker/stats_collector.go
// Project: RootGuard WebApp
// Purpose: Background collector for Docker container metrics
//
// Responsibilities:
//
// - poll Docker stats periodically
// - calculate CPU usage
// - normalize memory usage
// - smooth CPU values
// - update metrics cache
//
// =====================================================

package docker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// =====================================================
// previousStats
//
// Stores previous Docker stats snapshot per container.
// Required for CPU delta calculation.
// =====================================================

var previousStats = make(map[string]types.StatsJSON)

// =====================================================
// StartStatsCollector
//
// Starts a background goroutine which periodically
// polls Docker stats and updates the metrics cache.
// =====================================================

func StartStatsCollector() {

	go func() {

		log.Println("Docker stats collector started")

		for {

			updateStats()

			time.Sleep(1 * time.Second)

		}

	}()
}

// =====================================================
// updateStats
//
// Fetch stats for all running containers
// and update the metrics cache.
// =====================================================

func updateStats() {

	ctx := context.Background()

	cli, err := NewClient()
	if err != nil {
		log.Println("docker client error:", err)
		return
	}

	// -------------------------------------------------
	// Fetch running containers
	// -------------------------------------------------

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: false,
	})

	if err != nil {
		log.Println("container list error:", err)
		return
	}

	// -------------------------------------------------
	// Iterate containers
	// -------------------------------------------------

	for _, c := range containers {

		stats, err := cli.ContainerStats(ctx, c.ID, false)
		if err != nil {
			continue
		}

		var current types.StatsJSON

		if err := json.NewDecoder(stats.Body).Decode(&current); err != nil {
			stats.Body.Close()
			continue
		}

		stats.Body.Close()

		prev, exists := previousStats[c.ID]

		var cpu float64

		if exists {

			cpu = calculateCPU(prev, current)

		} else {

			cpu = 0
		}

		mem := calculateMemoryMB(current)

		previousStats[c.ID] = current

		// -------------------------------------------------
		// CPU smoothing
		// -------------------------------------------------

		statsCache.Lock()

		old := statsCache.data[c.ID]

		smoothedCPU := (old.CPU + cpu) / 2

		statsCache.data[c.ID] = ContainerMetrics{
			CPU:     smoothedCPU,
			Memory:  mem,
			LastCPU: cpu,
		}

		statsCache.Unlock()
	}
}
