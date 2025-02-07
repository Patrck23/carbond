package userRegistration

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

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

type Permissions struct {
	Allow RWXD `json:"allow"` // Store Allow as JSON
	Deny  RWXD `json:"deny"`  // Store Deny as JSON
}

// Ensure `Permissions` implements `Valuer` and `Scanner`
func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Permissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert database value to byte slice")
	}
	return json.Unmarshal(bytes, p)
}

// RoleResourcePermission struct
type RoleResourcePermission struct {
	gorm.Model
	RoleCode     string      `gorm:"index"`      // Foreign key to the role
	ResourceCode string      `gorm:"index"`      // Foreign key to the resource
	Permissions  Permissions `gorm:"type:jsonb"` // Use jsonb for PostgreSQL, json for MySQL
	CreatedBy    string      `gorm:"size:100" json:"created_by"`
	UpdatedBy    string      `gorm:"size:100" json:"updated_by"`
}

// RoleWildCardPermission permissions for resource/role
type RoleWildCardPermission struct {
	gorm.Model
	RoleCode     string      `gorm:"size:50;index"` // Restrict RoleCode to 50 characters
	ResourceCode string      `gorm:"size:100"`      // Restrict ResourceCode to 100 characters    // ResourceCode allows define resource mask using "*" ("resource.*")
	Permissions  Permissions `gorm:"type:json"`
	CreatedBy    string      `gorm:"size:100" json:"created_by"`
	UpdatedBy    string      `gorm:"size:100" json:"updated_by"`
}
