package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/config"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/sigitisme/amf-loan-service/internal/infrastructure/database"
	"github.com/sigitisme/amf-loan-service/internal/infrastructure/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	err = database.Migrate(db)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	borrowerRepo := repository.NewBorrowerRepository(db)
	investorRepo := repository.NewInvestorRepository(db)

	ctx := context.Background()

	log.Println("Creating mock users, borrowers, and investors...")

	// Create mock borrowers
	borrowers := []struct {
		email          string
		password       string
		fullName       string
		phoneNumber    string
		address        string
		identityNumber string
	}{
		{
			email:          "borrower1@example.com",
			password:       "password123",
			fullName:       "John Doe",
			phoneNumber:    "+1234567890",
			address:        "123 Main St, New York, NY",
			identityNumber: "B001234567",
		},
		{
			email:          "borrower2@example.com",
			password:       "password123",
			fullName:       "Alice Johnson",
			phoneNumber:    "+1234567891",
			address:        "456 Oak Ave, Los Angeles, CA",
			identityNumber: "B001234568",
		},
		{
			email:          "borrower3@example.com",
			password:       "password123",
			fullName:       "Bob Smith",
			phoneNumber:    "+1234567892",
			address:        "789 Pine St, Chicago, IL",
			identityNumber: "B001234569",
		},
	}

	for _, b := range borrowers {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(b.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", b.email, err)
			continue
		}

		// Create user
		user := &domain.User{
			ID:        uuid.New(),
			Email:     b.email,
			Password:  string(hashedPassword),
			Role:      domain.RoleBorrower,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = userRepo.Create(ctx, user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", b.email, err)
			continue
		}

		// Create borrower profile
		borrower := &domain.Borrower{
			ID:             uuid.New(),
			UserID:         user.ID,
			FullName:       b.fullName,
			PhoneNumber:    b.phoneNumber,
			Address:        b.address,
			IdentityNumber: b.identityNumber,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		err = borrowerRepo.Create(ctx, borrower)
		if err != nil {
			log.Printf("Failed to create borrower profile for %s: %v", b.email, err)
			continue
		}

		log.Printf("✅ Created borrower: %s (%s)", b.fullName, b.email)
	}

	// Create mock investors
	investors := []struct {
		email          string
		password       string
		fullName       string
		phoneNumber    string
		address        string
		identityNumber string
	}{
		{
			email:          "investor1@example.com",
			password:       "password123",
			fullName:       "Emma Wilson",
			phoneNumber:    "+1987654320",
			address:        "321 Elm St, Boston, MA",
			identityNumber: "I001234567",
		},
		{
			email:          "investor2@example.com",
			password:       "password123",
			fullName:       "Michael Brown",
			phoneNumber:    "+1987654321",
			address:        "654 Maple Ave, Seattle, WA",
			identityNumber: "I001234568",
		},
		{
			email:          "investor3@example.com",
			password:       "password123",
			fullName:       "Sarah Davis",
			phoneNumber:    "+1987654322",
			address:        "987 Cedar Blvd, Miami, FL",
			identityNumber: "I001234569",
		},
		{
			email:          "investor4@example.com",
			password:       "password123",
			fullName:       "David Lee",
			phoneNumber:    "+1987654323",
			address:        "147 Birch Lane, Denver, CO",
			identityNumber: "I001234570",
		},
		{
			email:          "investor5@example.com",
			password:       "password123",
			fullName:       "Lisa Martinez",
			phoneNumber:    "+1987654324",
			address:        "258 Willow Dr, Austin, TX",
			identityNumber: "I001234571",
		},
	}

	for _, i := range investors {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(i.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", i.email, err)
			continue
		}

		// Create user
		user := &domain.User{
			ID:        uuid.New(),
			Email:     i.email,
			Password:  string(hashedPassword),
			Role:      domain.RoleInvestor,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = userRepo.Create(ctx, user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", i.email, err)
			continue
		}

		// Create investor profile
		investor := &domain.Investor{
			ID:             uuid.New(),
			UserID:         user.ID,
			FullName:       i.fullName,
			PhoneNumber:    i.phoneNumber,
			Address:        i.address,
			IdentityNumber: i.identityNumber,
			TotalInvested:  0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		err = investorRepo.Create(ctx, investor)
		if err != nil {
			log.Printf("Failed to create investor profile for %s: %v", i.email, err)
			continue
		}

		log.Printf("✅ Created investor: %s (%s)", i.fullName, i.email)
	}

	// Create field validator and field officer
	staffUsers := []struct {
		email    string
		password string
		role     domain.UserRole
		name     string
	}{
		{
			email:    "validator@amf.com",
			password: "validator123",
			role:     domain.RoleFieldValidator,
			name:     "Field Validator",
		},
		{
			email:    "officer@amf.com",
			password: "officer123",
			role:     domain.RoleFieldOfficer,
			name:     "Field Officer",
		},
	}

	for _, s := range staffUsers {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", s.email, err)
			continue
		}

		// Create user
		user := &domain.User{
			ID:        uuid.New(),
			Email:     s.email,
			Password:  string(hashedPassword),
			Role:      s.role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = userRepo.Create(ctx, user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", s.email, err)
			continue
		}

		log.Printf("✅ Created %s: %s", s.name, s.email)
	}

	log.Println("🎉 Mock data creation completed!")
	log.Println("")
	log.Println("📋 Created Accounts:")
	log.Println("👤 Borrowers:")
	log.Println("   - borrower1@example.com (John Doe)")
	log.Println("   - borrower2@example.com (Alice Johnson)")
	log.Println("   - borrower3@example.com (Bob Smith)")
	log.Println("")
	log.Println("💰 Investors:")
	log.Println("   - investor1@example.com (Emma Wilson)")
	log.Println("   - investor2@example.com (Michael Brown)")
	log.Println("   - investor3@example.com (Sarah Davis)")
	log.Println("   - investor4@example.com (David Lee)")
	log.Println("   - investor5@example.com (Lisa Martinez)")
	log.Println("")
	log.Println("🏢 Staff:")
	log.Println("   - validator@amf.com (Field Validator)")
	log.Println("   - officer@amf.com (Field Officer)")
	log.Println("")
	log.Println("🔑 All passwords: password123 (except staff: validator123/officer123)")
}
