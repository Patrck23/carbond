package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"car-bond/internals/config"
	"car-bond/internals/models/customerRegistration"
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/models/carRegistration"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() {
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		fmt.Println("Error parsing str to int")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", 
		config.Config("DB_HOST"), 
		config.Config("DB_USER"), 
		config.Config("DB_PASS"), 
		config.Config("DB_NAME"), 
		port, 
		config.Config("DB_SSLMODE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	// --- Customer --- //
	db.AutoMigrate(&customerRegistration.CustomerAddress{})
	db.AutoMigrate(&customerRegistration.CustomerContact{})
	db.AutoMigrate(&customerRegistration.Customer{})
	// --- Company --- //
	db.AutoMigrate(&companyRegistration.Company{})
	// --- Car --- //
	db.AutoMigrate(&carRegistration.Car{})
	db.AutoMigrate(&carRegistration.CarExpense{})
	DB = Dbinstance{
		Db: db,
	}
}
