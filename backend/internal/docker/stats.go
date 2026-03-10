// =====================================================
// File: backend/internal/docker/stats.go
// Purpose: Collect container list for dashboard
// =====================================================

package docker

import (
	"context"
	"math"
	"strings"

	"github.com/docker/docker/api/types/container"

	"github.com/foxly-it/rootguard-webapp/backend/internal/models"
)

func GetDashboardStats() (*models.DashboardResponse, error) {

	ctx := context.Background()

	cli, err := NewClient()
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})

	if err != nil {
		return nil, err
	}

	resp := &models.DashboardResponse{}

	for _, c := range containers {

		status := "stopped"

		if c.State == "running" {
			status = "running"
			resp.Running++
		}

		resp.Total++

		name := ""

		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		cpu := 0.0
		mem := 0.0

		// ---------------------------------
		// Read metrics from cache
		// ---------------------------------

		statsCache.RLock()

		if m, ok := statsCache.data[c.ID]; ok {

			cpu = math.Round(m.CPU*100) / 100
			mem = math.Round(m.Memory*10) / 10
		}

		statsCache.RUnlock()

		containerStat := models.ContainerStats{
			Name:   name,
			Image:  c.Image,
			Status: status,
			CPU:    cpu,
			Memory: mem,
			Uptime: c.Status,
		}

		resp.Containers = append(resp.Containers, containerStat)
	}

	return resp, nil
}
