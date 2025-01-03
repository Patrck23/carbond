package seeder

import (
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/models/metaData"
	"car-bond/internals/models/userRegistration"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func SeedDatabase(db *gorm.DB) {

	companies := []companyRegistration.Company{
		{
			Name:      "SHERAZ TRADING (u) ltd",
			StartDate: "1990-12-01",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "K.K. MADNA (J) ltd",
			StartDate: "1992-07-01",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the companies table already has data
	var companyCount int64
	db.Model(&companyRegistration.Company{}).Count(&companyCount)
	if companyCount == 0 {
		if err := db.Create(&companies).Error; err != nil {
			log.Fatalf("Failed to seed companies: %v", err)
		} else {
			log.Println("Company data seeded successfully")
		}
	} else {
		log.Println("Companies table already seeded, skipping...")
	}

	// Hashing password for users
	passwordHash, err := hashPassword("Admin123")
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	users := []userRegistration.User{
		{
			Surname:   "",
			Firstname: "John",
			Othername: "Doe",
			Gender:    "",
			Title:     "",
			Username:  "Admin",
			Email:     "admin@example.com",
			Password:  passwordHash,
			CompanyID: 1,
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	var userCount int64
	db.Model(&userRegistration.User{}).Count(&userCount)
	if userCount == 0 {
		if err := db.Create(&users).Error; err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		} else {
			log.Println("User data seeded successfully")
		}
	} else {
		log.Println("Users table already seeded, skipping...")
	}

	// Expense categories
	expenses := []metaData.ExpenseCategory{
		{
			Name:      "Auction Fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "VAT (Value Added Tax)",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Recycle Fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Carrier car fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Commission Fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Broker Commission Fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port Fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Inspection Fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Freight Cost",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Car Duty Fee",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Maintenance and Repair Fees",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Rent or lease and morgage",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Marketing and advertising",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Licences and permits",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Utility Bills",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Employ Salaries",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Insurance",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Technology and software",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Transport costs",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Dealership Amenities",
			Category:  "company",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the expense category table already has data
	var expenseCount int64
	db.Model(&metaData.ExpenseCategory{}).Count(&expenseCount)
	if expenseCount == 0 {
		if err := db.Create(&expenses).Error; err != nil {
			log.Fatalf("Failed to seed expense categories: %v", err)
		} else {
			log.Println("Expense category data seeded successfully")
		}
	} else {
		log.Println("Expense category table already seeded, skipping...")
	}

	unitsLength := []metaData.LeightUnit{
		{
			Name:      "Meter",
			Symbol:    "M",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Kilometer",
			Symbol:    "KM",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Centimeter",
			Symbol:    "CM",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the lengthunits table already has data
	var lenCount int64
	db.Model(&metaData.LeightUnit{}).Count(&lenCount)
	if lenCount == 0 {
		if err := db.Create(&unitsLength).Error; err != nil {
			log.Fatalf("Failed to seed Length Units: %v", err)
		} else {
			log.Println("Length units data seeded successfully")
		}
	} else {
		log.Println("Length units table already seeded, skipping...")
	}

	unitsWeight := []metaData.LeightUnit{
		{
			Name:      "Gram",
			Symbol:    "G",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Kilogram",
			Symbol:    "KG",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Tonnes",
			Symbol:    "Ton",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the weightunits table already has data
	var weiCount int64
	db.Model(&metaData.WeightUnit{}).Count(&weiCount)
	if weiCount == 0 {
		if err := db.Create(&unitsWeight).Error; err != nil {
			log.Fatalf("Failed to seed Weight Units: %v", err)
		} else {
			log.Println("Weight units data seeded successfully")
		}
	} else {
		log.Println("Weight units table already seeded, skipping...")
	}

}
