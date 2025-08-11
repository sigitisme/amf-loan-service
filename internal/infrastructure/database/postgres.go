package database

import (
	"fmt"
	"log"

	"github.com/sigitisme/amf-loan-service/internal/config"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	return db.AutoMigrate(
		&domain.User{},
		&domain.Borrower{},
		&domain.Investor{},
		&domain.Loan{},
		&domain.Approval{},
		&domain.Investment{},
		&domain.Disbursement{},
	)
}
