package repository

import (
	"car-bond/internals/middleware"
	"car-bond/internals/models/customerRegistration"
	"car-bond/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(customer *customerRegistration.Customer) error
	GetPaginatedCustomers(c *fiber.Ctx) (*utils.Pagination, []customerRegistration.Customer, error)
	GetCustomerAddresses(customerID uint) ([]customerRegistration.CustomerAddress, error)
	GetCustomerContacts(customerID uint) ([]customerRegistration.CustomerContact, error)
	// UpdateCustomer(customer *customerRegistration.Customer) error
	UpdateCustomer(id string, updates map[string]interface{}) error
	GetCustomerByID(id string) (customerRegistration.Customer, error)
	DeleteByID(id string) error
	SearchPaginatedCustomers(c *fiber.Ctx) (*utils.Pagination, []customerRegistration.Customer, error)

	// Contact
	CreateCustomerContact(address *customerRegistration.CustomerContact) error
	GetCustomerContactsByCustomerId(customerId string) ([]customerRegistration.CustomerContact, error)
	GetPaginatedContacts(c *fiber.Ctx, customerId string) (*utils.Pagination, []customerRegistration.CustomerContact, error)
	GetCustomerContactByIdAndCustomerId(id, customerId string) (*customerRegistration.CustomerContact, error)
	GetCustomerContactById(id string) (*customerRegistration.CustomerContact, error)
	UpdateCustomerContact(contact *customerRegistration.CustomerContact) error
	DeleteCustomerContactById(id, customerId string) error

	// Address
	CreateCustomerAddress(address *customerRegistration.CustomerAddress) error
	GetPaginatedAddresses(c *fiber.Ctx, companyId string) (*utils.Pagination, []customerRegistration.CustomerAddress, error)
	GetCustomerAddressesByCustomerId(customerId string) ([]customerRegistration.CustomerContact, error)
	GetCustomerAddressByIdAndCustomerId(id, customerId string) (*customerRegistration.CustomerAddress, error)
	GetCustomerAddressById(id string) (*customerRegistration.CustomerAddress, error)
	UpdateCustomerAddress(address *customerRegistration.CustomerAddress) error
	DeleteCustomerAddressById(id, customerId string) error
}

type CustomerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryImpl{db: db}
}

func (r *CustomerRepositoryImpl) CreateCustomer(customer *customerRegistration.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepositoryImpl) GetCustomerAddresses(customerID uint) ([]customerRegistration.CustomerAddress, error) {
	var addresses []customerRegistration.CustomerAddress
	err := r.db.Where("customer_id = ?", customerID).Find(&addresses).Error
	return addresses, err
}

// ====================

func (r *CustomerRepositoryImpl) GetCustomerContacts(customerID uint) ([]customerRegistration.CustomerContact, error) {
	var contacts []customerRegistration.CustomerContact
	err := r.db.Where("customer_id = ?", customerID).Find(&contacts).Error
	return contacts, err
}

// ==============

func (r *CustomerRepositoryImpl) GetPaginatedCustomers(c *fiber.Ctx) (*utils.Pagination, []customerRegistration.Customer, error) {
	_, companyID, err := middleware.GetUserAndCompanyFromSession(c)
	if err != nil {
		return nil, nil, err
	}

	// Scope the query by company_id
	query := r.db.Where("company_id = ?", companyID)

	pagination, customers, err := utils.Paginate(c, query, customerRegistration.Customer{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, customers, nil
}

func (r *CustomerRepositoryImpl) GetCustomerByID(id string) (customerRegistration.Customer, error) {
	var customer customerRegistration.Customer
	err := r.db.First(&customer, "id = ?", id).Error
	return customer, err
}

func (r *CustomerRepositoryImpl) UpdateCustomer(id string, updates map[string]interface{}) error {
	return r.db.Model(&customerRegistration.Customer{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a customer by ID
func (r *CustomerRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&customerRegistration.Customer{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CreateCustomerContact creates a new customer contact in the database
func (r *CustomerRepositoryImpl) CreateCustomerContact(address *customerRegistration.CustomerContact) error {
	return r.db.Create(address).Error
}

func (r *CustomerRepositoryImpl) GetCustomerContactsByCustomerId(customerId string) ([]customerRegistration.CustomerContact, error) {
	var contacts []customerRegistration.CustomerContact
	if err := r.db.Where("customer_id = ?", customerId).Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *CustomerRepositoryImpl) GetPaginatedContacts(c *fiber.Ctx, companyId string) (*utils.Pagination, []customerRegistration.CustomerContact, error) {
	pagination, contacts, err := utils.Paginate(c, r.db.Preload("Customer").Where("company_id = ?", companyId), customerRegistration.CustomerContact{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, contacts, nil
}

func (r *CustomerRepositoryImpl) GetCustomerContactByIdAndCustomerId(id, customerId string) (*customerRegistration.CustomerContact, error) {
	var contact customerRegistration.CustomerContact
	result := r.db.Preload("Customer").Where("id = ? AND customer_id = ?", id, customerId).First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (r *CustomerRepositoryImpl) GetCustomerContactById(id string) (*customerRegistration.CustomerContact, error) {
	var contact customerRegistration.CustomerContact
	result := r.db.Preload("Customer").Where("id = ?", id).First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (r *CustomerRepositoryImpl) UpdateCustomerContact(contact *customerRegistration.CustomerContact) error {
	return r.db.Save(contact).Error
}

// Delete Company contact by ID
func (r *CustomerRepositoryImpl) DeleteCustomerContactById(id, customerId string) error {
	if err := r.db.Delete(&customerRegistration.CustomerContact{}, "id = ? AND customer_id = ?", id, customerId).Error; err != nil {
		return err
	}
	return nil
}

func (r *CustomerRepositoryImpl) GetPaginatedAddresses(c *fiber.Ctx, companyId string) (*utils.Pagination, []customerRegistration.CustomerAddress, error) {
	pagination, contacts, err := utils.Paginate(c, r.db.Preload("Customer").Where("company_id = ?", companyId), customerRegistration.CustomerAddress{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, contacts, nil
}

// CreateCustomerAddress creates a new customer address in the database
func (r *CustomerRepositoryImpl) CreateCustomerAddress(address *customerRegistration.CustomerAddress) error {
	return r.db.Create(address).Error
}

func (r *CustomerRepositoryImpl) GetCustomerAddressesByCustomerId(customerId string) ([]customerRegistration.CustomerContact, error) {
	var addresses []customerRegistration.CustomerContact
	err := r.db.Preload("Customer").Where("customer_id = ?", customerId).Find(&addresses).Error
	return addresses, err
}

func (r *CustomerRepositoryImpl) GetCustomerAddressByIdAndCustomerId(id, customerId string) (*customerRegistration.CustomerAddress, error) {
	var contact customerRegistration.CustomerAddress
	result := r.db.Preload("Customer").Where("id = ? AND customer_id = ?", id, customerId).First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (r *CustomerRepositoryImpl) GetCustomerAddressById(id string) (*customerRegistration.CustomerAddress, error) {
	var address customerRegistration.CustomerAddress
	result := r.db.Preload("Customer").Where("id = ?", id).First(&address)
	if result.Error != nil {
		return nil, result.Error
	}
	return &address, nil
}

func (r *CustomerRepositoryImpl) UpdateCustomerAddress(address *customerRegistration.CustomerAddress) error {
	return r.db.Save(address).Error
}

// Delete Company contact by ID
func (r *CustomerRepositoryImpl) DeleteCustomerAddressById(id, customerId string) error {
	if err := r.db.Delete(&customerRegistration.CustomerAddress{}, "id = ? AND customer_id = ?", id, customerId).Error; err != nil {
		return err
	}
	return nil
}

func (r *CustomerRepositoryImpl) SearchPaginatedCustomers(c *fiber.Ctx) (*utils.Pagination, []customerRegistration.Customer, error) {
	// Get query parameters from request
	surname := c.Query("surname")
	firstname := c.Query("firstname")
	gender := c.Query("gender")
	nationality := c.Query("nationality")

	// Start building the query
	query := r.db.Model(&customerRegistration.Customer{})

	_, companyID, err := middleware.GetUserAndCompanyFromSession(c)
	if err != nil {
		return nil, nil, err
	}
	// Scope the query by company_id
	query = query.Where("company_id = ?", companyID)

	// Apply filters based on provided parameters
	if surname != "" {
		query = query.Where("LOWER(surname) LIKE LOWER(?)", "%"+surname+"%")
	}
	if firstname != "" {
		query = query.Where("LOWER(othername) LIKE LOWER(?)", "%"+firstname+"%")
	}
	if gender != "" {
		query = query.Where("LOWER(gender) = LOWER(?)", gender)
	}
	if nationality != "" {
		query = query.Where("LOWER(nationality) = LOWER(?)", nationality)
	}

	// Call the pagination helper
	pagination, customers, err := utils.Paginate(c, query, customerRegistration.Customer{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, customers, nil
}
