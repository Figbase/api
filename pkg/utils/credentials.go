package utils

import (
	"fmt"

	"github.com/Figbase/api/pkg/repository"
)

// GetCredentialsByRole func for getting credentials from a role name.
func GetCredentialsByRole(role string) ([]string, error) {
	// Define credentials variable.
	var credentials []string

	// Switch given role.
	switch role {
	case repository.AdminRoleName:
		// Admin credentials (all access).
		credentials = []string{
			repository.AppCreateCredential,
			repository.AppUpdateCredential,
			repository.AppDeleteCredential,
		}
	case repository.ModeratorRoleName:
		// Moderator credentials (only some access).
		credentials = []string{
			repository.AppCreateCredential,
			repository.AppUpdateCredential,
		}
	case repository.UserRoleName:
		// Simple user credentials (less acess).
		credentials = []string{
			repository.AppCreateCredential,
		}
	default:
		// Return error message.
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}

	return credentials, nil
}
