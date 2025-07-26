package repository

import (
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	CreateCompany(company *companyRegistration.Company) error
	GetPaginatedCompanies(c *fiber.Ctx) (*utils.Pagination, []companyRegistration.Company, error)
	GetCompanyByID(id string) (companyRegistration.Company, error)
	GetCompanyLocations(companyID uint) ([]companyRegistration.CompanyLocation, error)
	UpdateCompany(company *companyRegistration.Company) error
	DeleteByID(id string) error

	// Expenses
	CreateCompanyExpense(expense *companyRegistration.CompanyExpense) error
	GetPaginatedExpenses(c *fiber.Ctx, companyId string) (*utils.Pagination, []companyRegistration.CompanyExpense, error)
	FindCompanyExpenseByIdAndCompanyId(id, companyId string) (*companyRegistration.CompanyExpense, error)
	FindCompanyExpenseById(id string) (*companyRegistration.CompanyExpense, error)
	UpdateCompanyExpense(expense *companyRegistration.CompanyExpense) error
	FindCompanyExpenseByCompanyAndId(companyId, expenseId string) (*companyRegistration.CompanyExpense, error)
	DeleteCompanyExpense(expense *companyRegistration.CompanyExpense) error
	FindCompanyExpensesByCompanyIdAndExpenseDate(companyId, expenseDate string) ([]companyRegistration.CompanyExpense, error)
	FindCompanyExpensesByCompanyIdAndExpenseDescription(companyId, expenseDescription string) ([]companyRegistration.CompanyExpense, error)
	FindCompanyExpensesByCompanyIdAndCurrency(companyId, currency string) ([]companyRegistration.CompanyExpense, error)
	FindCompanyExpensesByThree(companyId, expenseDate, currency string) ([]companyRegistration.CompanyExpense, error)
	GetCompanyExpensesByFour(companyId, expenseDate, expenseDescription, currency string) ([]companyRegistration.CompanyExpense, error)
	GetPaginatedAllExpenses(c *fiber.Ctx) (*utils.Pagination, []companyRegistration.CompanyExpense, error)

	// Locations
	CreateCompanyLocation(location *companyRegistration.CompanyLocation) error
	GetAllCompanyLocations(companyId string) ([]companyRegistration.CompanyLocation, error)
	GetLocationByCompanyId(id, companyId string) (*companyRegistration.CompanyLocation, error)
	FindLocationById(Id string) (*companyRegistration.CompanyLocation, error)
	UpdateCompanyLocation(location *companyRegistration.CompanyLocation) error
	DeleteLocationByID(id string) error
}

type CompanyRepositoryImpl struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &CompanyRepositoryImpl{db: db}
}

func (r *CompanyRepositoryImpl) CreateCompany(company *companyRegistration.Company) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepositoryImpl) GetPaginatedCompanies(c *fiber.Ctx) (*utils.Pagination, []companyRegistration.Company, error) {
	pagination, companies, err := utils.Paginate(c, r.db, companyRegistration.Company{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, companies, nil
}

func (r *CompanyRepositoryImpl) GetCompanyLocations(companyID uint) ([]companyRegistration.CompanyLocation, error) {
	var locations []companyRegistration.CompanyLocation
	err := r.db.Where("company_id = ?", companyID).Find(&locations).Error
	return locations, err
}

func (r *CompanyRepositoryImpl) GetCompanyByID(id string) (companyRegistration.Company, error) {
	var company companyRegistration.Company
	err := r.db.First(&company, "id = ?", id).Error
	return company, err
}

func (r *CompanyRepositoryImpl) UpdateCompany(company *companyRegistration.Company) error {
	return r.db.Save(company).Error
}

// DeleteByID deletes a company by ID
func (r *CompanyRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&companyRegistration.Company{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CreateCompanyExpense creates a new company expense in the database
func (r *CompanyRepositoryImpl) CreateCompanyExpense(expense *companyRegistration.CompanyExpense) error {
	return r.db.Create(expense).Error
}

func (r *CompanyRepositoryImpl) GetPaginatedExpenses(c *fiber.Ctx, companyId string) (*utils.Pagination, []companyRegistration.CompanyExpense, error) {
	pagination, expenses, err := utils.Paginate(c, r.db.Preload("Company").Where("company_id = ?", companyId), companyRegistration.CompanyExpense{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, expenses, nil
}

func (r *CompanyRepositoryImpl) FindCompanyExpenseByIdAndCompanyId(id, companyId string) (*companyRegistration.CompanyExpense, error) {
	var expense companyRegistration.CompanyExpense
	result := r.db.Preload("Company").Where("id = ? AND company_id = ?", id, companyId).First(&expense)
	if result.Error != nil {
		return nil, result.Error
	}
	return &expense, nil
}

func (r *CompanyRepositoryImpl) FindCompanyExpenseById(id string) (*companyRegistration.CompanyExpense, error) {
	var expense companyRegistration.CompanyExpense
	if err := r.db.First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *CompanyRepositoryImpl) UpdateCompanyExpense(expense *companyRegistration.CompanyExpense) error {
	return r.db.Save(expense).Error
}

func (r *CompanyRepositoryImpl) FindCompanyExpenseByCompanyAndId(companyId, expenseId string) (*companyRegistration.CompanyExpense, error) {
	var expense companyRegistration.CompanyExpense
	if err := r.db.Where("id = ? AND company_id = ?", expenseId, companyId).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *CompanyRepositoryImpl) DeleteCompanyExpense(expense *companyRegistration.CompanyExpense) error {
	return r.db.Delete(expense).Error
}

func (r *CompanyRepositoryImpl) FindCompanyExpensesByCompanyIdAndExpenseDate(companyId, expenseDate string) ([]companyRegistration.CompanyExpense, error) {
	var expenses []companyRegistration.CompanyExpense
	if err := r.db.Where("company_id = ? AND expense_date = ?", companyId, expenseDate).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CompanyRepositoryImpl) FindCompanyExpensesByCompanyIdAndExpenseDescription(companyId, expenseDescription string) ([]companyRegistration.CompanyExpense, error) {
	var expenses []companyRegistration.CompanyExpense
	if err := r.db.Where("company_id = ? AND description = ?", companyId, expenseDescription).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CompanyRepositoryImpl) FindCompanyExpensesByCompanyIdAndCurrency(companyId, currency string) ([]companyRegistration.CompanyExpense, error) {
	var expenses []companyRegistration.CompanyExpense
	if err := r.db.Where("company_id = ? AND currency = ?", companyId, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CompanyRepositoryImpl) FindCompanyExpensesByThree(companyId, expenseDate, currency string) ([]companyRegistration.CompanyExpense, error) {
	var expenses []companyRegistration.CompanyExpense
	if err := r.db.Where("company_id = ? AND expense_date = ? AND currency = ?", companyId, expenseDate, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CompanyRepositoryImpl) GetCompanyExpensesByFour(companyId, expenseDate, expenseDescription, currency string) ([]companyRegistration.CompanyExpense, error) {
	var expenses []companyRegistration.CompanyExpense
	if err := r.db.Where("company_id = ? AND expense_date = ? AND description = ? AND currency = ?", companyId, expenseDate, expenseDescription, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CompanyRepositoryImpl) GetPaginatedAllExpenses(c *fiber.Ctx) (*utils.Pagination, []companyRegistration.CompanyExpense, error) {
	pagination, expenses, err := utils.Paginate(c, r.db, companyRegistration.CompanyExpense{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, expenses, nil
}

func (r *CompanyRepositoryImpl) CreateCompanyLocation(location *companyRegistration.CompanyLocation) error {
	return r.db.Create(location).Error
}

func (r *CompanyRepositoryImpl) GetAllCompanyLocations(companyId string) ([]companyRegistration.CompanyLocation, error) {
	var locations []companyRegistration.CompanyLocation
	if err := r.db.Preload("Company").Where("company_id = ?", companyId).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *CompanyRepositoryImpl) GetLocationByCompanyId(id, companyId string) (*companyRegistration.CompanyLocation, error) {
	var location companyRegistration.CompanyLocation
	result := r.db.Preload("Company").Where("company_id = ? AND id = ?", companyId).First(&location)
	if result.Error != nil {
		return nil, result.Error
	}
	return &location, nil
}

func (r *CompanyRepositoryImpl) FindLocationById(Id string) (*companyRegistration.CompanyLocation, error) {
	var location companyRegistration.CompanyLocation
	result := r.db.Preload("Company").Where("id = ?", Id).First(&location)
	if result.Error != nil {
		return nil, result.Error
	}
	return &location, nil
}

func (r *CompanyRepositoryImpl) UpdateCompanyLocation(location *companyRegistration.CompanyLocation) error {
	return r.db.Save(location).Error
}

// Delete Company Location by ID
func (r *CompanyRepositoryImpl) DeleteLocationByID(id string) error {
	if err := r.db.Delete(&companyRegistration.CompanyLocation{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
