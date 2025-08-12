package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/sigitisme/amf-loan-service/internal/handlers"
	"github.com/sigitisme/amf-loan-service/internal/middleware"
)

func SetupRoutes(
	r *gin.Engine,
	authService domain.AuthService,
	loanService domain.LoanService,
	investmentService domain.InvestmentService,
) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	loanHandler := handlers.NewLoanHandler(loanService)
	investmentHandler := handlers.NewInvestmentHandler(investmentService)

	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// Loan routes
		loans := api.Group("/loans")
		{
			loans.POST("", loanHandler.CreateLoan)   // Borrowers only
			loans.GET("", loanHandler.GetLoans)      // All authenticated users
			loans.GET("/my", loanHandler.GetMyLoans) // Borrowers only - specific endpoint for borrower's loans
			loans.GET("/:id", loanHandler.GetLoan)   // All authenticated users

			// Approval route - field validators only
			loans.POST("/:id/approve",
				middleware.RoleMiddleware(domain.RoleFieldValidator),
				loanHandler.ApproveLoan)

			// Disbursement route - field officers only
			loans.POST("/:id/disburse",
				middleware.RoleMiddleware(domain.RoleFieldOfficer),
				loanHandler.DisburseLoan)

			// Investment routes for loans - using same :id parameter
			loans.GET("/:id/investments", investmentHandler.GetLoanInvestments)
		}

		// Investment routes
		investments := api.Group("/investments")
		{
			investments.POST("", investmentHandler.Invest)             // Investors only
			investments.GET("/my", investmentHandler.GetMyInvestments) // Investors only
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "amf-loan-service",
		})
	})
}
