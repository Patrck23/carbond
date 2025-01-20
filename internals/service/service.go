package service

import (
	"car-bond/internals/models/userRegistration"

	"github.com/gofiber/fiber/v2"
)

type SecurityService interface {
	// GetGroup retrieves a group by code
	GetGroup(c *fiber.Ctx, code string) (*userRegistration.Group, error)
	// GetAllGroups retrieves all not deleted groups
	GetAllGroups(c *fiber.Ctx) ([]*userRegistration.Group, error)
	// GetRole retrieves a role by code
	GetRole(c *fiber.Ctx, code string) (*userRegistration.Role, error)
	// GetAllRoles retrieves all not deleted roles
	GetAllRoles(c *fiber.Ctx) ([]*userRegistration.Role, error)
	// GetRolesForGroups retrieves roles assigned on groups
	GetRolesForGroups(c *fiber.Ctx, groups []string) ([]string, error)
	// GetResource retrieves a resource by code
	GetResource(c *fiber.Ctx, code string) (*userRegistration.Resource, error)
	// GetAllResources retrieves all not deleted resources
	GetAllResources(c *fiber.Ctx) ([]*userRegistration.Resource, error)
	// GetGrantedPermissions calculates permissions on the resource for the roles and applies allow/deny logic
	GetGrantedPermissions(c *fiber.Ctx, resource string, roles []string) (*userRegistration.RWXD, error)
	// CheckPermissions checks if the roles have the requested perms on the given resource
	CheckPermissions(c *fiber.Ctx, resource string, roles []string, requestedPermissions []string) (bool, error)
	// GetExplicitPermissions returns permissions on resource / roles setup explicitly
	GetExplicitPermissions(c *fiber.Ctx, resources []string, roles []string) ([]*userRegistration.RoleResourcePermission, error)
	// GetWildCardPermissions returns wildcard permissions on roles
	GetWildCardPermissions(c *fiber.Ctx, roles []string) ([]*userRegistration.RoleWildCardPermission, error)
}

type SecurityStorage interface {
	// GetGroup retrieves a group by code
	GetGroup(c *fiber.Ctx, code string) (*userRegistration.Group, error)
	// GetGroups retrieves all not deleted groups
	GetGroups(c *fiber.Ctx) ([]*userRegistration.Group, error)
	// GetRole retrieves a role by code
	GetRole(c *fiber.Ctx, code string) (*userRegistration.Role, error)
	// GetAllRoles retrieves all not deleted roles
	GetAllRoles(c *fiber.Ctx) ([]*userRegistration.Role, error)
	// GetAllRoleCodes retrieves all role codes
	GetAllRoleCodes(c *fiber.Ctx) ([]string, error)
	// GetResource retrieves a resource by code
	GetResource(c *fiber.Ctx, code string) (*userRegistration.Resource, error)
	// GetAllResources retrieves all not deleted resources
	GetAllResources(c *fiber.Ctx) ([]*userRegistration.Resource, error)
	// ResourceExplicitPermissionsExists checks if there are explicit (no wildcard) permissions on the resource
	ResourceExplicitPermissionsExists(c *fiber.Ctx, code string) (bool, error)
	// GetRoleCodesForGroups retrieves role codes for groups
	GetRoleCodesForGroups(c *fiber.Ctx, groups []string) ([]string, error)
	// GroupsWithRoleExists checks if there are groups with assigned role
	GroupsWithRoleExists(c *fiber.Ctx, role string) (bool, error)
	// GetPermissions retrieves permissions granted to roles on resource
	GetPermissions(c *fiber.Ctx, resource string, roles []string) ([]*userRegistration.Permissions, error)
	// GetWildcardPermissions retrieves wildcard permissions granted to roles on resource
	GetWildcardPermissions(c *fiber.Ctx, resource string, roles []string) ([]*userRegistration.Permissions, error)
}

// SecurityServiceWrapper is a wrapper for the concrete implementation of userRegistration.SecurityService
type SecurityServiceWrapper struct {
	Service SecurityService
}

// ConcreteSecurityStorage represents the actual implementation of userRegistration.SecurityStorage
type ConcreteSecurityStorage struct {
	Storage SecurityStorage
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
