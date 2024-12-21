package controllers

import (
	"car-bond/internals/database"
	"car-bond/internals/models/customerRegistration"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Create a customer
func CreateCustomer(c *fiber.Ctx) error {
	db := database.DB.Db
	customer := new(customerRegistration.Customer)
	// Store the body in the customer and return error if encountered
	err := c.BodyParser(customer)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	err = db.Create(&customer).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create customer", "data": err})
	}
	// Return the created customer
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "customer created", "data": customer})
}

// Get All customers from db
func GetAllCustomers(c *fiber.Ctx) error {
	db := database.DB.Db
	var customers []customerRegistration.Customer
	// find all customers in the database
	db.Find(&customers)
	// If no customer found, return an error
	if len(customers) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customers not found"})
	}
	// return customers
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "Customers Found", "data": customers})
}

// GetSingleCustomer from db
func GetSingleCustomer(c *fiber.Ctx) error {
	db := database.DB.Db
	// get id params
	id := c.Params("id")
	var customer customerRegistration.Customer
	// find single customer in the database by id
	db.Find(&customer, "id = ?", id)
	if customer.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer not found"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Customer Found", "data": customer})
}

// update a customer in db
func UpdateCustomer(c *fiber.Ctx) error {
	type updateCustomer struct {
		Surname     string `json:"surname"`
		Firstname   string `json:"firstname"`
		Othername   string `json:"othername"`
		Gender      string `json:"gender"`
		Nationality string `json:"nationality"`
		Age         uint   `json:"age"`
		DOB         string `gorm:"type:date" json:"dob"`
		Telephone   string `json:"telephone"`
		Email       string `json:"email"`
		NIN         string `json:"nin"`
		UpdatedBy   string `json:"updated_by"`
	}
	db := database.DB.Db
	var customer customerRegistration.Customer
	// get id params
	id := c.Params("id")
	// find single customer in the database by id
	db.Find(&customer, "id = ?", id)
	if customer.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer not found", "data": nil})
	}
	var updateCustomerData updateCustomer
	err := c.BodyParser(&updateCustomerData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	customer.Surname = updateCustomerData.Surname
	customer.Firstname = updateCustomerData.Firstname
	customer.Othername = updateCustomerData.Othername
	customer.Gender = updateCustomerData.Gender
	customer.Nationality = updateCustomerData.Nationality
	customer.Age = updateCustomerData.Age
	customer.DOB = updateCustomerData.DOB
	customer.Telephone = updateCustomerData.Telephone
	customer.Email = updateCustomerData.Email
	customer.NIN = updateCustomerData.NIN
	customer.UpdatedBy = updateCustomerData.UpdatedBy
	// Save the Changes
	db.Save(&customer)
	// Return the updated customer
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "customers Found", "data": customer})
}

// delete customer in db by ID
func DeleteCustomerByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var customer customerRegistration.Customer
	// get id params
	id := c.Params("id")
	// find single customer in the database by id
	db.Find(&customer, "id = ?", id)
	if customer.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer not found", "data": nil})
	}
	err := db.Delete(&customer, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete user", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Customer deleted"})
}

// =================================================================

// Get all customer contacts

func GetAllContacts(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	var contacts []customerRegistration.CustomerContact

	// Fetch all contacts with associated customers
	if err := db.Preload("Customer").Find(&contacts).Error; err != nil {
		// Handle database query error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch contacts",
		})
	}

	// Check if no contacts are found
	if len(contacts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No contacts found",
		})
	}

	// Return success response with contacts
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Contacts found",
		"data":    contacts,
	})
}

// Create a customer contact

func CreateCustomerContact(c *fiber.Ctx) error {
	db := database.DB.Db
	contact := new(customerRegistration.CustomerContact)
	err := c.BodyParser(contact)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&contact).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create customer contact", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Customer contact created successfully", "data": contact})
}

// Get customer contacts
func GetCustomerContacts(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	customerId := c.Params("customerId")

	var contacts []customerRegistration.CustomerContact

	// Fetch all contacts with associated customer details
	if err := db.Preload("Customer").Where("customer_id = ?", customerId).Find(&contacts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch customer contact",
		})
	}

	return c.Status(fiber.StatusOK).JSON(contacts)
}

// Get Customer Contact by ID

func GetCustomerContactById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve the contact ID and customer ID from the request parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Query the database for the contact by its ID and customer ID
	var contact customerRegistration.CustomerContact
	result := db.Preload("Customer").Where("id = ? AND customer_id = ?", id, customerId).First(&contact)

	// Handle potential database query errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Contact not found for the specified customer",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the contact",
			"error":   result.Error.Error(),
		})
	}

	// Return the fetched contact
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Contact fetched successfully",
		"data":    contact,
	})
}

// Update Customer Contact

func UpdateCustomerContact(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCustomerContactInput struct {
		ContactType        string `json:"contact_type"`
		ContactInformation string `json:"contact_information"`
		UpdatedBy          string `json:"updated_by"`
	}

	db := database.DB.Db
	contactID := c.Params("id")

	// Find the contact record by ID
	var contact customerRegistration.CustomerContact
	if err := db.First(&contact, contactID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Contact not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Database error", "error": err.Error()})
	}

	// Parse the request body into the input struct
	var input UpdateCustomerContactInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "error": err.Error()})
	}

	// Update the fields of the contact record
	contact.ContactType = input.ContactType
	contact.ContactInformation = input.ContactInformation
	contact.UpdatedBy = input.UpdatedBy

	// Save the changes to the database
	if err := db.Save(&contact).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update contact", "error": err.Error()})
	}

	// Return the updated contact
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "contact updated successfully",
		"data":    contact,
	})
}

// Delete Customer contact by ID

func DeleteCustomerContactById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Parse parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Check if the customer exists
	var customer customerRegistration.Customer
	if err := db.First(&customer, customerId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer not found"})
	}

	// Check if the customer contact exists and belongs to the customer
	var customerContact customerRegistration.CustomerContact
	if err := db.Where("id = ? AND customer_id = ?", id, customerId).First(&customerContact).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer contact not found or does not belong to the specified customer"})
	}

	// Delete the customer contact
	if err := db.Delete(&customerContact).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete customer contact", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Customer contact deleted successfully"})
}

// =================================================================

// Get all customer addresses

func GetAllAddresses(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	var addresses []customerRegistration.CustomerAddress

	// Fetch all addresses with associated customers
	if err := db.Preload("Customer").Find(&addresses).Error; err != nil {
		// Handle database query error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch addresses",
		})
	}

	// Check if no addresses are found
	if len(addresses) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No addresses found",
		})
	}

	// Return success response with addresses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Addresses found",
		"data":    addresses,
	})
}

// Create a customer address

func CreateCustomerAddress(c *fiber.Ctx) error {
	db := database.DB.Db
	address := new(customerRegistration.CustomerAddress)
	err := c.BodyParser(address)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&address).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create customer address", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Customer address created successfully", "data": address})
}

// Get customer addresses
func GetCustomerAddresses(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	customerId := c.Params("customerId")

	var addresses []customerRegistration.CustomerAddress

	// Fetch all addresses with associated customer details
	if err := db.Preload("Customer").Where("customer_id = ?", customerId).Find(&addresses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch customer address",
		})
	}

	return c.Status(fiber.StatusOK).JSON(addresses)
}

// Get Customer Address by ID

func GetCustomerAddressById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve the address ID and customer ID from the request parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Query the database for the address by its ID and customer ID
	var address customerRegistration.CustomerAddress
	result := db.Preload("Customer").Where("id = ? AND customer_id = ?", id, customerId).First(&address)

	// Handle potential database query errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Address not found for the specified customer",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the address",
			"error":   result.Error.Error(),
		})
	}

	// Return the fetched address
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Address fetched successfully",
		"data":    address,
	})
}

// Update Customer Address

func UpdateCustomerAddress(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCustomerAddressInput struct {
		District  string `json:"district"`
		Subcounty string `json:"subcounty"`
		Parish    string `json:"parish"`
		Village   string `json:"village"`
		UpdatedBy string `json:"updated_by"`
	}

	db := database.DB.Db
	addressID := c.Params("id")

	// Find the address record by ID
	var address customerRegistration.CustomerAddress
	if err := db.First(&address, addressID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Address not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Database error", "error": err.Error()})
	}

	// Parse the request body into the input struct
	var input UpdateCustomerAddressInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "error": err.Error()})
	}

	// Update the fields of the address record
	address.District = input.District
	address.Subcounty = input.Subcounty
	address.Parish = input.Parish
	address.Village = input.Village
	address.UpdatedBy = input.UpdatedBy

	// Save the changes to the database
	if err := db.Save(&address).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update address", "error": err.Error()})
	}

	// Return the updated address
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "address updated successfully",
		"data":    address,
	})
}

// Delete Customer address by ID

func DeleteCustomerAddressById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Parse parameters
	id := c.Params("id")
	customerId := c.Params("customerId")

	// Check if the customer exists
	var customer customerRegistration.Customer
	if err := db.First(&customer, customerId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer not found"})
	}

	// Check if the customer address exists and belongs to the customer
	var customerAddress customerRegistration.CustomerAddress
	if err := db.Where("id = ? AND customer_id = ?", id, customerId).First(&customerAddress).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Customer address not found or does not belong to the specified customer"})
	}

	// Delete the customer address
	if err := db.Delete(&customerAddress).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete customer address", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Customer address deleted successfully"})
}
