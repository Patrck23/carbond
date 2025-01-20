package userRegistration

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Code        string `gorm:"unique;not null"` // Unique and not null
	Name        string `gorm:"not null"`        // Not null
	Description string `gorm:"size:255"`        // Limits string size
	Internal    bool   `gorm:"default:false"`   // Defaults to false
	CreatedBy   string `gorm:"size:100" json:"created_by"`
	UpdatedBy   string `gorm:"size:100" json:"updated_by"`
	Roles       []Role `gorm:"foreignKey:GroupID" json:"roles"`
}

type Role struct {
	gorm.Model
	Code        string `gorm:"not null"`      // Unique and not null
	Name        string `gorm:"not null"`      // Not null
	Description string `gorm:"size:255"`      // Limits string size
	Internal    bool   `gorm:"default:false"` // Defaults to false
	CreatedBy   string `gorm:"size:100" json:"created_by"`
	UpdatedBy   string `gorm:"size:100" json:"updated_by"`
	GroupID     uint   `json:"group_id"`                         // Foreign key to Group
	Group       Group  `gorm:"foreignKey:GroupID;references:ID"` // Belongs to a Group
}

type Resource struct {
	gorm.Model
	Code        string `gorm:"unique;not null"` // Unique and not null
	Name        string `gorm:"not null"`        // Not null
	Description string `gorm:"size:255"`        // Limits string size
	Internal    bool   `gorm:"default:false"`   // Defaults to false
	CreatedBy   string `gorm:"size:100" json:"created_by"`
	UpdatedBy   string `gorm:"size:100" json:"updated_by"`
}

// RWXD represents a set of permissions (Read, Write, Execute, Delete)
type RWXD struct {
	R bool `json:"r"` // Read
	W bool `json:"w"` // Write
	X bool `json:"x"` // Execute
	D bool `json:"d"` // Delete
}

// Implement the Valuer interface for RWXD (convert to JSON for database storage)
func (r RWXD) Value() (driver.Value, error) {
	return json.Marshal(r) // Convert RWXD to JSON
}

// Implement the Scanner interface for RWXD (convert JSON from database to struct)
func (r *RWXD) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert database value to byte slice")
	}
	return json.Unmarshal(bytes, r) // Convert JSON to RWXD
}

// Permissions specifies allow/deny permissions on a resource
type Permissions struct {
	Allow RWXD `gorm:"type:json"` // Store Allow as JSON
	Deny  RWXD `gorm:"type:json"` // Store Deny as JSON
}

type RoleResourcePermission struct {
	gorm.Model
	RoleCode     string      `gorm:"index"` // Foreign key to the role
	ResourceCode string      `gorm:"index"` // Foreign key to the resource
	Permissions  Permissions `gorm:"type:json"`
	CreatedBy    string      `gorm:"size:100" json:"created_by"`
	UpdatedBy    string      `gorm:"size:100" json:"updated_by"`
}

// RoleWildCardPermission permissions for resource/role
type RoleWildCardPermission struct {
	gorm.Model
	RoleCode        string      `gorm:"size:50;index"` // Restrict RoleCode to 50 characters
	ResourcePattern string      `gorm:"size:100"`      // Restrict ResourcePattern to 100 characters    // ResourcePattern allows define resource mask using "*" ("resource.*")
	Permissions     Permissions `gorm:"type:json"`
	CreatedBy       string      `gorm:"size:100" json:"created_by"`
	UpdatedBy       string      `gorm:"size:100" json:"updated_by"`
}

type SecurityService interface {
	// GetGroup retrieves a group by code
	GetGroup(c *fiber.Ctx, code string) (*Group, error)
	// GetAllGroups retrieves all not deleted groups
	GetAllGroups(c *fiber.Ctx) ([]*Group, error)
	// GetRole retrieves a role by code
	GetRole(c *fiber.Ctx, code string) (*Role, error)
	// GetAllRoles retrieves all not deleted roles
	GetAllRoles(c *fiber.Ctx) ([]*Role, error)
	// GetRolesForGroups retrieves roles assigned on groups
	GetRolesForGroups(c *fiber.Ctx, groups []string) ([]string, error)
	// GetResource retrieves a resource by code
	GetResource(c *fiber.Ctx, code string) (*Resource, error)
	// GetAllResources retrieves all not deleted resources
	GetAllResources(c *fiber.Ctx) ([]*Resource, error)
	// GetGrantedPermissions calculates permissions on the resource for the roles and applies allow/deny logic
	GetGrantedPermissions(c *fiber.Ctx, resource string, roles []string) (*RWXD, error)
	// CheckPermissions checks if the roles have the requested perms on the given resource
	CheckPermissions(c *fiber.Ctx, resource string, roles []string, requestedPermissions []string) (bool, error)
	// GetExplicitPermissions returns permissions on resource / roles setup explicitly
	GetExplicitPermissions(c *fiber.Ctx, resources []string, roles []string) ([]*RoleResourcePermission, error)
	// GetWildCardPermissions returns wildcard permissions on roles
	GetWildCardPermissions(c *fiber.Ctx, roles []string) ([]*RoleWildCardPermission, error)
}

type SecurityStorage interface {
	// GetGroup retrieves a group by code
	GetGroup(c *fiber.Ctx, code string) (*Group, error)
	// GetGroups retrieves all not deleted groups
	GetGroups(c *fiber.Ctx) ([]*Group, error)
	// GetRole retrieves a role by code
	GetRole(c *fiber.Ctx, code string) (*Role, error)
	// GetAllRoles retrieves all not deleted roles
	GetAllRoles(c *fiber.Ctx) ([]*Role, error)
	// GetAllRoleCodes retrieves all role codes
	GetAllRoleCodes(c *fiber.Ctx) ([]string, error)
	// GetResource retrieves a resource by code
	GetResource(c *fiber.Ctx, code string) (*Resource, error)
	// GetAllResources retrieves all not deleted resources
	GetAllResources(c *fiber.Ctx) ([]*Resource, error)
	// ResourceExplicitPermissionsExists checks if there are explicit (no wildcard) permissions on the resource
	ResourceExplicitPermissionsExists(c *fiber.Ctx, code string) (bool, error)
	// GetRoleCodesForGroups retrieves role codes for groups
	GetRoleCodesForGroups(c *fiber.Ctx, groups []string) ([]string, error)
	// GroupsWithRoleExists checks if there are groups with assigned role
	GroupsWithRoleExists(c *fiber.Ctx, role string) (bool, error)
	// GetPermissions retrieves permissions granted to roles on resource
	GetPermissions(c *fiber.Ctx, resource string, roles []string) ([]*Permissions, error)
	// GetWildcardPermissions retrieves wildcard permissions granted to roles on resource
	GetWildcardPermissions(c *fiber.Ctx, resource string, roles []string) ([]*Permissions, error)
}
