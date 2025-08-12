package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Simple integration-style tests that test the business logic

// Test Loan Service Business Logic
func TestLoanService_BusinessLogic_CreateLoan(t *testing.T) {
	// Arrange
	mockLoanRepo := new(mockLoanRepository)
	mockApprovalRepo := new(mockApprovalRepository)
	mockDisbursementRepo := new(mockDisbursementRepository)
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockBorrowerRepo := new(mockBorrowerRepository)

	loanService := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockBorrowerRepo)

	userID := uuid.New()
	borrowerID := uuid.New()
	principalAmount := 100000.0
	rate := 0.12

	borrower := &domain.Borrower{
		ID:       borrowerID,
		UserID:   userID,
		FullName: "Test Borrower",
	}

	mockBorrowerRepo.On("GetByUserID", mock.Anything, userID).Return(borrower, nil)
	mockLoanRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)

	// Act
	loan, err := loanService.CreateLoan(context.Background(), userID, principalAmount, rate)

	// Assert - Test Business Logic
	assert.NoError(t, err)
	assert.NotNil(t, loan)

	// Verify business calculations
	expectedROI := rate * 0.8 // 80% of borrower rate
	expectedTotalInterest := principalAmount * rate

	assert.Equal(t, principalAmount, loan.PrincipalAmount)
	assert.Equal(t, principalAmount, loan.RemainingInvestment) // Initially all remaining
	assert.Equal(t, 0.0, loan.InvestedAmount)                  // Initially no investment
	assert.Equal(t, rate, loan.Rate)
	assert.Equal(t, expectedROI, loan.ROI)
	assert.Equal(t, expectedTotalInterest, loan.TotalInterest)
	assert.Equal(t, domain.LoanStateProposed, loan.State)

	mockBorrowerRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
}

// Test Investment Service Business Logic
func TestInvestmentService_BusinessLogic_ProcessInvestment(t *testing.T) {
	// Arrange
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockLoanRepo := new(mockLoanRepository)
	mockInvestorRepo := new(mockInvestorRepository)
	mockKafkaProducer := new(mockKafkaProducer)
	mockNotificationService := new(mockNotificationService)

	investmentService := NewInvestmentService(mockInvestmentRepo, mockLoanRepo, mockInvestorRepo, mockKafkaProducer, mockNotificationService)

	loanID := uuid.New()
	investorID := uuid.New()
	investmentAmount := 30000.0

	event := domain.InvestmentEvent{
		ID:         uuid.New(),
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     investmentAmount,
		Timestamp:  time.Now(),
	}

	// Loan with partial funding opportunity
	loan := &domain.Loan{
		ID:                  loanID,
		State:               domain.LoanStateApproved,
		InvestedAmount:      20000.0, // Already partially funded
		RemainingInvestment: 80000.0, // 100k total - 20k invested
		PrincipalAmount:     100000.0,
	}

	mockLoanRepo.On("GetByIDWithLock", mock.Anything, loanID).Return(loan, nil)

	// Capture the loan state changes
	var capturedInvestment *domain.Investment
	var capturedLoan *domain.Loan
	mockInvestmentRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*domain.Investment"), mock.AnythingOfType("*domain.Loan")).
		Run(func(args mock.Arguments) {
			capturedInvestment = args.Get(1).(*domain.Investment)
			capturedLoan = args.Get(2).(*domain.Loan)
		}).Return(nil)

	// Act
	err := investmentService.ProcessInvestment(context.Background(), event)

	// Assert - Test Business Logic
	assert.NoError(t, err)

	// Verify investment record
	assert.Equal(t, event.ID, capturedInvestment.ID)
	assert.Equal(t, event.LoanID, capturedInvestment.LoanID)
	assert.Equal(t, event.InvestorID, capturedInvestment.InvestorID)
	assert.Equal(t, event.Amount, capturedInvestment.Amount)
	assert.Equal(t, "completed", capturedInvestment.Status)

	// Verify loan calculations
	expectedInvestedAmount := 20000.0 + investmentAmount // Previous + new
	expectedRemainingInvestment := 100000.0 - expectedInvestedAmount

	assert.Equal(t, expectedInvestedAmount, capturedLoan.InvestedAmount)
	assert.Equal(t, expectedRemainingInvestment, capturedLoan.RemainingInvestment)
	assert.Equal(t, domain.LoanStateApproved, capturedLoan.State) // Still approved, not fully funded

	mockLoanRepo.AssertExpectations(t)
	mockInvestmentRepo.AssertExpectations(t)
}

// Test Investment Service - Full Funding Scenario
func TestInvestmentService_BusinessLogic_FullyFunded(t *testing.T) {
	// Arrange
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockLoanRepo := new(mockLoanRepository)
	mockInvestorRepo := new(mockInvestorRepository)
	mockKafkaProducer := new(mockKafkaProducer)
	mockNotificationService := new(mockNotificationService)

	investmentService := NewInvestmentService(mockInvestmentRepo, mockLoanRepo, mockInvestorRepo, mockKafkaProducer, mockNotificationService)

	loanID := uuid.New()
	investorID := uuid.New()
	investmentAmount := 50000.0 // This will complete the funding

	event := domain.InvestmentEvent{
		ID:         uuid.New(),
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     investmentAmount,
		Timestamp:  time.Now(),
	}

	// Loan that will be fully funded after this investment
	loan := &domain.Loan{
		ID:                  loanID,
		State:               domain.LoanStateApproved,
		InvestedAmount:      50000.0, // Already half funded
		RemainingInvestment: 50000.0, // Exactly the investment amount
		PrincipalAmount:     100000.0,
	}

	mockLoanRepo.On("GetByIDWithLock", mock.Anything, loanID).Return(loan, nil)

	// Capture the loan state changes
	var capturedLoan *domain.Loan
	mockInvestmentRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*domain.Investment"), mock.AnythingOfType("*domain.Loan")).
		Run(func(args mock.Arguments) {
			capturedLoan = args.Get(2).(*domain.Loan)
		}).Return(nil)

	// Mock the fully funded flow
	mockKafkaProducer.On("PublishFullyFundedLoan", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)
	mockNotificationService.On("SendAgreementLetters", mock.Anything, loanID).Return(nil)

	// Act
	err := investmentService.ProcessInvestment(context.Background(), event)

	// Assert - Test Business Logic for Full Funding
	assert.NoError(t, err)

	// Verify loan is fully funded
	assert.Equal(t, 100000.0, capturedLoan.InvestedAmount)        // Total principal
	assert.Equal(t, 0.0, capturedLoan.RemainingInvestment)        // No remaining investment
	assert.Equal(t, domain.LoanStateInvested, capturedLoan.State) // Changed to invested state

	// Verify fully funded events were triggered
	mockKafkaProducer.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
	mockInvestmentRepo.AssertExpectations(t)
}

// Test Notification Service Business Logic
func TestNotificationService_BusinessLogic_GenerateURLs(t *testing.T) {
	// Arrange
	mockLoanRepo := new(mockLoanRepository)
	mockInvestmentRepo := new(mockInvestmentRepository)

	notificationService := NewNotificationService(mockLoanRepo, mockInvestmentRepo).(*notificationService)

	loanID := uuid.New()
	investorID := uuid.New()
	investmentID := uuid.New()

	// Act - Test URL Generation Business Logic
	url := notificationService.generateAgreementLetterURL(loanID, investorID, investmentID)

	// Assert - Test URL Format
	expectedPrefix := "https://amf-documents.s3.amazonaws.com/agreements"
	expectedSuffix := ".pdf"

	assert.Contains(t, url, expectedPrefix)
	assert.Contains(t, url, loanID.String())
	assert.Contains(t, url, investorID.String())
	assert.Contains(t, url, investmentID.String())
	assert.True(t, strings.HasSuffix(url, expectedSuffix))

	// Verify URL structure follows expected pattern
	expectedURL := expectedPrefix + "/loan_" + loanID.String() +
		"/investor_" + investorID.String() +
		"/agreement_" + investmentID.String() + expectedSuffix
	assert.Equal(t, expectedURL, url)
}
