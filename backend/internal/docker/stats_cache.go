// =====================================================
// File: backend/internal/docker/stats_cache.go
// Project: RootGuard WebApp
// Purpose: In-memory metrics cache for container stats
//
// Responsibilities:
//
// - store container metrics
// - provide thread-safe read/write access
//
// =====================================================

package docker

import "sync"

// =====================================================
// ContainerMetrics
//
// Holds normalized metrics for one container.
// =====================================================

type ContainerMetrics struct {
	CPU     float64
	Memory  float64
	LastCPU float64
}

// =====================================================
// statsCache
//
// Thread-safe in-memory cache used by:
//
// - Stats Collector (writer)
// - Dashboard API (reader)
//
// =====================================================

var statsCache = struct {
	sync.RWMutex
	data map[string]ContainerMetrics
}{
	data: make(map[string]ContainerMetrics),
}
