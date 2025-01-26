package controllers

import (
	"errors"
	"net/mail"
	"time"

	"car-bond/internals/config"
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/models/userRegistration"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func getUserByEmail(e string, db *gorm.DB) (*userRegistration.User, error) {
	var user userRegistration.User
	if err := db.Where(&userRegistration.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string, db *gorm.DB) (*userRegistration.User, error) {
	var user userRegistration.User
	if err := db.Where(&userRegistration.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Utility function to validate email format
func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// CheckPasswordHash compares the password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx, db *gorm.DB) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Group    string `json:"group"`
		Location string `json:"location"`
	}

	// Parse the request body
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid login request format",
			"data":    err.Error(),
		})
	}

	identity := input.Identity
	password := input.Password
	var user *userRegistration.User
	var err error

	// Fetch user by email or username
	if isEmail(identity) {
		user, err = getUserByEmail(identity, db)
	} else {
		user, err = getUserByUsername(identity, db)
	}

	// Handle errors during user retrieval
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error retrieving user",
			"data":    err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid identity or password",
		})
	}

	// Validate password
	if !CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid identity or password",
		})
	}

	// Fetch roles based on the user's GroupID
	var roles []userRegistration.Role
	if err := db.Where("group_id = ?", user.GroupID).Find(&roles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error retrieving roles",
			"data":    err.Error(),
		})
	}

	// Extract role codes from the roles
	roleCodes := []string{}
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// Ensure that user.GroupID is valid
	if user.GroupID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User does not belong to a valid group",
			"data":    nil,
		})
	}

	// Fetch the group based on the user's GroupID
	var group userRegistration.Group
	if err := db.Where("id = ?", user.GroupID).First(&group).Error; err != nil {
		// Check if the error is due to a record not being found
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Group not found",
				"data":    nil,
			})
		}

		// Handle other errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error retrieving group",
			"data":    err.Error(),
		})
	}

	group_code := group.Code

	// company_id

	// Ensure that user.GroupID is valid
	if user.GroupID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User does not belong to a valid group",
			"data":    nil,
		})
	}

	// Fetch the group based on the user's GroupID
	var company companyRegistration.Company
	if err := db.Where("id = ?", user.CompanyID).First(&company).Error; err != nil {
		// Check if the error is due to a record not being found
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Group not found",
				"data":    nil,
			})
		}

		// Handle other errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error retrieving group",
			"data":    err.Error(),
		})
	}

	var companyLocation companyRegistration.CompanyLocation
	if err := db.Where("company_id = ?", company.ID).First(&companyLocation).Error; err != nil {
		// Check if the error is due to a record not being found
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Location not found",
				"data":    nil,
			})
		}

		// Handle other errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error retrieving location",
			"data":    err.Error(),
		})
	}

	country := companyLocation.Country

	// Retrieve session from context
	session := c.Locals("session").(*session.Session)

	// Store userId, username, and roles in the session
	session.Set("username", user.Username)
	session.Set("userId", user.ID)
	session.Set("roles", roleCodes)

	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save session",
		})
	}

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["roles"] = roleCodes
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	secretKey := config.Config("SECRET")
	if secretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Server configuration error",
		})
	}

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate token",
			"data":    err.Error(),
		})
	}

	// Populate UserData
	userData := UserData{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Group:    group_code,
		Location: country,
	}

	// Return the token and user data
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"token":   t,
		"user":    userData,
	})
}
