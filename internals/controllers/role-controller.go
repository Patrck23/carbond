package controllers

import (
	"car-bond/internals/models/userRegistration"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GroupService defines the interface for group operations
type GroupService interface {
	CreateGroup(group *userRegistration.Group) error
	GetGroupByCode(code string) (*userRegistration.Group, error)
	GetAllGroups() ([]userRegistration.Group, error)
	UpdateGroup(code string, group *userRegistration.Group) error
	DeleteGroup(code string) error
}

// DatabaseService implements GroupService using a database
type DatabaseService struct {
	db *gorm.DB
}

// NewDatabaseService creates a new instance of DatabaseService
func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{db: db}
}

// CreateGroup creates a new group in the database
func (s *DatabaseService) CreateGroup(group *userRegistration.Group) error {
	return s.db.Create(group).Error
}

// GetGroupByCode retrieves a group by its code from the database
func (s *DatabaseService) GetGroupByCode(code string) (*userRegistration.Group, error) {
	var group userRegistration.Group
	if err := s.db.Where("code = ?", code).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// GetAllGroups retrieves all groups from the database
func (s *DatabaseService) GetAllGroups() ([]userRegistration.Group, error) {
	var groups []userRegistration.Group
	if err := s.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// UpdateGroup updates only the specified fields of an existing group in the database
func (s *DatabaseService) UpdateGroup(code string, group *userRegistration.Group) error {
	// Find the group by its code
	var existingGroup userRegistration.Group
	if err := s.db.First(&existingGroup, "code = ?", code).Error; err != nil {
		return err
	}

	// Only update the fields: name, description, internal, and updated_by
	existingGroup.Name = group.Name
	existingGroup.Description = group.Description
	existingGroup.Internal = group.Internal
	existingGroup.UpdatedBy = group.UpdatedBy

	// Save the updated group
	return s.db.Save(&existingGroup).Error
}

// DeleteGroup deletes a group by its code from the database
func (s *DatabaseService) DeleteGroup(code string) error {
	// Find the group by its code and delete it
	if err := s.db.Delete(&userRegistration.Group{}, "code = ?", code).Error; err != nil {
		return err
	}
	return nil
}

// GroupController handles group-related endpoints
type GroupController struct {
	service GroupService
}

// NewGroupController creates a new instance of GroupController
func NewGroupController(service GroupService) *GroupController {
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

// RoleService defines the interface for role operations
type RoleService interface {
	CreateRole(role *userRegistration.Role) error
	GetRoleByCode(code string) (*userRegistration.Role, error)
	GetAllRoles() ([]userRegistration.Role, error)
	GetRolesForGroups(groupCodes []string) ([]userRegistration.Role, error)
}

// CreateRole creates a new role in the database
func (s *DatabaseService) CreateRole(role *userRegistration.Role) error {
	return s.db.Create(role).Error
}

// GetRoleByCode retrieves a role by its code from the database
func (s *DatabaseService) GetRoleByCode(code string) (*userRegistration.Role, error) {
	var role userRegistration.Role
	if err := s.db.Preload("Group").Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAllRoles retrieves all roles from the database
func (s *DatabaseService) GetAllRoles() ([]userRegistration.Role, error) {
	var roles []userRegistration.Role
	if err := s.db.Preload("Group").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRolesForGroups fetches roles associated with the provided group codes
func (s *DatabaseService) GetRolesForGroups(groupCodes []string) ([]userRegistration.Role, error) {
	var roles []userRegistration.Role
	query := `
    SELECT DISTINCT roles.*
    FROM roles
    JOIN groups ON groups.id = roles.group_id
    WHERE roles.group_id IN (
        SELECT id
        FROM groups
        WHERE code IN ?
    )
	`
	if err := s.db.Preload("Group").Raw(query, groupCodes).Scan(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// RoleController handles role-related endpoints
type RoleController struct {
	service RoleService
}

// NewRoleController creates a new instance of RoleController
func NewRoleController(service RoleService) *RoleController {
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

// ==========================================================================

// ResourceService defines the interface for resource operations
type ResourceService interface {
	CreateResource(resource *userRegistration.Resource) error
	GetResourceByCode(code string) (*userRegistration.Resource, error)
	GetAllResources() ([]userRegistration.Resource, error)
}

// CreateResource creates a new resource in the database
func (s *DatabaseService) CreateResource(resource *userRegistration.Resource) error {
	return s.db.Create(resource).Error
}

// GetResourceByCode retrieves a resource by its code
func (s *DatabaseService) GetResourceByCode(code string) (*userRegistration.Resource, error) {
	var resource userRegistration.Resource
	if err := s.db.Where("code = ?", code).First(&resource).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetAllResources retrieves all non-deleted resources
func (s *DatabaseService) GetAllResources() ([]userRegistration.Resource, error) {
	var resources []userRegistration.Resource
	if err := s.db.Where("internal = ?", false).Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

// ResourceController handles resource-related endpoints
type ResourceController struct {
	service ResourceService
}

// NewResourceController creates a new instance of ResourceController
func NewResourceController(service ResourceService) *ResourceController {
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

// PermissionService defines the interface for permission operations
type PermissionService interface {
	CreatePermission(permission *userRegistration.RoleResourcePermission) (*userRegistration.RoleResourcePermission, error)
	CreateWildCardPermission(permission *userRegistration.RoleWildCardPermission) (*userRegistration.RoleWildCardPermission, error)
	GetGrantedPermissions(roleCode, resourceCode string) (*userRegistration.Permissions, error)
	CheckPermissions(roleCode, resourceCode string, requestedPerms userRegistration.Permissions) (userRegistration.Permissions, error)
	GetExplicitPermissions(roleCode, resourceCode string) (userRegistration.RoleResourcePermission, error)
	GetWildCardPermissions(roleCode, resourceCode string) ([]userRegistration.RoleWildCardPermission, error)
	ResourceExplicitPermissionsExists(resourceCode string) (bool, error)
	GroupsWithRoleExists(roleCode string) (bool, error)
	GetPermissions(roleCodes []string, resourceCode string) ([]userRegistration.RoleResourcePermission, error)
}

// CreatePermission creates a new role-resource permission
func (s *DatabaseService) CreatePermission(permission *userRegistration.RoleResourcePermission) (*userRegistration.RoleResourcePermission, error) {
	if err := s.db.Create(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

// CreateWildCardPermission creates a new wildcard permission
func (s *DatabaseService) CreateWildCardPermission(permission *userRegistration.RoleWildCardPermission) (*userRegistration.RoleWildCardPermission, error) {
	if err := s.db.Create(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

// GetGrantedPermissions calculates permissions for roles on a resource
func (s *DatabaseService) GetGrantedPermissions(roleCode, resourceCode string) (*userRegistration.Permissions, error) {
	var roleResourcePermission userRegistration.RoleResourcePermission

	// Query to retrieve the entire row
	if err := s.db.Raw(`
		SELECT permissions
		FROM role_resource_permissions
		WHERE role_code = ? AND resource_code = ?
		LIMIT 1
	`, roleCode, resourceCode).Scan(&roleResourcePermission).Error; err != nil {
		return nil, err
	}

	// Return only the permissions field
	return &roleResourcePermission.Permissions, nil
}

// Check checks the permissions for a given role and resource
func (s *DatabaseService) CheckPermissions(roleCode, resourceCode string, requestedPerms userRegistration.Permissions) (userRegistration.Permissions, error) {
	var permissions userRegistration.Permissions
	if err := s.db.Raw(`
        SELECT permissions
        FROM role_resource_permissions
        WHERE role_code = ? AND resource_code = ?
    `, roleCode, resourceCode).Scan(&permissions.Allow).Error; err != nil {
		return permissions, err
	}
	return permissions, nil
}

// Get retrieves explicitly defined permissions for a role on a resource
func (s *DatabaseService) GetExplicitPermissions(roleCode, resourceCode string) (userRegistration.RoleResourcePermission, error) {
	var permissions userRegistration.RoleResourcePermission

	// Use the `Select` method to retrieve only the permissions column.
	if err := s.db.Where("role_code = ? AND resource_code = ?", roleCode, resourceCode).First(&permissions).Error; err != nil {
		fmt.Println("Error fetching permissions:", err)
		return permissions, err // Return the empty struct and the error
	}

	return permissions, nil // Return the permissions and no error
}

// Get retrieves wildcard permissions for a role based on a resource pattern
func (s *DatabaseService) GetWildCardPermissions(roleCode, resourceCode string) ([]userRegistration.RoleWildCardPermission, error) {
	var permissions []userRegistration.RoleWildCardPermission
	if err := s.db.Where("role_code = ? AND resource_code = ?", roleCode, resourceCode).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// Exists checks if there are explicit (non-wildcard) permissions on the resource
func (s *DatabaseService) ResourceExplicitPermissionsExists(resourceCode string) (bool, error) {
	var count int64
	err := s.db.Model(&userRegistration.RoleResourcePermission{}).
		Where("resource_code = ?", resourceCode).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Exists checks if there are groups assigned a specific role
func (s *DatabaseService) GroupsWithRoleExists(groupCode string) (bool, error) {
	var group userRegistration.Group

	// First, find the group by its Code
	if err := s.db.Where("code = ?", groupCode).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Group does not exist, so return false
		}
		return false, err // Return error if something goes wrong
	}

	// Now count roles assigned to the found group
	var count int64
	err := s.db.Model(&userRegistration.Role{}).
		Where("group_id = ?", group.ID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Get retrieves permissions granted to roles on a specific resource
func (s *DatabaseService) GetPermissions(roleCodes []string, resourceCode string) ([]userRegistration.RoleResourcePermission, error) {
	var permissions []userRegistration.RoleResourcePermission
	err := s.db.Where("role_code IN (?) AND resource_code = ?", roleCodes, resourceCode).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// PermissionController handles permission-related endpoints
type PermissionController struct {
	service PermissionService
}

// NewPermissionController creates a new instance of PermissionController
func NewPermissionController(service PermissionService) *PermissionController {
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
