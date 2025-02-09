package controllers

import (
	"car-bond/internals/models/customerRegistration"
	"car-bond/internals/utils"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"path/filepath"

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

type CustomerController struct {
	repo CustomerRepository
}

func NewCustomerController(repo CustomerRepository) *CustomerController {
	return &CustomerController{repo: repo}
}

// ============================================

func (r *CustomerRepositoryImpl) CreateCustomer(customer *customerRegistration.Customer) error {
	return r.db.Create(customer).Error
}

func (h *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	// Parse form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Convert age from string to uint
	ageStr := c.FormValue("age")
	age, err := strconv.Atoi(ageStr) // Convert string to int
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid age provided",
			"data":    err.Error(),
		})
	}

	// Create a directory for storing files
	uploadDir := "./uploads/customer_files"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create upload directory",
			"data":    err.Error(),
		})
	}

	// Extract file if provided
	var filePath string
	if files, ok := form.File["upload_file"]; ok && len(files) > 0 {
		file := files[0]
		ext := strings.ToLower(filepath.Ext(file.Filename))

		// Validate the file type (photo or PDF)
		allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
		if !contains(allowedExtensions, ext) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid file type. Only JPG, JPEG, PNG, and PDF are allowed.",
			})
		}

		// Generate unique file name to avoid conflicts
		cleanFileName := strings.ReplaceAll(file.Filename, " ", "_")
		filePath = filepath.Join(uploadDir, cleanFileName)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to upload file",
				"data":    err.Error(),
			})
		}
	}

	// Create a new Customer instance
	customer := &customerRegistration.Customer{
		Surname:     c.FormValue("surname"),
		Firstname:   c.FormValue("firstname"),
		Othername:   c.FormValue("othername"),
		Gender:      c.FormValue("gender"),
		Nationality: c.FormValue("nationality"),
		Age:         uint(age),
		DOB:         c.FormValue("dob"),
		Telephone:   c.FormValue("telephone"),
		Email:       c.FormValue("email"),
		NIN:         c.FormValue("nin"),
		CreatedBy:   c.FormValue("created_by"),
		UpdatedBy:   c.FormValue("updated_by"),
		UploadFile:  filePath, // Store file path
	}

	// Attempt to create the customer record using the repository
	if err := h.repo.CreateCustomer(customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create customer",
			"data":    err.Error(),
		})
	}

	// Return the newly created customer record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Customer created successfully",
		"data":    customer,
	})
}

// =====================

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
	pagination, customers, err := utils.Paginate(c, r.db, customerRegistration.Customer{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, customers, nil
}

func (h *CustomerController) GetAllCustomers(c *fiber.Ctx) error {
	pagination, customers, err := h.repo.GetPaginatedCustomers(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customers",
			"data":    err.Error(),
		})
	}

	// Initialize a response slice to hold customers with their addresses and contacts
	var response []fiber.Map

	// Iterate over all customers to fetch associated customer addresses and contacts
	for _, customer := range customers {
		// Fetch customer addresses and contacts
		addresses, err := h.repo.GetCustomerAddresses(customer.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve customer ports for customer ID " + strconv.Itoa(int(customer.ID)),
				"data":    err.Error(),
			})
		}

		contacts, err := h.repo.GetCustomerContacts(customer.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve contacts for customer ID " + strconv.Itoa(int(customer.ID)),
				"data":    err.Error(),
			})
		}

		// Combine customer, ports, and contacts into a single response map
		response = append(response, fiber.Map{
			"customer":  customer,
			"addresses": addresses,
			"contacts":  contacts,
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "customers retrieved successfully",
		"data":    response,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =====================

func (r *CustomerRepositoryImpl) GetCustomerByID(id string) (customerRegistration.Customer, error) {
	var customer customerRegistration.Customer
	err := r.db.First(&customer, "id = ?", id).Error
	return customer, err
}

// GetSingleCustomer fetches a customer with its associated contacts and addresses from the database
func (h *CustomerController) GetSingleCustomer(c *fiber.Ctx) error {
	// Get the Customer ID from the route parameters
	id := c.Params("id")

	// Fetch the customer by ID
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Customer",
			"data":    err.Error(),
		})
	}

	// Fetch customer addresses associated with the customer
	addresses, err := h.repo.GetCustomerAddresses(customer.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customer addresses",
			"data":    err.Error(),
		})
	}

	// Fetch customer contacts associated with the customer
	contacts, err := h.repo.GetCustomerContacts(customer.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customer contacts",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"customer": customer,
		"address":  addresses,
		"contact":  contacts,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Customer and associated data retrieved successfully",
		"data":    response,
	})
}

// =======================

func (r *CustomerRepositoryImpl) UpdateCustomer(id string, updates map[string]interface{}) error {
	return r.db.Model(&customerRegistration.Customer{}).Where("id = ?", id).Updates(updates).Error
}

// Define the UpdateCustomer struct
type UpdateCustomerPayload struct {
	Surname     string `json:"surname"`
	Firstname   string `json:"firstname"`
	Othername   string `json:"othername"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Age         uint   `json:"age"`
	DOB         string `json:"dob"`
	Telephone   string `json:"telephone"`
	Email       string `json:"email"`
	NIN         string `json:"nin"`
	UpdatedBy   string `json:"updated_by"`
}

func (h *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
	// Get the customer ID from the route parameters
	id := c.Params("id")

	// Find the customer in the database
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customer",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateCustomerPayload struct
	var payload UpdateCustomerPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateCustomerPayload{} == payload) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	// Convert payload to a map for partial update
	updates := make(map[string]interface{})

	if payload.Surname != "" {
		updates["surname"] = payload.Surname
	}
	if payload.Firstname != "" {
		updates["firstname"] = payload.Firstname
	}
	if payload.Othername != "" {
		updates["othername"] = payload.Othername
	}
	if payload.Gender != "" {
		updates["gender"] = payload.Gender
	}
	if payload.Nationality != "" {
		updates["nationality"] = payload.Nationality
	}
	if payload.Age != 0 {
		updates["age"] = payload.Age
	}
	if payload.DOB != "" {
		updates["dob"] = payload.DOB
	}
	if payload.Telephone != "" {
		updates["telephone"] = payload.Telephone
	}
	if payload.Email != "" {
		updates["email"] = payload.Email
	}
	if payload.NIN != "" {
		updates["nin"] = payload.NIN
	}
	if payload.UpdatedBy != "" {
		updates["updated_by"] = payload.UpdatedBy
	}

	// Handle file upload
	file, err := c.FormFile("upload_file")
	if err == nil { // If a new file is uploaded
		uploadDir := "./uploads/customer_files/"

		// Ensure the directory exists
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create upload directory",
				"data":    err.Error(),
			})
		}

		// Replace spaces in file name with underscores
		safeFileName := strings.ReplaceAll(file.Filename, " ", "_")
		uploadPath := fmt.Sprintf("%s%s", uploadDir, safeFileName)

		// **Delete old file if it exists**
		if customer.UploadFile != "" {
			oldFilePath := customer.UploadFile
			if _, err := os.Stat(oldFilePath); err == nil {
				if err := os.Remove(oldFilePath); err != nil {
					return c.Status(500).JSON(fiber.Map{
						"status":  "error",
						"message": "Failed to delete existing file",
						"data":    err.Error(),
					})
				}
			}
		}

		// Save the new file
		if err := c.SaveFile(file, uploadPath); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to upload file",
				"data":    err.Error(),
			})
		}

		// Update the file path in the database
		updates["upload_file"] = uploadPath
	}

	// Update the customer in the database
	if err := h.repo.UpdateCustomer(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update customer",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Customer updated successfully",
		"data":    updates,
	})
}

// ======================

// DeleteByID deletes a customer by ID
func (r *CustomerRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&customerRegistration.Customer{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCustomerByID deletes a Customer by its ID
func (h *CustomerController) DeleteCustomerByID(c *fiber.Ctx) error {
	// Get the Customer ID from the route parameters
	id := c.Params("id")

	// Find the Customer in the database
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find Customer",
			"data":    err.Error(),
		})
	}

	// Delete the Customer
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Customer",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Customer deleted successfully",
		"data":    customer,
	})
}

// =======================

func (h *CustomerController) FetchCustomerUpload(c *fiber.Ctx) error {
	// Get the customer ID from the route parameters
	id := c.Params("id")

	// Find the customer in the database
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customer",
			"data":    err.Error(),
		})
	}

	// Check if the upload path exists
	if customer.UploadFile == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No file uploaded for this customer",
		})
	}

	// Serve the file from the upload path
	filePath := customer.UploadFile
	return c.SendFile(filePath)
}

// // =================================================================

// CreateCustomerContact creates a new customer contact in the database
func (r *CustomerRepositoryImpl) CreateCustomerContact(address *customerRegistration.CustomerContact) error {
	return r.db.Create(address).Error
}

// CreateCustomerContact handles the creation of a customer contact
func (h *CustomerController) CreateCustomerContact(c *fiber.Ctx) error {
	// Parse the request body into a customerContact struct
	customerContact := new(customerRegistration.CustomerContact)
	if err := c.BodyParser(customerContact); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the customer address in the database
	if err := h.repo.CreateCustomerContact(customerContact); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create customer contact",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "customer address created successfully",
		"data":    customerContact,
	})
}

// ========================

func (r *CustomerRepositoryImpl) GetCustomerContactsByCustomerId(customerId string) ([]customerRegistration.CustomerContact, error) {
	var contacts []customerRegistration.CustomerContact
	if err := r.db.Where("customer_id = ?", customerId).Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (h *CustomerController) GetCustomerContactsByCustomerId(c *fiber.Ctx) error {
	// Retrieve customerId from the request parameters
	customerId := c.Params("id")

	// Fetch contacts using the repository
	contacts, err := h.repo.GetCustomerContactsByCustomerId(customerId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer contacts not found for the specified date",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch customer contacts",
			"error":   err.Error(),
		})
	}

	// Return the fetched contacts
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "customer contacts fetched successfully",
		"data":    contacts,
	})
}

// ======================

func (r *CustomerRepositoryImpl) GetPaginatedContacts(c *fiber.Ctx, companyId string) (*utils.Pagination, []customerRegistration.CustomerContact, error) {
	pagination, contacts, err := utils.Paginate(c, r.db.Preload("Customer").Where("company_id = ?", companyId), customerRegistration.CustomerContact{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, contacts, nil
}

func (h *CustomerController) GetCustomerContactsByCompanyId(c *fiber.Ctx) error {
	companyId := c.Params("companyId")

	// Fetch paginated contacts using the repository
	pagination, contacts, err := h.repo.GetPaginatedContacts(c, companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve contacts",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "contacts and associated data retrieved successfully",
		"data":    contacts,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ======================

// Get Customer Contact by ID

func (r *CustomerRepositoryImpl) GetCustomerContactByIdAndCustomerId(id, customerId string) (*customerRegistration.CustomerContact, error) {
	var contact customerRegistration.CustomerContact
	result := r.db.Preload("Customer").Where("id = ? AND customer_id = ?", id, customerId).First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (h *CustomerController) GetCustomerContactById(c *fiber.Ctx) error {
	// Retrieve the contact ID and customer ID from the request parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Fetch the company contact from the repository
	contact, err := h.repo.GetCustomerContactByIdAndCustomerId(id, customerId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Contact not found for the specified customer",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the contact",
			"error":   err.Error(),
		})
	}

	// Return the fetched contact
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Contact fetched successfully",
		"data":    contact,
	})
}

// ==========================

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

func (h *CustomerController) UpdateCustomerContact(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCustomerContactInput struct {
		CustomerID         uint   `json:"customer_id" validate:"required"`
		ContactType        string `json:"contact_type" validate:"required"`
		ContactInformation string `json:"contact_information" validate:"required"`
		UpdatedBy          string `json:"updated_by" validate:"required"`
	}

	// Parse the expense ID from the request parameters
	expenseID := c.Params("id")

	// Parse and validate the request body
	var input UpdateCustomerContactInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Use a validation library to validate the input
	if validationErr := utils.ValidateStruct(input); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  validationErr,
		})
	}

	// Fetch the expense record using the repository
	expense, err := h.repo.GetCustomerContactById(expenseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch expense",
			"error":   err.Error(),
		})
	}

	// Update the expense fields
	expense.CustomerID = input.CustomerID
	expense.ContactType = input.ContactType
	expense.ContactInformation = input.ContactInformation
	expense.UpdatedBy = input.UpdatedBy

	// Save the updated expense using the repository
	if err := h.repo.UpdateCustomerContact(expense); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update expense",
			"error":   err.Error(),
		})
	}

	// Return the updated expense
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

// ==========================

// Delete Company contact by ID
func (r *CustomerRepositoryImpl) DeleteCustomerContactById(id, customerId string) error {
	if err := r.db.Delete(&customerRegistration.CustomerContact{}, "id = ? AND customer_id = ?", id, customerId).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCompanyByID deletes a company by its ID
func (h *CustomerController) DeleteCustomerContactById(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Find the contact in the database
	contact, err := h.repo.GetCustomerContactById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer contact not found or does not belong to the specified customer",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find contact",
			"data":    err.Error(),
		})
	}

	// Delete the contact
	if err := h.repo.DeleteCustomerContactById(id, customerId); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete contact",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "contact deleted successfully",
		"data":    contact,
	})
}

// =================================================================

func (r *CustomerRepositoryImpl) GetPaginatedAddresses(c *fiber.Ctx, companyId string) (*utils.Pagination, []customerRegistration.CustomerAddress, error) {
	pagination, contacts, err := utils.Paginate(c, r.db.Preload("Customer").Where("company_id = ?", companyId), customerRegistration.CustomerAddress{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, contacts, nil
}

func (h *CustomerController) GetCustomerAddressesByCompanyId(c *fiber.Ctx) error {
	companyId := c.Params("companyId")

	// Fetch paginated contacts using the repository
	pagination, contacts, err := h.repo.GetPaginatedAddresses(c, companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve addresses",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "addresses and associated data retrieved successfully",
		"data":    contacts,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ================================

// CreateCustomerAddress creates a new customer address in the database
func (r *CustomerRepositoryImpl) CreateCustomerAddress(address *customerRegistration.CustomerAddress) error {
	return r.db.Create(address).Error
}

// CreateCustomerAddress handles the creation of a customer address
func (h *CustomerController) CreateCustomerAddress(c *fiber.Ctx) error {
	// Parse the request body into a customerAddress struct
	customerAddress := new(customerRegistration.CustomerAddress)
	if err := c.BodyParser(customerAddress); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the customer address in the database
	if err := h.repo.CreateCustomerAddress(customerAddress); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create customer address",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "customer address created successfully",
		"data":    customerAddress,
	})
}

// ===================================

func (r *CustomerRepositoryImpl) GetCustomerAddressesByCustomerId(customerId string) ([]customerRegistration.CustomerContact, error) {
	var addresses []customerRegistration.CustomerContact
	err := r.db.Preload("Customer").Where("customer_id = ?", customerId).Find(&addresses).Error
	return addresses, err
}

func (h *CustomerController) GetCustomerAddressesByCustomerId(c *fiber.Ctx) error {
	customerId := c.Params("id")

	// Fetch company expenses using the service layer
	expenses, err := h.repo.GetCustomerAddressesByCustomerId(customerId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch customer addresses",
			"error":   err.Error(),
		})
	}

	// If no expenses found
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Customer addresses not found for the specified filters",
		})
	}

	// Return the fetched company expenses
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Customer addresses fetched successfully",
		"data":    expenses,
	})
}

// ==============================

// // Get Customer Address by ID

func (r *CustomerRepositoryImpl) GetCustomerAddressByIdAndCustomerId(id, customerId string) (*customerRegistration.CustomerAddress, error) {
	var contact customerRegistration.CustomerAddress
	result := r.db.Preload("Customer").Where("id = ? AND customer_id = ?", id, customerId).First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (h *CustomerController) GetCustomerAddressById(c *fiber.Ctx) error {
	// Retrieve the contact ID and customer ID from the request parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Fetch the company contact from the repository
	contact, err := h.repo.GetCustomerAddressByIdAndCustomerId(id, customerId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Address not found for the specified customer",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the address",
			"error":   err.Error(),
		})
	}

	// Return the fetched contact
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Address fetched successfully",
		"data":    contact,
	})
}

// ========================================

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

func (h *CustomerController) UpdateCustomerAddress(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCustomerAddressInput struct {
		CustomerID uint   `json:"customer_id" validate:"required"`
		District   string `json:"district" validate:"required"`
		Subcounty  string `json:"subcounty" validate:"required"`
		Parish     string `json:"parish" validate:"required"`
		Village    string `json:"village" validate:"required"`
		UpdatedBy  string `json:"updated_by" validate:"required"`
	}

	// Parse the expense ID from the request parameters
	expenseID := c.Params("id")

	// Parse and validate the request body
	var input UpdateCustomerAddressInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Use a validation library to validate the input
	if validationErr := utils.ValidateStruct(input); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  validationErr,
		})
	}

	// Fetch the expense record using the repository
	expense, err := h.repo.GetCustomerAddressById(expenseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch expense",
			"error":   err.Error(),
		})
	}

	// Update the expense fields
	expense.CustomerID = input.CustomerID
	expense.District = input.District
	expense.Subcounty = input.Subcounty
	expense.Parish = input.Parish
	expense.Village = input.Village
	expense.UpdatedBy = input.UpdatedBy

	// Save the updated expense using the repository
	if err := h.repo.UpdateCustomerAddress(expense); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update expense",
			"error":   err.Error(),
		})
	}

	// Return the updated expense
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

// ==========================

// Delete Company contact by ID
func (r *CustomerRepositoryImpl) DeleteCustomerAddressById(id, customerId string) error {
	if err := r.db.Delete(&customerRegistration.CustomerAddress{}, "id = ? AND customer_id = ?", id, customerId).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCompanyByID deletes a company by its ID
func (h *CustomerController) DeleteCustomerAddressById(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Find the contact in the database
	contact, err := h.repo.GetCustomerAddressById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer contact not found or does not belong to the specified customer",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find contact",
			"data":    err.Error(),
		})
	}

	// Delete the contact
	if err := h.repo.DeleteCustomerAddressById(id, customerId); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete contact",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "contact deleted successfully",
		"data":    contact,
	})
}

// ========================================

// // Delete Customer address by ID

// func DeleteCustomerAddressById(c *fiber.Ctx) error {
// 	// Initialize database instance
// 	db := database.DB.Db

// 	// Parse parameters
// 	id := c.Params("id")
// 	customerId := c.Params("customerId")

// 	// Check if the customer exists
// 	var customer customerRegistration.Customer
// 	if err := db.First(&customer, customerId).Error; err != nil {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer not found"})
// 	}

// 	// Check if the customer address exists and belongs to the customer
// 	var customerAddress customerRegistration.CustomerAddress
// 	if err := db.Where("id = ? AND customer_id = ?", id, customerId).First(&customerAddress).Error; err != nil {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer address not found or does not belong to the specified customer"})
// 	}

// 	// Delete the customer address
// 	if err := db.Delete(&customerAddress).Error; err != nil {
// 		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete customer address", "data": err.Error()})
// 	}

// 	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Customer address deleted successfully"})
// }

// ======================

func (r *CustomerRepositoryImpl) SearchPaginatedCustomers(c *fiber.Ctx) (*utils.Pagination, []customerRegistration.Customer, error) {
	// Get query parameters from request
	surname := c.Query("surname")
	firstname := c.Query("firstname")
	gender := c.Query("gender")
	nationality := c.Query("nationality")

	// Start building the query
	query := r.db.Model(&customerRegistration.Customer{})

	// Apply filters based on provided parameters
	if surname != "" {
		query = query.Where("surname LIKE ?", "%"+surname+"%")
	}
	if firstname != "" {
		query = query.Where("othername LIKE ?", "%"+firstname+"%")
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if nationality != "" {
		query = query.Where("nationality = ?", nationality)
	}

	// Call the pagination helper
	pagination, customers, err := utils.Paginate(c, query, customerRegistration.Customer{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, customers, nil
}

func (h *CustomerController) SearchCustomers(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, customers, err := h.repo.SearchPaginatedCustomers(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customers",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Customers retrieved successfully",
		"pagination": pagination,
		"data":       customers,
	})
}
