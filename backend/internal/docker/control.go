// =====================================================
// File: backend/internal/docker/control.go
// Project: RootGuard WebApp
// Purpose: Docker container lifecycle operations
// =====================================================

package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

// -----------------------------------------------------
// StartContainer
// -----------------------------------------------------

func StartContainer(name string) error {

	cli, err := NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	id, err := FindContainer(name)
	if err != nil {
		return err
	}

	if id == "" {
		return nil
	}

	return cli.ContainerStart(ctx, id, container.StartOptions{})
}

// -----------------------------------------------------
// StopContainer
// -----------------------------------------------------

func StopContainer(name string) error {

	cli, err := NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	id, err := FindContainer(name)
	if err != nil {
		return err
	}

	if id == "" {
		return nil
	}

	timeout := 10

	return cli.ContainerStop(ctx, id, container.StopOptions{
		Timeout: &timeout,
	})
}

// -----------------------------------------------------
// RestartContainer
// -----------------------------------------------------

func RestartContainer(name string) error {

	cli, err := NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	id, err := FindContainer(name)
	if err != nil {
		return err
	}

	if id == "" {
		return nil
	}

	timeout := 10

	return cli.ContainerRestart(ctx, id, container.StopOptions{
		Timeout: &timeout,
	})
}
