// =====================================================
// File: backend/internal/models/system.go
// Project: RootGuard WebApp
// Purpose: Data models for system overview API
//
// Used by:
// - /api/system endpoint
//
// =====================================================

package models

// =====================================================
// SystemResponse
//
// Contains high level Docker system information used
// by the RootGuard dashboard overview cards.
// =====================================================

type SystemResponse struct {
	DockerVersion string  `json:"docker_version"`
	Containers    int     `json:"containers"`
	Running       int     `json:"running"`
	CPU           float64 `json:"cpu_total"`
	Memory        float64 `json:"memory_total"`
}
