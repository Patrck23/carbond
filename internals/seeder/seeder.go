package seeder

import (
	"car-bond/internals/models/companyRegistration"
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
			Name:      "DAJIMA MOTORS (U) ltd",
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
}
