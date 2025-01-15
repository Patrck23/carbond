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

	// Hashing password for users
	passwordHash, err := hashPassword("Admin123")
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	groups := []userRegistration.Group{
		{
			Code:        "admin_group",
			Name:        "Admin Group",
			Description: "Group for administrative users",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "manager_group",
			Name:        "Manager Group",
			Description: "Group for manager users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "user_group",
			Name:        "User Group",
			Description: "Group for regular users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
	}

	// Seed Roles
	roles := []userRegistration.Role{
		{
			Code:        "admin_role",
			Name:        "Admin Role",
			Description: "Role with full permissions",
			GroupID:     1,
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "manager_role",
			Name:        "Manager Role",
			Description: "Role with permissions to manage resources and settings, but no full admin rights",
			GroupID:     2,
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "user_role",
			Name:        "User Role",
			Description: "Role with limited permissions",
			GroupID:     3,
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
	}

	// Seed Resources
	resources := []userRegistration.Resource{
		{
			Code:        "resource.*",
			Name:        "Dashboard",
			Description: "Access everything",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "resource.dashboard",
			Name:        "Dashboard",
			Description: "Access to the dashboard",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "resource.settings",
			Name:        "Settings",
			Description: "Access to settings page",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "users.admin",
			Name:        "Dashboard",
			Description: "Access to all users",
			Internal:    true,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "users.read",
			Name:        "Dashboard",
			Description: "Read users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
		{
			Code:        "users.write",
			Name:        "Dashboard",
			Description: "Edit users",
			Internal:    false,
			CreatedBy:   "Seeder",
			UpdatedBy:   "",
		},
	}

	// Seed RoleResourcePermissions
	roleResourcePermissions := []userRegistration.RoleResourcePermission{
		{
			RoleCode:     "admin_role",
			ResourceCode: "resource.*",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: true, D: true},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "manager_role",
			ResourceCode: "resource.settings",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: false, D: false}, // Allow read and write, but no execute or delete
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:     "user_role",
			ResourceCode: "resource.dashboard",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: false, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
	}

	// Seed RoleWildCardPermissions
	roleWildCardPermissions := []userRegistration.RoleWildCardPermission{
		{
			RoleCode:        "admin_role",
			ResourcePattern: "resource.*",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: true, D: true},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		}, {
			RoleCode:        "manager_role",
			ResourcePattern: "resource.*",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: false, D: false},
				Deny:  userRegistration.RWXD{R: false, W: false, X: false, D: false},
			},
			CreatedBy: "Seeder",
			UpdatedBy: "",
		},
		{
			RoleCode:        "user_role",
			ResourcePattern: "resource.settings",
			Permissions: userRegistration.Permissions{
				Allow: userRegistration.RWXD{R: true, W: true, X: false, D: false},
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

}
