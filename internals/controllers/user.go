package controllers

import (
	"car-bond/internals/models/userRegistration"
	"car-bond/internals/repository"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	repo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{repo: repo}
}

// ============================================

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ======================

func (h *UserController) CreateUser(c *fiber.Ctx) error {
	// Struct to return minimal user details
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	// Parse request body into user object
	user := new(userRegistration.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"data":    err.Error(),
		})
	}

	// Hash the user's password
	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't hash password",
			"data":    err.Error(),
		})
	}
	user.Password = hash

	// Create the user in the database
	if err := h.repo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create user",
			"data":    err.Error(),
		})
	}

	// Prepare response with minimal user details
	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Created user",
		"data":    newUser,
	})
}

// =====================

func (h *UserController) GetAllUsers(c *fiber.Ctx) error {
	pagination, users, err := h.repo.GetPaginatedUsers(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve users",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Users retrieved successfully",
		"data":    users,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ======================

// GetCarSale fetches a sale with its associated contacts and addresses from the database
func (h *UserController) GetUserByID(c *fiber.Ctx) error {
	// Get the user ID from the route parameters
	id := c.Params("id")

	// Fetch the user by ID
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve User",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User and associated data retrieved successfully",
		"data":    user,
	})
}

// ======================

func (h *UserController) UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Surname   string `json:"surname"`
		Firstname string `json:"firstname"`
		Othername string `json:"othername"`
		Gender    string `json:"gender"`
		Title     string `json:"title"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		CompanyID uint   `json:"company_id"`
		GroupID   uint   `json:"group_id"`
		UpdatedBy string `json:"updated_by"`
	}

	// Parse and validate request body
	var input UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Retrieve user ID from request parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User ID is required",
		})
	}

	// Fetch the user from the repository
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch user",
			"error":   err.Error(),
		})
	}

	// Update user fields
	user.Surname = input.Surname
	user.Firstname = input.Firstname
	user.Othername = input.Othername
	user.Gender = input.Gender
	user.Title = input.Title
	user.Email = input.Email
	user.CompanyID = input.CompanyID
	user.GroupID = input.GroupID
	user.UpdatedBy = input.UpdatedBy

	// Save the updated user
	if err := h.repo.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully updated",
		"data":    user,
	})
}

// ====================

// DeleteUser delete user
func (h *UserController) DeleteUserByID(c *fiber.Ctx) error {

	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")

	// Find the user in the database
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "user not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find user",
			"data":    err.Error(),
		})
	}

	// Delete the user
	if err := h.repo.DeleteUserByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete user",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
		"data":    user,
	})
}

// =========================

// GetUsersByCompany retrieves paginated users by company ID
func (h *UserController) GetUsersByCompany(c *fiber.Ctx) error {
	// Get companyId from URL params and convert to uint
	companyIdStr := c.Params("companyId")
	companyId, err := strconv.ParseUint(companyIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid company ID",
			"data":    err.Error(),
		})
	}

	// Fetch paginated users by company ID
	pagination, users, err := h.repo.GetPaginatedUsersByCompanyId(c, uint(companyId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve users",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Users retrieved successfully",
		"data":    users,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}
