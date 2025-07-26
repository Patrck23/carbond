package repository

import (
	"car-bond/internals/models/userRegistration"
	"car-bond/internals/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *userRegistration.User) error
	GetPaginatedUsers(c *fiber.Ctx) (*utils.Pagination, []userRegistration.User, error)
	GetUserByID(id string) (*userRegistration.User, error)
	UpdateUser(user *userRegistration.User) error
	DeleteUserByID(id string) error
	GetPaginatedUsersByCompanyId(c *fiber.Ctx, companyId uint) (*utils.Pagination, []userRegistration.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *userRegistration.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) GetPaginatedUsers(c *fiber.Ctx) (*utils.Pagination, []userRegistration.User, error) {
	pagination, users, err := utils.Paginate(c, r.db.Preload("Company"), userRegistration.User{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, users, nil
}

func (r *UserRepositoryImpl) GetUserByID(id string) (*userRegistration.User, error) {
	var user userRegistration.User
	if err := r.db.Preload("Company").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil // Return a pointer
}

func (r *UserRepositoryImpl) UpdateUser(user *userRegistration.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) DeleteUserByID(id string) error {
	if err := r.db.Delete(&userRegistration.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// GetPaginatedUsersByCompanyId retrieves paginated users by company ID
func (r *UserRepositoryImpl) GetPaginatedUsersByCompanyId(c *fiber.Ctx, companyId uint) (*utils.Pagination, []userRegistration.User, error) {
	// Build the query with Preload
	query := r.db.Preload("Company").Where("company_id = ?", companyId)

	// Paginate results
	pagination, users, err := utils.Paginate(c, query, userRegistration.User{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to paginate users: %w", err)
	}

	return &pagination, users, nil
}
