package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"psi-src/internals/config"
	"psi-src/internals/models/customerRegistration"

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.Config("DB_HOST"), config.Config("DB_USER"), config.Config("DB_PASS"), config.Config("DB_NAME"), port, config.Config("DB_SSLMODE"))
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
	db.AutoMigrate(&customerRegistration.ClientAddress{})
	db.AutoMigrate(&customerRegistration.ClientAllergy{})
	db.AutoMigrate(&customerRegistration.ClientChronicCondition{})
	db.AutoMigrate(&customerRegistration.ClientContact{})
	db.AutoMigrate(&customerRegistration.Client{})
	db.AutoMigrate(&customerRegistration.ClientIdentifier{})
	db.AutoMigrate(&customerRegistration.ClientQueue{})
	db.AutoMigrate(&customerRegistration.ClientVisit{})
	DB = Dbinstance{
		Db: db,
	}
}
