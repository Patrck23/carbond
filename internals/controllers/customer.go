package controllers

import (
	"car-bond/internals/models/customerRegistration"
	"car-bond/internals/repository"
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

type CustomerController struct {
	repo repository.CustomerRepository
}

func NewCustomerController(repo repository.CustomerRepository) *CustomerController {
	return &CustomerController{repo: repo}
}

// ============================================

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
