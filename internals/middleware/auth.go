package middleware

import (
	"car-bond/internals/config"
	"car-bond/internals/models/userRegistration"
	"fmt"
	"strings"

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
	var allPermissions []userRegistration.Permissions // Slice to hold permissions from all roles
	var aggregatedPermissions userRegistration.Permissions

	// Iterate over all role codes to accumulate permissions
	for _, roleCode := range roleCodes {
		var roleResourcePermission userRegistration.RoleResourcePermission

		// Query to get the explicit permissions
		if err := s.db.Raw(
			`SELECT permissions
            FROM role_resource_permissions
            WHERE role_code = ? AND resource_code = ?`, roleCode, resourceCode).Scan(&roleResourcePermission).Error; err != nil {
			return aggregatedPermissions, err
		}

		// Append explicit permissions to the slice
		allPermissions = append(allPermissions, roleResourcePermission.Permissions)

		// Query to get the wildcard permissions
		var wildCardResourcePermissions userRegistration.RoleWildCardPermission
		if err := s.db.Raw(
			`SELECT permissions
            FROM role_wild_card_permissions
            WHERE role_code = ? AND resource_code = ?`, roleCode, resourceCode).Scan(&wildCardResourcePermissions).Error; err != nil {
			return aggregatedPermissions, err
		}

		// Append wildcard permissions to the slice
		allPermissions = append(allPermissions, wildCardResourcePermissions.Permissions)
	}

	// Aggregate the permissions across all roles (using OR logic to combine)
	for _, permissions := range allPermissions {
		aggregatedPermissions.Allow.R = aggregatedPermissions.Allow.R || permissions.Allow.R
		aggregatedPermissions.Allow.W = aggregatedPermissions.Allow.W || permissions.Allow.W
		aggregatedPermissions.Allow.X = aggregatedPermissions.Allow.X || permissions.Allow.X
		aggregatedPermissions.Allow.D = aggregatedPermissions.Allow.D || permissions.Allow.D

		// Aggregate Deny permissions if needed
		aggregatedPermissions.Deny.R = aggregatedPermissions.Deny.R || permissions.Deny.R
		aggregatedPermissions.Deny.W = aggregatedPermissions.Deny.W || permissions.Deny.W
		aggregatedPermissions.Deny.X = aggregatedPermissions.Deny.X || permissions.Deny.X
		aggregatedPermissions.Deny.D = aggregatedPermissions.Deny.D || permissions.Deny.D
	}

	// Print the final aggregated permissions
	fmt.Printf("Aggregated Permissions: %+v\n", aggregatedPermissions)

	// Return the aggregated permissions
	return aggregatedPermissions, nil
}

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

		// Check if the requested permissions are allowed
		for _, perm := range requestedPerms {
			switch perm {
			case "R":
				if !permissions.Allow.R {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Read permission denied",
					})
				}
			case "W":
				if !permissions.Allow.W {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Write permission denied",
					})
				}
			case "X":
				if !permissions.Allow.X {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Execute permission denied",
					})
				}
			case "D":
				if !permissions.Allow.D {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Delete permission denied",
					})
				}
			}
		}

		// If all checks pass, proceed to the next handler
		return c.Next()
	}
}

func RequireGroupMembership(allowedGroups ...string) fiber.Handler {
	allowedSet := make(map[string]struct{}, len(allowedGroups))
	for _, group := range allowedGroups {
		allowedSet[group] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		userRoles := getRolesFromRequest(c)
		if len(userRoles) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No roles found in session",
			})
		}

		for _, role := range userRoles {
			parts := strings.Split(role, ".")
			if len(parts) > 1 {
				suffix := parts[len(parts)-1]
				if _, ok := allowedSet[suffix]; ok {
					return c.Next() // Group match found
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied: user does not belong to required group",
		})
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
