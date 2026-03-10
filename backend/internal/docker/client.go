// =====================================================
// File: backend/internal/docker/client.go
// Purpose: Create Docker API client
// =====================================================

package docker

import (
	"github.com/docker/docker/client"
)

// -----------------------------------------------------
// NewClient
// Creates Docker client using environment configuration
// -----------------------------------------------------

func NewClient() (*client.Client, error) {

	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
}
