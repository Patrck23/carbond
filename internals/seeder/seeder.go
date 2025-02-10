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
			Name:      "SHERAZ TRADING (U) ltd",
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

	companyLocations := []companyRegistration.CompanyLocation{
		{
			CompanyID: 1,
			Address:   "",
			Telephone: "",
			Country:   "Uganda",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			CompanyID: 2,
			Address:   "",
			Telephone: "",
			Country:   "Japan",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the companies table already has data
	var companyLocationCount int64
	db.Model(&companyRegistration.CompanyLocation{}).Count(&companyCount)
	if companyLocationCount == 0 {
		if err := db.Create(&companyLocations).Error; err != nil {
			log.Fatalf("Failed to seed companies: %v", err)
		} else {
			log.Println("Company location data seeded successfully")
		}
	} else {
		log.Println("Company location table already seeded, skipping...")
	}

	groups := []userRegistration.Group{
		{
			Code:        "admin",
			Name:        "Admin Group",
			Description: "Group for administrative users",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "manager",
			Name:        "Manager Group",
			Description: "Group for manager users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "user",
			Name:        "User Group",
			Description: "Group for regular users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
	}

	// Seed Roles
	roles := []userRegistration.Role{
		// Admin Group Roles
		{
			Code:        "resource.admin",
			Name:        "Admin Role",
			Description: "Full resource permissions",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     1,
		},
		{
			Code:        "users.admin",
			Name:        "Admin Role",
			Description: "Full user permissions",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     1,
		},
		{
			Code:        "documents.admin",
			Name:        "Admin Role",
			Description: "Full document permissions",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     1,
		},
		// Manager Group Roles
		{
			Code:        "resource.read",
			Name:        "Manager Role",
			Description: "Read-only resource permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     2,
		},
		{
			Code:        "resource.write",
			Name:        "Manager Role",
			Description: "Write resource permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     2,
		},
		{
			Code:        "users.read",
			Name:        "Manager Role",
			Description: "Read-only user permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     2,
		},
		{
			Code:        "users.write",
			Name:        "Manager Role",
			Description: "Write user permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     2,
		},
		{
			Code:        "documents.read",
			Name:        "Manager Role",
			Description: "Read-only document permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     2,
		},
		{
			Code:        "documents.write",
			Name:        "Manager Role",
			Description: "Write document permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     2,
		},
		// User Group Roles
		{
			Code:        "resource.write",
			Name:        "User Role",
			Description: "Write resource permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     3,
		},
		{
			Code:        "users.write",
			Name:        "User Role",
			Description: "Write user permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     3,
		},
		{
			Code:        "documents.write",
			Name:        "User Role",
			Description: "Write document permissions",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
			GroupID:     3,
		},
	}

	// Seed Resources
	resources := []userRegistration.Resource{
		{
			Code:        "resource.*",
			Name:        "Resources",
			Description: "Access everything",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "resource.all",
			Name:        "Read all Resouces",
			Description: "Access to the dashboard",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "resource.my",
			Name:        "personal resource",
			Description: "Access to settings page",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "users.*",
			Name:        "All users",
			Description: "Access to all users",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "users.all",
			Name:        "Access to all users",
			Description: "Read users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "users.my",
			Name:        "Just the user",
			Description: "Edit users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "documents.*",
			Name:        "All document access",
			Description: "Access to all documents",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "documents.all",
			Name:        "Read documents",
			Description: "Read all documents",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "documents.my",
			Name:        "View user documents",
			Description: "Edit my documents",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
	}

	// Seed RoleResourcePermissions
	roleResourcePermissions := []userRegistration.RoleResourcePermission{
		{
			RoleCode:     "resource.admin",
			ResourceCode: "resource.*",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: true, D: true},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "resource.reader",
			ResourceCode: "resource.all",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: false, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "resource.writer",
			ResourceCode: "resource.my",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "users.admin",
			ResourceCode: "users.*",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: true, D: true},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "users.reader",
			ResourceCode: "users.all",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: false, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "users.writer",
			ResourceCode: "users.my",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "documents.admin",
			ResourceCode: "documents.*",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: true, D: true},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "documents.reader",
			ResourceCode: "documents.all",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: false, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "documents.writer",
			ResourceCode: "documents.my",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Seed RoleWildCardPermissions
	roleWildCardPermissions := []userRegistration.RoleWildCardPermission{
		{
			RoleCode:     "resource.admin",
			ResourceCode: "resource.*",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: true, D: true},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		}, {
			RoleCode:     "resource.writer",
			ResourceCode: "resource.my",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "resource.reader",
			ResourceCode: "resource.all",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: false, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: true, D: true},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the tables have data before seeding
	// Check Group
	var groupCount int64
	db.Model(&userRegistration.Group{}).Count(&groupCount)
	if groupCount == 0 {
		if err := db.Create(&groups).Error; err != nil {
			log.Fatalf("Failed to seed Group data: %v", err)
		} else {
			log.Println("Group data seeded successfully")
		}
	}

	// Check Role
	var roleCount int64
	db.Model(&userRegistration.Role{}).Count(&roleCount)
	if roleCount == 0 {
		if err := db.Create(&roles).Error; err != nil {
			log.Fatalf("Failed to seed Role data: %v", err)
		} else {
			log.Println("Role data seeded successfully")
		}
	}

	// Check Resource
	var resourceCount int64
	db.Model(&userRegistration.Resource{}).Count(&resourceCount)
	if resourceCount == 0 {
		if err := db.Create(&resources).Error; err != nil {
			log.Fatalf("Failed to seed Resource data: %v", err)
		} else {
			log.Println("Resource data seeded successfully")
		}
	}

	// Check RoleResourcePermission
	var roleResourcePermissionCount int64
	db.Model(&userRegistration.RoleResourcePermission{}).Count(&roleResourcePermissionCount)
	if roleResourcePermissionCount == 0 {
		if err := db.Create(&roleResourcePermissions).Error; err != nil {
			log.Fatalf("Failed to seed RoleResourcePermission data: %v", err)
		} else {
			log.Println("RoleResourcePermission data seeded successfully")
		}
	}

	// Check RoleWildCardPermission
	var roleWildCardPermissionCount int64
	db.Model(&userRegistration.RoleWildCardPermission{}).Count(&roleWildCardPermissionCount)
	if roleWildCardPermissionCount == 0 {
		if err := db.Create(&roleWildCardPermissions).Error; err != nil {
			log.Fatalf("Failed to seed RoleWildCardPermission data: %v", err)
		} else {
			log.Println("RoleWildCardPermission data seeded successfully")
		}
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
			GroupID:   1,
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
			Name:      "Carrier car fee(RISKO)",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Company Commission Fee",
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
			Name:      "Inspection Fee(JEVIC)",
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
			Name:      "Road Tax",
			Category:  "car",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Shipping fee",
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

	unitsWeight := []metaData.WeightUnit{
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

	curruncies := []metaData.Currency{
		{
			Name:      "Kenyan Shilling",
			Symbol:    "KSh",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Ugandan Shilling",
			Symbol:    "USh",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Tanzanian Shilling",
			Symbol:    "TSh",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Rwandan Franc",
			Symbol:    "RWF",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "South Sudanese Pound",
			Symbol:    "SSP",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Sudanese Pound",
			Symbol:    "SDG",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Japanese Yen",
			Symbol:    "JPY",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Pound Sterling",
			Symbol:    "GBP",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "United States Dollar",
			Symbol:    "USD",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Chinese Yuan (Renminbi)",
			Symbol:    "CNY",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the currencies table already has data
	var curCount int64
	db.Model(&metaData.Currency{}).Count(&curCount)
	if curCount == 0 {
		if err := db.Create(&curruncies).Error; err != nil {
			log.Fatalf("Failed to seed currencies: %v", err)
		} else {
			log.Println("Currency data seeded successfully")
		}
	} else {
		log.Println("Currency table already seeded, skipping...")
	}

	ports := []metaData.Port{
		{
			Name:      "Port of Tokyo",
			Location:  "Tokyo",
			Category:  "Designated Major Port",
			Function:  "International trade, Passenger ferry",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Yokohama",
			Location:  "Yokohama",
			Category:  "Designated Major Port",
			Function:  "Container shipping, Passenger ferry",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Osaka",
			Location:  "Osaka",
			Category:  "Designated Major Port",
			Function:  "Industrial activities, Commercial activities",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Kobe",
			Location:  "Kobe",
			Category:  "Designated Major Port",
			Function:  "International trade, Commercial activities",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Nagoya",
			Location:  "Nagoya",
			Category:  "Designated Major Port",
			Function:  "Cargo, Automobile shipping, Machinery",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Chiba",
			Location:  "Chiba",
			Category:  "Designated Important Port",
			Function:  "Industrial activities, Commercial activities",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Kitakyushu",
			Location:  "Kitakyushu",
			Category:  "Designated Important Port",
			Function:  "Industrial hub",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Fukuoka",
			Location:  "Fukuoka",
			Category:  "Designated Important Port",
			Function:  "Cargo, Ferries",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Hakata",
			Location:  "Fukuoka",
			Category:  "Designated Important Port",
			Function:  "Cargo, Passenger ferries",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Sapporo",
			Location:  "Sapporo",
			Category:  "Fishing Port",
			Function:  "Fishing activities",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Shimizu",
			Location:  "Shizuoka",
			Category:  "Designated Major Port",
			Function:  "Cargo, Ferry",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Sendai",
			Location:  "Sendai",
			Category:  "Designated Important Port",
			Function:  "Cargo, Passenger ferries",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Niigata",
			Location:  "Niigata",
			Category:  "Designated Important Port",
			Function:  "Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Oita",
			Location:  "Oita",
			Category:  "Designated Important Port",
			Function:  "Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Hiroshima",
			Location:  "Hiroshima",
			Category:  "Designated Important Port",
			Function:  "Cargo, Passenger ferries",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Tokuyama",
			Location:  "Yamaguchi",
			Category:  "Designated Important Port",
			Function:  "Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Akita",
			Location:  "Akita",
			Category:  "Designated Important Port",
			Function:  "Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Maizuru",
			Location:  "Kyoto",
			Category:  "Designated Important Port",
			Function:  "Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Nagasaki",
			Location:  "Nagasaki",
			Category:  "Designated Important Port",
			Function:  "Cargo, Passenger ferries",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Miyako",
			Location:  "Miyagi",
			Category:  "Designated Important Port",
			Function:  "Fishing, Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Chugoku",
			Location:  "Chugoku",
			Category:  "Designated Important Port",
			Function:  "Fishing, Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Akashi",
			Location:  "Hyogo",
			Category:  "Designated Important Port",
			Function:  "Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Kushiro",
			Location:  "Hokkaido",
			Category:  "Designated Important Port",
			Function:  "Fishing",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Wakayama",
			Location:  "Wakayama",
			Category:  "Designated Important Port",
			Function:  "Cargo",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Takamatsu",
			Location:  "Kagawa",
			Category:  "Designated Important Port",
			Function:  "Cargo, Passenger ferries",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			Name:      "Port of Otaru",
			Location:  "Hokkaido",
			Category:  "Designated Important Port",
			Function:  "Fishing",
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Check if the currencies table already has data
	var portCount int64
	db.Model(&metaData.Port{}).Count(&portCount)
	if portCount == 0 {
		if err := db.Create(&ports).Error; err != nil {
			log.Fatalf("Failed to seed ports: %v", err)
		} else {
			log.Println("Port data seeded successfully")
		}
	} else {
		log.Println("Port table already seeded, skipping...")
	}

	modes := []metaData.PaymentMode{
		{
			Mode:        "Credit/Debit Card",
			Description: "Widely accepted for online transactions.",
			Category:    "Bank",
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Mode:        "Airtel Money",
			Description: "A popular digital payment method allowing users to make purchases via mobile phones.",
			Category:    "Mobile Money",
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Mode:        "MTN Mobile Money",
			Description: "A popular digital payment method allowing users to make purchases via mobile phones.",
			Category:    "Mobile Money",
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Mode:        "Direct Cash Payments",
			Description: "Payments made directly in cash to intermediaries for international purchases.",
			Category:    "Cash",
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Mode:        "Bank Transfer",
			Description: "Payments made directly from bank accounts.",
			Category:    "Bank",
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
	}

	// Check if the currencies table already has data
	var modeCount int64
	db.Model(&metaData.PaymentMode{}).Count(&modeCount)
	if modeCount == 0 {
		if err := db.Create(&modes).Error; err != nil {
			log.Fatalf("Failed to seed payment modes: %v", err)
		} else {
			log.Println("Payment modes data seeded successfully")
		}
	} else {
		log.Println("Payment modes already seeded, skipping...")
	}

}
