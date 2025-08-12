package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test User Entity Creation
func TestUser_Creation(t *testing.T) {
	// Arrange & Act
	user := User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      RoleInvestor,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Assert
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, RoleInvestor, user.Role)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

// Test User Roles
func TestUserRoles(t *testing.T) {
	// Test all user roles exist
	assert.Equal(t, UserRole("borrower"), RoleBorrower)
	assert.Equal(t, UserRole("investor"), RoleInvestor)
	assert.Equal(t, UserRole("field_officer"), RoleFieldOfficer)
	assert.Equal(t, UserRole("field_validator"), RoleFieldValidator)
}

// Test Loan States
func TestLoanStates(t *testing.T) {
	// Test all loan states exist
	assert.Equal(t, LoanState("proposed"), LoanStateProposed)
	assert.Equal(t, LoanState("approved"), LoanStateApproved)
	assert.Equal(t, LoanState("invested"), LoanStateInvested)
	assert.Equal(t, LoanState("disbursed"), LoanStateDisbursed)
}

// Test Loan Entity Creation
func TestLoan_Creation(t *testing.T) {
	// Arrange & Act
	borrowerID := uuid.New()
	loan := Loan{
		ID:                  uuid.New(),
		BorrowerID:          borrowerID,
		PrincipalAmount:     100000.0,
		InvestedAmount:      0,
		RemainingInvestment: 100000.0,
		Rate:                0.12,
		ROI:                 0.096, // 80% of rate
		TotalInterest:       12000.0,
		State:               LoanStateProposed,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Assert
	assert.NotEmpty(t, loan.ID)
	assert.Equal(t, borrowerID, loan.BorrowerID)
	assert.Equal(t, 100000.0, loan.PrincipalAmount)
	assert.Equal(t, 0.0, loan.InvestedAmount)
	assert.Equal(t, 100000.0, loan.RemainingInvestment)
	assert.Equal(t, 0.12, loan.Rate)
	assert.Equal(t, 0.096, loan.ROI)
	assert.Equal(t, 12000.0, loan.TotalInterest)
	assert.Equal(t, LoanStateProposed, loan.State)
}

// Test Borrower Entity Creation
func TestBorrower_Creation(t *testing.T) {
	// Arrange & Act
	userID := uuid.New()
	borrower := Borrower{
		ID:             uuid.New(),
		UserID:         userID,
		FullName:       "John Doe",
		PhoneNumber:    "+1234567890",
		Address:        "123 Main St",
		IdentityNumber: "ID123456789",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Assert
	assert.NotEmpty(t, borrower.ID)
	assert.Equal(t, userID, borrower.UserID)
	assert.Equal(t, "John Doe", borrower.FullName)
	assert.Equal(t, "+1234567890", borrower.PhoneNumber)
	assert.Equal(t, "123 Main St", borrower.Address)
	assert.Equal(t, "ID123456789", borrower.IdentityNumber)
}

// Test Investor Entity Creation
func TestInvestor_Creation(t *testing.T) {
	// Arrange & Act
	userID := uuid.New()
	investor := Investor{
		ID:             uuid.New(),
		UserID:         userID,
		FullName:       "Jane Smith",
		PhoneNumber:    "+0987654321",
		Address:        "456 Oak St",
		IdentityNumber: "ID987654321",
		TotalInvested:  50000.0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Assert
	assert.NotEmpty(t, investor.ID)
	assert.Equal(t, userID, investor.UserID)
	assert.Equal(t, "Jane Smith", investor.FullName)
	assert.Equal(t, "+0987654321", investor.PhoneNumber)
	assert.Equal(t, "456 Oak St", investor.Address)
	assert.Equal(t, "ID987654321", investor.IdentityNumber)
	assert.Equal(t, 50000.0, investor.TotalInvested)
}

// Test Investment Entity Creation
func TestInvestment_Creation(t *testing.T) {
	// Arrange & Act
	loanID := uuid.New()
	investorID := uuid.New()
	investment := Investment{
		ID:                 uuid.New(),
		LoanID:             loanID,
		InvestorID:         investorID,
		Amount:             25000.0,
		Status:             "completed",
		AgreementLetterURL: "https://example.com/agreement.pdf",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Assert
	assert.NotEmpty(t, investment.ID)
	assert.Equal(t, loanID, investment.LoanID)
	assert.Equal(t, investorID, investment.InvestorID)
	assert.Equal(t, 25000.0, investment.Amount)
	assert.Equal(t, "completed", investment.Status)
	assert.Equal(t, "https://example.com/agreement.pdf", investment.AgreementLetterURL)
}

// Test Investment Event Creation
func TestInvestmentEvent_Creation(t *testing.T) {
	// Arrange & Act
	loanID := uuid.New()
	investorID := uuid.New()
	timestamp := time.Now()

	event := InvestmentEvent{
		ID:         uuid.New(),
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     30000.0,
		Timestamp:  timestamp,
	}

	// Assert
	assert.NotEmpty(t, event.ID)
	assert.Equal(t, loanID, event.LoanID)
	assert.Equal(t, investorID, event.InvestorID)
	assert.Equal(t, 30000.0, event.Amount)
	assert.Equal(t, timestamp, event.Timestamp)
}

// Test Login Response Creation
func TestLoginResponse_Creation(t *testing.T) {
	// Arrange & Act
	userID := uuid.New()
	expiresAt := time.Now().Add(time.Hour)

	response := LoginResponse{
		UserID:    userID,
		Email:     "user@example.com",
		Token:     "jwt-token-here",
		ExpiresAt: expiresAt,
	}

	// Assert
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "user@example.com", response.Email)
	assert.Equal(t, "jwt-token-here", response.Token)
	assert.Equal(t, expiresAt, response.ExpiresAt)
}
