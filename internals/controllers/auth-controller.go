package controllers

import (
	"errors"
	"net/mail"
	"time"

	"car-bond/internals/config"
	"car-bond/internals/database"
	"car-bond/internals/models/userRegistration"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*userRegistration.User, error) {
	db := database.DB.Db
	var user userRegistration.User
	if err := db.Where(&userRegistration.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*userRegistration.User, error) {
	db := database.DB.Db
	var user userRegistration.User
	if err := db.Where(&userRegistration.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	// Struct for user data
	type UserData struct {
		ID       uint     `json:"id"`
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Roles    []string `json:"roles"` // Added roles to the response
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
		user, err = getUserByEmail(identity)
	} else {
		user, err = getUserByUsername(identity)
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
	db := database.DB.Db
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

	// Retrieve session from context
	session := c.Locals("session").(*session.Session)

	// Store userId, username and roles in the session
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
	claims["roles"] = roleCodes // Include roles in the token
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
		Roles:    roleCodes,
	}

	// Return the token and user data
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"token":   t,
		"user":    userData,
	})
}
