package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test Notification Service - Happy Flow
func TestNotificationService_SendAgreementLetters_Success(t *testing.T) {
	// Arrange
	mockLoanRepo := new(mockLoanRepository)
	mockInvestmentRepo := new(mockInvestmentRepository)

	notificationService := NewNotificationService(mockLoanRepo, mockInvestmentRepo)

	loanID := uuid.New()

	// Mock investments with investor and user data
	investments := []domain.Investment{
		{
			ID:         uuid.New(),
			LoanID:     loanID,
			InvestorID: uuid.New(),
			Amount:     25000,
			Investor: domain.Investor{
				FullName: "John Investor",
				User: domain.User{
					Email: "john@example.com",
				},
			},
		},
		{
			ID:         uuid.New(),
			LoanID:     loanID,
			InvestorID: uuid.New(),
			Amount:     30000,
			Investor: domain.Investor{
				FullName: "Jane Investor",
				User: domain.User{
					Email: "jane@example.com",
				},
			},
		},
	}

	mockInvestmentRepo.On("GetByLoanID", mock.Anything, loanID).Return(investments, nil)
	// Mock UpdateAgreementLetterURL for each investment
	for _, investment := range investments {
		mockInvestmentRepo.On("UpdateAgreementLetterURL", mock.Anything, investment.ID, mock.AnythingOfType("string")).Return(nil)
	}

	// Act
	err := notificationService.SendAgreementLetters(context.Background(), loanID)

	// Assert
	assert.NoError(t, err)

	mockInvestmentRepo.AssertExpectations(t)
}

// Test Notification Service - Generate Agreement Letter URL
func TestNotificationService_GenerateAgreementLetterURL(t *testing.T) {
	// Arrange
	mockLoanRepo := new(mockLoanRepository)
	mockInvestmentRepo := new(mockInvestmentRepository)

	notificationService := NewNotificationService(mockLoanRepo, mockInvestmentRepo).(*notificationService)

	loanID := uuid.New()
	investorID := uuid.New()
	investmentID := uuid.New()

	// Act
	url := notificationService.generateAgreementLetterURL(loanID, investorID, investmentID)

	// Assert
	expectedURL := "https://amf-documents.s3.amazonaws.com/agreements/loan_" +
		loanID.String() + "/investor_" + investorID.String() + "/agreement_" + investmentID.String() + ".pdf"
	assert.Equal(t, expectedURL, url)
	assert.Contains(t, url, "https://amf-documents.s3.amazonaws.com/agreements")
	assert.Contains(t, url, loanID.String())
	assert.Contains(t, url, investorID.String())
	assert.Contains(t, url, investmentID.String())
	assert.Contains(t, url, ".pdf")
}

// Test Notification Service - No Investments
func TestNotificationService_SendAgreementLetters_NoInvestments(t *testing.T) {
	// Arrange
	mockLoanRepo := new(mockLoanRepository)
	mockInvestmentRepo := new(mockInvestmentRepository)

	notificationService := NewNotificationService(mockLoanRepo, mockInvestmentRepo)

	loanID := uuid.New()

	// Return empty investments array
	emptyInvestments := []domain.Investment{}

	mockInvestmentRepo.On("GetByLoanID", mock.Anything, loanID).Return(emptyInvestments, nil)

	// Act
	err := notificationService.SendAgreementLetters(context.Background(), loanID)

	// Assert
	assert.NoError(t, err) // Should not error even with no investments

	mockInvestmentRepo.AssertExpectations(t)
}
