package controllers

import (
	"strconv"

	"car-bond/internals/database"
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/models/userRegistration"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB.Db
	var user userRegistration.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// Get All users from db
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB.Db
	var users []userRegistration.User
	// find all users in the database
	db.Preload("Company").Find(&users)
	// If no customer found, return an error
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "users not found"})
	}
	// return users
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "users Found", "data": users})
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var user userRegistration.User
	db.Find(&user, id)
	if user.Username == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	db := database.DB.Db
	user := new(userRegistration.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	}

	user.Password = hash
	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Surname   string                      `gorm:"size:100;not null" json:"surname"`
		Firstname string                      `gorm:"size:100;not null" json:"firstname"`
		Othername string                      `gorm:"size:100" json:"othername"`
		Gender    string                      `gorm:"size:10;not null" json:"gender"`
		Title     string                      `gorm:"size:50" json:"title"`
		Email     string                      `gorm:"uniqueIndex;not null" json:"email"`
		Password  string                      `gorm:"size:255;not null" json:"password"`
		CompanyID uint                        `json:"company_id"`
		Company   companyRegistration.Company `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
		UpdatedBy string                      `gorm:"size:100" json:"updated_by"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	db := database.DB.Db
	var user userRegistration.User

	db.First(&user, id)
	user.Surname = uui.Surname
	user.Firstname = uui.Firstname
	user.Othername = uui.Othername
	user.Gender = uui.Gender
	user.Title = uui.Title
	user.Email = uui.Email
	user.CompanyID = uui.CompanyID
	user.UpdatedBy = uui.UpdatedBy
	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	if !validUser(id, pi.Password) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})
	}

	db := database.DB.Db
	var user userRegistration.User

	db.First(&user, id)

	db.Delete(&user)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
