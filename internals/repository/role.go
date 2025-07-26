package repository

import (
	"car-bond/internals/models/userRegistration"
	"errors"
	"fmt"

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

// =============================

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

// ===========================

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
