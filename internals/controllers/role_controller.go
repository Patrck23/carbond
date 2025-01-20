package controllers

import (
	"car-bond/internals/models/userRegistration"
	"strings"

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

// GetRolesForGroups retrieves roles assigned to multiple groups.
func GetRolesForGroups(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

	// Retrieve the group codes from the query parameters
	groupCodes := c.Query("group_codes")
	if groupCodes == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group codes are required",
		})
	}

	// Split the group codes into a slice
	groupCodeList := strings.Split(groupCodes, ",")

	// Define a slice to hold the roles
	var roles []userRegistration.Role

	// Query the database to fetch roles associated with the groups
	query := `
        SELECT DISTINCT roles.* 
        FROM roles
        JOIN group_roles ON group_roles.role_code = roles.code
        WHERE group_roles.group_code IN ?
    `
	if err := db.Raw(query, groupCodeList).Scan(&roles).Error; err != nil {
		// Return an internal server error if the query fails
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve roles for the specified groups",
		})
	}

	// Return the retrieved roles as JSON
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

// GetGrantedPermissions calculates permissions for roles on a resource
// GetGrantedPermissions calculates permissions for roles on a resource
func GetGrantedPermissions(c *fiber.Ctx) error {
	db := database.DB.Db
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")

	// Validate input parameters
	if roleCode == "" || resourceCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both role_code and resource_code are required",
		})
	}

	// Struct to hold the permissions
	var permissions userRegistration.Permissions

	// Query to calculate the permissions
	if err := db.Raw(`
        SELECT 
            COALESCE(MAX(role_resource_permissions.allow_r), false) AS r,
            COALESCE(MAX(role_resource_permissions.allow_w), false) AS w,
            COALESCE(MAX(role_resource_permissions.allow_x), false) AS x,
            COALESCE(MAX(role_resource_permissions.allow_d), false) AS d
        FROM role_resource_permissions
        WHERE role_code = ? AND resource_code = ?
    `, roleCode, resourceCode).Scan(&permissions).Error; err != nil {
		// Return an internal server error if the query fails
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve permissions",
		})
	}

	// Return the permissions as JSON
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
