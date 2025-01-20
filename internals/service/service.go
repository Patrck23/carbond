package service

import (
	"car-bond/internals/models/userRegistration"

	"github.com/gofiber/fiber/v2"
)

// SecurityServiceWrapper is a wrapper for the concrete implementation of userRegistration.SecurityService
type SecurityServiceWrapper struct {
	Service userRegistration.SecurityService
}

// ConcreteSecurityStorage represents the actual implementation of userRegistration.SecurityStorage
type ConcreteSecurityStorage struct {
	Storage userRegistration.SecurityStorage
}

// GetGrantedPermissions retrieves and merges permissions
func (s *ConcreteSecurityStorage) GetGrantedPermissions(c *fiber.Ctx, resource string, roles []string) (*userRegistration.RWXD, error) {
	// Get explicit permissions
	explicitPermissions, err := s.Storage.GetPermissions(c, resource, roles)
	if err != nil {
		return nil, err
	}

	// Get wildcard permissions
	wildCardPermissions, err := s.Storage.GetWildcardPermissions(c, resource, roles)
	if err != nil {
		return nil, err
	}

	permissions := append(explicitPermissions, wildCardPermissions...)

	// Merge all roles' permissions (explicit and wildcard)
	resPermissions := &userRegistration.RWXD{}
	for _, p := range permissions {
		resPermissions.R = (resPermissions.R || p.Allow.R) && !p.Deny.R
		resPermissions.W = (resPermissions.W || p.Allow.W) && !p.Deny.W
		resPermissions.X = (resPermissions.X || p.Allow.X) && !p.Deny.X
		resPermissions.D = (resPermissions.D || p.Allow.D) && !p.Deny.D
	}

	return resPermissions, nil
}

// CheckPermissions checks if the given roles allow access to the requested resource
func (s *SecurityServiceWrapper) CheckPermissions(c *fiber.Ctx, resource string, roles []string, requestedPermissions []string) (bool, error) {
	// Empty request means no access
	if len(requestedPermissions) == 0 {
		return false, nil
	}

	// Get granted permissions
	grantedPerms, err := s.Service.GetGrantedPermissions(c, resource, roles)
	if err != nil {
		return false, err
	}

	// Check all requested permissions are granted
	for _, p := range requestedPermissions {
		if !isPermissionGranted(p, grantedPerms) {
			return false, nil
		}
	}

	return true, nil
}

// Helper function to check if a specific permission is granted
func isPermissionGranted(permission string, perms *userRegistration.RWXD) bool {
	switch permission {
	case "R":
		return perms.R
	case "W":
		return perms.W
	case "X":
		return perms.X
	case "D":
		return perms.D
	default:
		return false
	}
}
