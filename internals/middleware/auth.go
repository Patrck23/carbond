package middleware

import (
	"car-bond/internals/config"
	"car-bond/internals/models/userRegistration"
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config("SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

type DatabaseService struct {
	db *gorm.DB
}

func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{db: db}
}

// CheckPermissions checks the permissions for a given role and resource.
func (s *DatabaseService) CheckPermissions(roleCodes []string, resourceCode string) (userRegistration.Permissions, error) {
	var aggregatedPermissions userRegistration.Permissions

	// Iterate over all role codes to accumulate permissions
	for _, roleCode := range roleCodes {
		var roleResourcePermission userRegistration.RoleResourcePermission

		// Query to get the permissions as a JSON field (assuming permissions is a JSONB field)
		if err := s.db.Raw(
			`SELECT permissions
            FROM role_resource_permissions
            WHERE role_code = ? AND resource_code = ?`, roleCode, resourceCode).Scan(&roleResourcePermission).Error; err != nil {
			return aggregatedPermissions, err
		}

		// Print the raw permissions for debugging before aggregation
		fmt.Printf("Permissions for role %s: %+v\n", roleCode, roleResourcePermission.Permissions)

		// Aggregate the permissions across all roles (using OR logic to combine)
		aggregatedPermissions.Allow.R = aggregatedPermissions.Allow.R || roleResourcePermission.Permissions.Allow.R
		aggregatedPermissions.Allow.W = aggregatedPermissions.Allow.W || roleResourcePermission.Permissions.Allow.W
		aggregatedPermissions.Allow.X = aggregatedPermissions.Allow.X || roleResourcePermission.Permissions.Allow.X
		aggregatedPermissions.Allow.D = aggregatedPermissions.Allow.D || roleResourcePermission.Permissions.Allow.D

		// Note: You can also aggregate Deny permissions if needed
		aggregatedPermissions.Deny.R = aggregatedPermissions.Deny.R || roleResourcePermission.Permissions.Deny.R
		aggregatedPermissions.Deny.W = aggregatedPermissions.Deny.W || roleResourcePermission.Permissions.Deny.W
		aggregatedPermissions.Deny.X = aggregatedPermissions.Deny.X || roleResourcePermission.Permissions.Deny.X
		aggregatedPermissions.Deny.D = aggregatedPermissions.Deny.D || roleResourcePermission.Permissions.Deny.D
	}

	// Print the final aggregated permissions
	fmt.Printf("Aggregated Permissions: %+v\n", aggregatedPermissions)

	// Return the aggregated permissions
	return aggregatedPermissions, nil
}

// PermissionMiddleware is a Fiber middleware that checks user permissions.
func PermissionMiddleware(service *DatabaseService, resourceCode string, requestedPerms []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles := getRolesFromRequest(c)

		// Check if requested permissions are empty
		if len(requestedPerms) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No permissions requested",
			})
		}

		// Get the granted permissions
		permissions, err := service.CheckPermissions(roles, resourceCode)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot check permissions",
			})
		}

		// Create a map to track permission checks
		permissionsMap := map[string]bool{
			"r": permissions.Allow.R,
			"w": permissions.Allow.W,
			"x": permissions.Allow.X,
			"d": permissions.Allow.D,
		}

		// Check if any requested permission is insufficient
		for _, perm := range requestedPerms {
			if allowed, exists := permissionsMap[perm]; !exists || !allowed {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": fmt.Sprintf("%s permission denied", perm),
				})
			}
		}

		// If all checks pass, proceed to the next handler
		return c.Next()
	}
}

func getRolesFromRequest(c *fiber.Ctx) []string {
	// Get the session object from Locals (assuming it's stored there)
	session, ok := c.Locals("session").(*session.Session)
	if !ok || session == nil {
		return nil // No session available
	}

	// Get the roles from the session, assuming roles are stored as a slice in the session
	rolesInterface := session.Get("roles")
	if rolesInterface == nil {
		return nil // No roles available
	}

	// If roles are stored as a slice of strings
	roles, ok := rolesInterface.([]string)
	if !ok {
		return nil // Roles aren't in the expected format
	}

	return roles
}
