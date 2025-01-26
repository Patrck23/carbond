package middleware

import (
	"car-bond/internals/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/session"
	// "car-bond/internals/service"
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

// var SecurityService *service.SecurityServiceWrapper

// // CheckPermissionsMiddleware is a middleware function that checks permissions for the requested route
// func CheckPermissionsMiddleware(resource string, requestedPermissions []string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		roles := getRolesFromRequest(c)

// 		if SecurityService == nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "Security service is not initialized",
// 			})
// 		}

// 		// Check if requested permissions are allowed
// 		allowed, err := SecurityService.CheckPermissions(c, resource, roles, requestedPermissions)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "Failed to check permissions",
// 			})
// 		}

// 		// If not allowed, return forbidden error
// 		if !allowed {
// 			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 				"error": "You do not have permission to access this resource",
// 			})
// 		}

// 		// If allowed, continue to the next handler
// 		return c.Next()
// 	}
// }

// func getRolesFromRequest(c *fiber.Ctx) []string {
// 	// Get the session object from Locals (assuming it's stored there)
// 	session, ok := c.Locals("session").(*session.Session)
// 	if !ok || session == nil {
// 		return nil // No session available
// 	}

// 	// Get the roles from the session, assuming roles are stored as a slice in the session
// 	rolesInterface := session.Get("roles")
// 	if rolesInterface == nil {
// 		return nil // No roles available
// 	}

// 	// If roles are stored as a slice of strings
// 	roles, ok := rolesInterface.([]string)
// 	if !ok {
// 		return nil // Roles aren't in the expected format
// 	}

// 	return roles
// }
