// =====================================================
// File: backend/internal/services/unbound.go
// Project: RootGuard WebApp
// Purpose: Detect Unbound resolver process
// =====================================================

package services

import (
	"os/exec"
)

// -----------------------------------------------------
// DetectUnbound
//
// Detection Strategy:
//
// 1. Verify that the unbound binary exists
// 2. Use pgrep to check if the process is running
//
// Exit codes:
//
// pgrep exit 0 → process exists
// pgrep exit 1 → process not found
// -----------------------------------------------------

func DetectUnbound() *Service {

	// -------------------------------------------------
	// Check if binary exists
	// -------------------------------------------------

	_, err := exec.LookPath("unbound")
	if err != nil {
		return nil
	}

	// -------------------------------------------------
	// Check if process is running
	// -------------------------------------------------

	cmd := exec.Command("pgrep", "-x", "unbound")

	err = cmd.Run()

	if err != nil {

		return &Service{
			Name:   "Unbound",
			Type:   "binary",
			Status: "stopped",
		}
	}

	return &Service{
		Name:   "Unbound",
		Type:   "binary",
		Status: "running",
	}
}
