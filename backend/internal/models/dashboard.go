// =====================================================
// File: backend/internal/models/dashboard.go
// Purpose: Dashboard response models
// =====================================================

package models

type ContainerStats struct {
	Name   string  `json:"name"`
	Image  string  `json:"image"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Uptime string  `json:"uptime"`
}

type DashboardResponse struct {
	Containers []ContainerStats `json:"containers"`
	Total      int              `json:"total"`
	Running    int              `json:"running"`
}
