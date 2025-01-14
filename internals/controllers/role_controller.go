package controllers

import (
	"car-bond/internals/models/userRegistration"

	"car-bond/internals/database"

	"github.com/gofiber/fiber/v2"
)

// CreateGroup creates a new group
func CreateGroup(c *fiber.Ctx) error {
	db := database.DB.Db
	var group userRegistration.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Create(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create group"})
	}
	return c.Status(fiber.StatusCreated).JSON(group)
}

// GetGroup retrieves a group by its code
func GetGroup(c *fiber.Ctx) error {
	db := database.DB.Db
	code := c.Params("code")
	var group userRegistration.Group
	if err := db.Where("code = ?", code).First(&group).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}
	return c.JSON(group)
}

// GetGroups retrieves all groups
func GetAllGroups(c *fiber.Ctx) error {
	db := database.DB.Db
	var groups []userRegistration.Group
	if err := db.Find(&groups).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve groups"})
	}
	return c.JSON(groups)
}

// UpdateGroup updates a group by code
func UpdateGroup(c *fiber.Ctx) error {
	db := database.DB.Db
	code := c.Params("code")
	var group userRegistration.Group
	if err := db.First(&group, "code = ?", code).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Save(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update group"})
	}
	return c.JSON(group)
}

// DeleteGroup deletes a group by code
func DeleteGroup(c *fiber.Ctx) error {
	db := database.DB.Db
	code := c.Params("code")
	if err := db.Delete(&userRegistration.Group{}, "code = ?", code).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot delete group"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// CreateRole creates a new role
func CreateRole(c *fiber.Ctx) error {
	db := database.DB.Db
	var role userRegistration.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Create(&role).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create role"})
	}
	return c.Status(fiber.StatusCreated).JSON(role)
}

// GetRole retrieves a role by its code
func GetRole(c *fiber.Ctx) error {
	db := database.DB.Db
	code := c.Params("code")
	var role userRegistration.Role
	if err := db.Where("code = ?", code).First(&role).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
	}
	return c.JSON(role)
}

// GetRoles retrieves all roles
func GetAllRoles(c *fiber.Ctx) error {
	db := database.DB.Db
	var roles []userRegistration.Role
	if err := db.Find(&roles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve roles"})
	}
	return c.JSON(roles)
}

// GetRolesForGroups retrieves roles assigned to groups
func GetRolesForGroups(c *fiber.Ctx) error {
	db := database.DB.Db
	groupCode := c.Query("group_code")
	var roles []userRegistration.Role
	if err := db.Raw(`
        SELECT roles.* 
        FROM roles
        JOIN group_roles ON group_roles.role_code = roles.code
        WHERE group_roles.group_code = ?
    `, groupCode).Scan(&roles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve roles for groups"})
	}
	return c.JSON(roles)
}

// CreateResource creates a new resource
func CreateResource(c *fiber.Ctx) error {
	db := database.DB.Db
	var resource userRegistration.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Create(&resource).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create resource"})
	}
	return c.Status(fiber.StatusCreated).JSON(resource)
}

// GetResource retrieves a resource by its code
func GetResource(c *fiber.Ctx) error {
	db := database.DB.Db
	code := c.Params("code")
	var resource userRegistration.Resource
	if err := db.Where("code = ?", code).First(&resource).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Resource not found"})
	}
	return c.JSON(resource)
}

// GetAllResources retrieves all non-deleted resources
func GetAllResources(c *fiber.Ctx) error {
	db := database.DB.Db
	var resources []userRegistration.Resource
	if err := db.Where("internal = ?", false).Find(&resources).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve resources"})
	}
	return c.JSON(resources)
}

// CreatePermission creates a new role-resource permission
func CreatePermission(c *fiber.Ctx) error {
	db := database.DB.Db
	var permission userRegistration.RoleResourcePermission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Create(&permission).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create permission"})
	}
	return c.Status(fiber.StatusCreated).JSON(permission)
}

// CreateWildCardPermission creates a new wildcard permission
func CreateWildCardPermission(c *fiber.Ctx) error {
	db := database.DB.Db
	var permission userRegistration.RoleWildCardPermission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Create(&permission).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create wildcard permission"})
	}
	return c.Status(fiber.StatusCreated).JSON(permission)
}

// ===============================

// GetGrantedPermissions calculates permissions for roles on a resource
func GetGrantedPermissions(c *fiber.Ctx) error {
	db := database.DB.Db
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")

	var permissions userRegistration.Permissions
	if err := db.Raw(`
        SELECT 
            COALESCE(MAX(permissions.allow_r), false) AS r,
            COALESCE(MAX(permissions.allow_w), false) AS w,
            COALESCE(MAX(permissions.allow_x), false) AS x,
            COALESCE(MAX(permissions.allow_d), false) AS d
        FROM role_resource_permissions
        WHERE role_code = ? AND resource_code = ?
    `, roleCode, resourceCode).Scan(&permissions.Allow).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot calculate permissions"})
	}

	return c.JSON(permissions)
}

// CheckPermissions checks if the role has the requested permissions on a resource
func CheckPermissions(c *fiber.Ctx) error {
	db := database.DB.Db
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")
	requestedPerms := userRegistration.Permissions{}
	if err := c.BodyParser(&requestedPerms); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse requested permissions"})
	}

	var permissions userRegistration.Permissions
	if err := db.Raw(`
        SELECT 
            COALESCE(MAX(permissions.allow_r), false) AS r,
            COALESCE(MAX(permissions.allow_w), false) AS w,
            COALESCE(MAX(permissions.allow_x), false) AS x,
            COALESCE(MAX(permissions.allow_d), false) AS d
        FROM role_resource_permissions
        WHERE role_code = ? AND resource_code = ?
    `, roleCode, resourceCode).Scan(&permissions.Allow).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot check permissions"})
	}

	// Compare requested permissions with allowed permissions
	if requestedPerms.Allow.R && !permissions.Allow.R ||
		requestedPerms.Allow.W && !permissions.Allow.W ||
		requestedPerms.Allow.X && !permissions.Allow.X ||
		requestedPerms.Allow.D && !permissions.Allow.D {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetExplicitPermissions retrieves explicitly defined permissions for a role on a resource
func GetExplicitPermissions(c *fiber.Ctx) error {
	db := database.DB.Db
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")

	var permissions []userRegistration.RoleResourcePermission
	if err := db.Where("role_code = ? AND resource_code = ?", roleCode, resourceCode).Find(&permissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve explicit permissions"})
	}

	return c.JSON(permissions)
}

// GetWildCardPermissions retrieves wildcard permissions for roles
func GetWildCardPermissions(c *fiber.Ctx) error {
	db := database.DB.Db
	roleCode := c.Query("role_code")
	resourcePattern := c.Query("resource_pattern")

	var permissions []userRegistration.RoleWildCardPermission
	if err := db.Where("role_code = ? AND resource_pattern LIKE ?", roleCode, resourcePattern).Find(&permissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve wildcard permissions"})
	}

	return c.JSON(permissions)
}

// ResourceExplicitPermissionsExists checks if there are explicit (non-wildcard) permissions on the resource
func ResourceExplicitPermissionsExists(c *fiber.Ctx) error {
	db := database.DB.Db
	resourceCode := c.Query("resource_code")

	var count int64
	err := db.Model(&userRegistration.RoleResourcePermission{}).
		Where("resource_code = ?", resourceCode).
		Count(&count).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking explicit permissions"})
	}

	exists := count > 0
	return c.JSON(fiber.Map{"exists": exists})
}

// GetRoleCodesForGroups retrieves role codes assigned to groups
func GetRoleCodesForGroups(c *fiber.Ctx) error {
	db := database.DB.Db
	groupCodes := c.Query("group_codes") // Comma-separated group codes

	var roleCodes []string
	err := db.Raw(`
        SELECT DISTINCT role_code
        FROM group_roles
        WHERE group_code IN (?)
    `, groupCodes).Scan(&roleCodes).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving role codes for groups"})
	}

	return c.JSON(roleCodes)
}

// GroupsWithRoleExists checks if there are groups assigned a specific role
func GroupsWithRoleExists(c *fiber.Ctx) error {
	db := database.DB.Db
	roleCode := c.Query("role_code")

	var count int64
	err := db.Table("group_roles").
		Where("role_code = ?", roleCode).
		Count(&count).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking if groups with role exist"})
	}

	exists := count > 0
	return c.JSON(fiber.Map{"exists": exists})
}

// GetPermissions retrieves permissions granted to roles on a specific resource
func GetPermissions(c *fiber.Ctx) error {
	db := database.DB.Db
	roleCodes := c.Query("role_codes") // Comma-separated role codes
	resourceCode := c.Query("resource_code")

	var permissions []userRegistration.RoleResourcePermission
	err := db.Where("role_code IN (?) AND resource_code = ?", roleCodes, resourceCode).
		Find(&permissions).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving permissions"})
	}

	return c.JSON(permissions)
}

// =============================

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

// ==================================

// SecurityServiceWrapper is a wrapper for the concrete implementation of userRegistration.SecurityService
type SecurityServiceWrapper struct {
	Service SecurityService
}

// ConcreteSecurityStorage represents the actual implementation of userRegistration.SecurityStorage
type ConcreteSecurityStorage struct {
	Storage SecurityStorage
}

// StorageInterface defines methods for interacting with storage
type StorageInterface interface {
	GetPermissions(ctx *fiber.Ctx, resource string, roles []string) ([]userRegistration.Permissions, error)
	GetWildcardPermissions(ctx *fiber.Ctx, resource string, roles []string) ([]userRegistration.Permissions, error)
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
