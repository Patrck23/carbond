package controllers

import (
	"car-bond/internals/database"
	"car-bond/internals/models/saleRegistration"

	"github.com/gofiber/fiber/v2"
)

// Create a car sale
func CreateCarSale(c *fiber.Ctx) error {
	db := database.DB.Db
	var sale saleRegistration.Sale

	if err := c.BodyParser(&sale); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if err := db.Create(&sale).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to create sale", "data": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Sale created successfully", "data": sale})
}

// Get all car sales
func GetAllCarSales(c *fiber.Ctx) error {
	db := database.DB.Db
	var sales []saleRegistration.Sale

	db.Find(&sales)
	if len(sales) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No sales found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Sales found", "data": sales})
}

// Get a single car sale by ID
func GetCarSale(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var sale saleRegistration.Sale

	db.First(&sale, "id = ?", id)
	if sale.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Sale not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Sale found", "data": sale})
}

// Update a car sale
func UpdateCarSale(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var sale saleRegistration.Sale

	db.First(&sale, "id = ?", id)
	if sale.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Sale not found"})
	}

	userID := c.Locals("user_id")
	isAdmin := c.Locals("is_admin").(bool)

	if sale.CreatedBy != userID && !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "You don't have the required permissions"})
	}

	var updatedSale saleRegistration.Sale
	if err := c.BodyParser(&updatedSale); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	db.Model(&sale).Updates(updatedSale)
	return c.JSON(fiber.Map{"status": "success", "message": "Sale updated successfully", "data": sale})
}

// Delete a car sale
func DeleteCarSale(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var sale saleRegistration.Sale

	db.First(&sale, "id = ?", id)
	if sale.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Sale not found"})
	}

	if err := db.Delete(&sale, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to delete sale", "data": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Sale deleted successfully"})
}

// Common search functions for filtering sales
func searchSales(c *fiber.Ctx, condition string, args ...interface{}) error {
	db := database.DB.Db
	var sales []saleRegistration.Sale

	db.Where(condition, args...).Find(&sales)
	if len(sales) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No results found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Results found", "data": sales})
}

// Search specific criteria
func SearchByCriteria(c *fiber.Ctx) error {
	criteria := c.Query("criteria")
	query := c.Query("query")

	switch criteria {
	case "unfinished_transactions":
		return searchSales(c, "transaction_status = ?", "Unfinished")
	case "customer_details":
		return searchSales(c, "customer_name LIKE ? OR customer_id LIKE ?", "%"+query+"%", "%"+query+"%")
	case "car_details":
		return searchSales(c, "car_model LIKE ? OR car_registration.registration_number LIKE ?", "%"+query+"%", "%"+query+"%")
	case "car_brand":
		return searchSales(c, "car_brand LIKE ?", "%"+query+"%")
	case "total_amount":
		return searchSales(c, "total_amount = ?", query)
	case "sale_date":
		return searchSales(c, "sale_date = ?", query)
	case "company":
		return searchSales(c, "company = ?", query)
	case "full_payment":
		return searchSales(c, "is_full_payment = ?", true)
	case "partial_payment":
		return searchSales(c, "is_full_payment = ?", false)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid search criteria"})
	}
}

//Create an invoice for a customer

func CreateInvoice(c *fiber.Ctx) error {
	db := database.DB.Db
	var sale saleRegistration.SalePayment

	if err := c.BodyParser(&sale); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if err := db.Create(&sale).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to create sale", "data": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Invoice created successfully", "data": sale})
}

//Get all invoices

func GetAllInvoices(c *fiber.Ctx) error {
	db := database.DB.Db
	var sales []saleRegistration.SalePayment

	db.Find(&sales)
	if len(sales) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No sales found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Sales found", "data": sales})
}

//Update an invoice

func UpdateInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var sale saleRegistration.SalePayment

	db.First(&sale, "id = ?", id)
	if sale.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Sale not found"})
	}

	userID := c.Locals("user_id")
	isAdmin := c.Locals("is_admin").(bool)

	if sale.CreatedBy != userID && !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "You don't have the required permissions"})
	}

	var updatedSale saleRegistration.SalePayment
	if err := c.BodyParser(&updatedSale); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}
	db.Model(&sale).Updates(updatedSale)

	return c.JSON(fiber.Map{"status": "success", "message": "Sale updated successfully", "data": sale})
}

//Delete invoice by ID

func DeleteInvoiceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var sale saleRegistration.SalePayment

	db.First(&sale, "id = ?", id)
	if sale.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Sale not found"})
	}

	if err := db.Delete(&sale, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to delete sale", "data": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Sale deleted successfully"})
}

// Get all Payments

func GetAllPayments(c *fiber.Ctx) error {
	db := database.DB.Db
	var payments []saleRegistration.SalePayment

	db.Find(&payments)
	if len(payments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No payments found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Payments found", "data": payments})
}

//Get payment by ID

func GetPaymentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var payment saleRegistration.SalePayment

	db.First(&payment, "id = ?", id)
	if payment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Payment not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Payment found", "data": payment})
}

// Update payment

func UpdatePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var payment saleRegistration.SalePayment

	db.First(&payment, "id = ?", id)
	if payment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Payment not found"})
	}

	userID := c.Locals("user_id")
	isAdmin := c.Locals("is_admin").(bool)

	if payment.CreatedBy != userID && !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "You don't have the required permissions"})
	}

	var updatedPayment saleRegistration.SalePayment
	if err := c.BodyParser(&updatedPayment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}
	db.Model(&payment).Updates(updatedPayment)
	return c.JSON(fiber.Map{"status": "success", "message": "Payment updated successfully", "data": payment})
}

//Delete payment By ID

func DeletePaymentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB.Db
	var payment saleRegistration.SalePayment

	db.First(&payment, "id = ?", id)
	if payment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Payment not found"})
	}

	if err := db.Delete(&payment, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to delete payment", "data": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Payment deleted successfully"})
}

// Get payment by ModeOfPayment

func GetPaymentByModeOfPayment(c *fiber.Ctx) error {
	mode := c.Params("mode")
	db := database.DB.Db
	var payments []saleRegistration.SalePayment

	db.Where("mode_of_payment = ?", mode).Find(&payments)
	if len(payments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No payments found for the given mode of payment"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Payments found", "data": payments})
}

//create payment

func CreatePayment(c *fiber.Ctx) error {
	db := database.DB.Db
	var payment saleRegistration.SalePayment

	if err := c.BodyParser(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if err := db.Create(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to create payment", "data": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Payment created successfully", "data": payment})
}
