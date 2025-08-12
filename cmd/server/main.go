package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sigitisme/amf-loan-service/internal/config"
	"github.com/sigitisme/amf-loan-service/internal/infrastructure/database"
	"github.com/sigitisme/amf-loan-service/internal/infrastructure/email"
	"github.com/sigitisme/amf-loan-service/internal/infrastructure/kafka"
	"github.com/sigitisme/amf-loan-service/internal/infrastructure/repository"
	"github.com/sigitisme/amf-loan-service/internal/routes"
	"github.com/sigitisme/amf-loan-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	borrowerRepo := repository.NewBorrowerRepository(db)
	investorRepo := repository.NewInvestorRepository(db)
	loanRepo := repository.NewLoanRepository(db)
	approvalRepo := repository.NewApprovalRepository(db)
	investmentRepo := repository.NewInvestmentRepository(db)
	disbursementRepo := repository.NewDisbursementRepository(db)

	// Initialize infrastructure services
	kafkaProducer := kafka.NewProducer(&cfg.Kafka)
	defer kafkaProducer.Close()

	emailService := email.NewService(&cfg.SMTP)

	// Initialize business services
	authService := service.NewAuthService(userRepo, borrowerRepo, investorRepo, &cfg.JWT)
	loanService := service.NewLoanService(loanRepo, approvalRepo, disbursementRepo, investmentRepo, borrowerRepo)
	investmentService := service.NewInvestmentService(investmentRepo, loanRepo, investorRepo, kafkaProducer)
	_ = service.NewNotificationService(loanRepo, investmentRepo, emailService) // Available for future use

	// Initialize and start Kafka consumer
	consumer := kafka.NewConsumer(&cfg.Kafka, investmentService)
	go func() {
		if err := consumer.StartConsumer(context.Background()); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()
	defer consumer.StopConsumer()

	// Setup Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(r, authService, loanService, investmentService)

	// Start server
	log.Printf("Server starting on port %s", cfg.API.Port)
	if err := r.Run(":" + cfg.API.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
