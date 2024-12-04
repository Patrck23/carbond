package controllers

import (
	"psi-src/internals/database"
	"psi-src/internals/models/customerRegistration"

	"github.com/gofiber/fiber/v2"
)

// Create a customer
func CreateCustomer(c *fiber.Ctx) error {
	db := database.DB.Db
	customer := new(registration.Customer)
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
	var customers []registration.Customer
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
	var customer registration.Customer
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
		Age         int    `json:"age"`
		DOB         string `gorm:"type:date" json:"dob"`
		Telephone   string `json:"telephone"`
		Email		string `json:"email"`
		NIN         string `json:"nin"`
		UpdatedBy   string `json:"updated_by"`
	}
	db := database.DB.Db
	var customer registration.Customer
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
	var customer registration.Customer
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
