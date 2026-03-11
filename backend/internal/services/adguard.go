// =====================================================
// File: backend/internal/services/adguard.go
// Project: RootGuard WebApp
// Purpose: Detect AdGuard Home process
// =====================================================

package services

import (
	"os/exec"
)

// -----------------------------------------------------
// DetectAdGuard
//
// Detection Strategy:
//
// 1. Verify that the AdGuardHome binary exists
// 2. Use pgrep to check if the process is running
// -----------------------------------------------------

func DetectAdGuard() *Service {

	// -------------------------------------------------
	// Check if binary exists
	// -------------------------------------------------

	_, err := exec.LookPath("AdGuardHome")
	if err != nil {
		return nil
	}

	// -------------------------------------------------
	// Check if process is running
	// -------------------------------------------------

	cmd := exec.Command("pgrep", "-f", "AdGuardHome")

	err = cmd.Run()

	if err != nil {

		return &Service{
			Name:   "AdGuard Home",
			Type:   "binary",
			Status: "stopped",
		}
	}

	return &Service{
		Name:   "AdGuard Home",
		Type:   "binary",
		Status: "running",
	}
}
