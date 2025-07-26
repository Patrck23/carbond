package controllers

import (
	"car-bond/internals/models/userRegistration"
	"car-bond/internals/repository"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GroupController handles group-related endpoints
type GroupController struct {
	service repository.GroupService
}

// NewGroupController creates a new instance of GroupController
func NewGroupController(service repository.GroupService) *GroupController {
	return &GroupController{service: service}
}

// CreateGroup handles the HTTP request to create a new group
func (gc *GroupController) CreateGroup(c *fiber.Ctx) error {
	var group userRegistration.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON payload",
			"data":    err.Error(),
		})
	}
	if err := gc.service.CreateGroup(&group); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create group",
			"data":    err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(group)
}

// GetGroup handles the HTTP request to retrieve a group by its code
func (gc *GroupController) GetGroup(c *fiber.Ctx) error {
	code := c.Params("code")
	group, err := gc.service.GetGroupByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}
	return c.JSON(group)
}

// GetAllGroups handles the HTTP request to retrieve all groups
func (gc *GroupController) GetAllGroups(c *fiber.Ctx) error {
	groups, err := gc.service.GetAllGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve groups"})
	}
	return c.JSON(groups)
}

// UpdateGroup handles the HTTP request to update an existing group
func (gc *GroupController) UpdateGroup(c *fiber.Ctx) error {
	code := c.Params("code")
	var group userRegistration.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := gc.service.UpdateGroup(code, &group); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found or cannot update"})
	}
	return c.JSON(fiber.Map{"success": "Group updated"})
}

// DeleteGroup handles the HTTP request to delete a group by its code
func (gc *GroupController) DeleteGroup(c *fiber.Ctx) error {
	code := c.Params("code")
	if err := gc.service.DeleteGroup(code); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot delete group"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ================================================================

// RoleController handles role-related endpoints
type RoleController struct {
	service repository.RoleService
}

// NewRoleController creates a new instance of RoleController
func NewRoleController(service repository.RoleService) *RoleController {
	return &RoleController{service: service}
}

// CreateRole handles the HTTP request to create a new role
func (rc *RoleController) CreateRole(c *fiber.Ctx) error {
	var role userRegistration.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := rc.service.CreateRole(&role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create role"})
	}
	return c.Status(fiber.StatusCreated).JSON(role)
}

// GetGroup handles the HTTP request to retrieve a group by its code
func (rc *RoleController) GetRole(c *fiber.Ctx) error {
	code := c.Params("code")
	role, err := rc.service.GetRoleByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
	}
	return c.JSON(role)
}

// GetAllRoles handles the HTTP request to retrieve all groups
func (rc *RoleController) GetAllRoles(c *fiber.Ctx) error {
	roles, err := rc.service.GetAllRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve roles"})
	}
	return c.JSON(roles)
}

// GetRolesForGroups handles the HTTP request to fetch roles for specific groups
func (rc *RoleController) GetRolesForGroups(c *fiber.Ctx) error {
	// Retrieve the group codes from the query parameters
	groupCodes := c.Query("group_codes")
	if groupCodes == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group codes are required",
		})
	}

	// Split the group codes into a slice
	groupCodeList := strings.Split(groupCodes, ",")

	// Call the service to get roles for the groups
	roles, err := rc.service.GetRolesForGroups(groupCodeList)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve roles for the specified groups",
		})
	}

	for _, role := range roles {
		fmt.Println("Role Name:", role.Name)
		fmt.Println("Group Name:", role.Group.Name) // Accessing preloaded Group
	}

	// Return the retrieved roles as JSON
	return c.JSON(roles)
}

// ResourceController handles resource-related endpoints
type ResourceController struct {
	service repository.ResourceService
}

// NewResourceController creates a new instance of ResourceController
func NewResourceController(service repository.ResourceService) *ResourceController {
	return &ResourceController{service: service}
}

// CreateResource handles the HTTP request to create a new resource
func (rc *ResourceController) CreateResource(c *fiber.Ctx) error {
	var resource userRegistration.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := rc.service.CreateResource(&resource); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create resource"})
	}
	return c.Status(fiber.StatusCreated).JSON(resource)
}

// GetResource handles the HTTP request to get a resource by code
func (rc *ResourceController) GetResource(c *fiber.Ctx) error {
	code := c.Params("code")
	resource, err := rc.service.GetResourceByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Resource not found"})
	}
	return c.JSON(resource)
}

// GetAllResources handles the HTTP request to get all non-deleted resources
func (rc *ResourceController) GetAllResources(c *fiber.Ctx) error {
	resources, err := rc.service.GetAllResources()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve resources"})
	}
	return c.JSON(resources)
}

// ============================================================================

// PermissionController handles permission-related endpoints
type PermissionController struct {
	service repository.PermissionService
}

// NewPermissionController creates a new instance of PermissionController
func NewPermissionController(service repository.PermissionService) *PermissionController {
	return &PermissionController{service: service}
}

// CreatePermission handles the HTTP request to create a new role-resource permission
func (pc *PermissionController) CreatePermission(c *fiber.Ctx) error {
	var permission userRegistration.RoleResourcePermission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	createdPermission, err := pc.service.CreatePermission(&permission)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create permission"})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPermission)
}

// CreateWildCardPermission handles the HTTP request to create a new wildcard permission
func (pc *PermissionController) CreateWildCardPermission(c *fiber.Ctx) error {
	var permission userRegistration.RoleWildCardPermission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	createdPermission, err := pc.service.CreateWildCardPermission(&permission)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create wildcard permission"})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPermission)
}

// GetGrantedPermissions handles the HTTP request to retrieve granted permissions
func (pc *PermissionController) GetGrantedPermissions(c *fiber.Ctx) error {
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")

	// Validate input parameters
	if roleCode == "" || resourceCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both role_code and resource_code are required",
		})
	}

	// Get granted permissions using the service
	permissions, err := pc.service.GetGrantedPermissions(roleCode, resourceCode)
	log.Printf("Permissions: %+v\n", permissions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve permissions",
		})
	}

	// Return the permissions as JSON
	return c.JSON(permissions)
}

func (pc *PermissionController) CheckPermissions(c *fiber.Ctx) error {
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")

	// Parse the requested permissions from the body
	requestedPerms := userRegistration.Permissions{}
	if err := c.BodyParser(&requestedPerms); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse requested permissions",
		})
	}

	// Get the granted permissions
	permissions, err := pc.service.CheckPermissions(roleCode, resourceCode, requestedPerms)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot check permissions",
		})
	}

	// Create a map to track permission checks
	permissionsMap := map[string]bool{
		"Read":    requestedPerms.Allow.R && !permissions.Allow.R || requestedPerms.Deny.R && permissions.Allow.R,
		"Write":   requestedPerms.Allow.W && !permissions.Allow.W || requestedPerms.Deny.W && permissions.Allow.W,
		"Execute": requestedPerms.Allow.X && !permissions.Allow.X || requestedPerms.Deny.X && permissions.Allow.X,
		"Delete":  requestedPerms.Allow.D && !permissions.Allow.D || requestedPerms.Deny.D && permissions.Allow.D,
	}

	// Check if any permission request is insufficient
	for perm, insufficient := range permissionsMap {
		if insufficient {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("%s permission denied", perm),
			})
		}
	}

	// If all checks pass, return OK
	return c.SendStatus(fiber.StatusOK)
}

// GetExplicitPermissions retrieves explicitly defined permissions for a role on a resource
func (pc *PermissionController) GetExplicitPermissions(c *fiber.Ctx) error {
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")

	permissions, err := pc.service.GetExplicitPermissions(roleCode, resourceCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve explicit permissions"})
	}

	return c.JSON(permissions)
}

// GetWildCardPermissions retrieves wildcard permissions for a role based on a resource pattern
func (pc *PermissionController) GetWildCardPermissions(c *fiber.Ctx) error {
	roleCode := c.Query("role_code")
	resourceCode := c.Query("resource_code")

	permissions, err := pc.service.GetWildCardPermissions(roleCode, resourceCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve wildcard permissions"})
	}

	return c.JSON(permissions)
}

// ResourceExplicitPermissionsExists checks if there are explicit (non-wildcard) permissions on the resource
func (pc *PermissionController) ResourceExplicitPermissionsExists(c *fiber.Ctx) error {
	resourceCode := c.Query("resource_code")

	exists, err := pc.service.ResourceExplicitPermissionsExists(resourceCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking explicit permissions"})
	}

	return c.JSON(fiber.Map{"exists": exists})
}

// GroupsWithRoleExists checks if there are groups assigned a specific role
func (pc *PermissionController) GroupsWithRoleExists(c *fiber.Ctx) error {
	groupCode := c.Query("group_code")

	exists, err := pc.service.GroupsWithRoleExists(groupCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking if groups with role exist"})
	}

	return c.JSON(fiber.Map{"exists": exists})
}

// GetPermissions retrieves permissions granted to roles on a specific resource
func (pc *PermissionController) GetPermissions(c *fiber.Ctx) error {
	// Get roleCodes and resourceCode from query parameters
	roleCodesStr := c.Query("role_codes") // Comma-separated role codes
	resourceCode := c.Query("resource_code")

	// Convert the comma-separated string into a slice of strings
	roleCodes := strings.Split(roleCodesStr, ",")

	// Call the service method to retrieve permissions
	permissions, err := pc.service.GetPermissions(roleCodes, resourceCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving permissions"})
	}

	// Return the retrieved permissions as a JSON response
	return c.JSON(permissions)
}
