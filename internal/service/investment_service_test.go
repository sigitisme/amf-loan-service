package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Kafka Producer
type mockKafkaProducer struct {
	mock.Mock
}

func (m *mockKafkaProducer) PublishInvestmentEvent(ctx context.Context, event domain.InvestmentEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockKafkaProducer) PublishFullyFundedLoan(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

// Mock Notification Service
type mockNotificationService struct {
	mock.Mock
}

func (m *mockNotificationService) SendAgreementLetters(ctx context.Context, loanID uuid.UUID) error {
	args := m.Called(ctx, loanID)
	return args.Error(0)
}

// Test Investment Request - Happy Flow
func TestInvestmentService_RequestInvestment_Success(t *testing.T) {
	// Arrange
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockLoanRepo := new(mockLoanRepository)
	mockInvestorRepo := new(mockInvestorRepository)
	mockKafkaProducer := new(mockKafkaProducer)
	mockNotificationService := new(mockNotificationService)

	investmentService := NewInvestmentService(mockInvestmentRepo, mockLoanRepo, mockInvestorRepo, mockKafkaProducer, mockNotificationService)

	userID := uuid.New()
	loanID := uuid.New()
	investorID := uuid.New()
	borrowerUserID := uuid.New() // Different from investor
	amount := 50000.0

	investor := &domain.Investor{
		ID:     investorID,
		UserID: userID,
	}

	loan := &domain.Loan{
		ID:                  loanID,
		State:               domain.LoanStateApproved,
		RemainingInvestment: 100000.0,
		Borrower: domain.Borrower{
			UserID: borrowerUserID, // Different user ID to prevent self-investment
		},
	}

	mockInvestorRepo.On("GetByUserID", mock.Anything, userID).Return(investor, nil)
	mockLoanRepo.On("GetByID", mock.Anything, loanID).Return(loan, nil)
	mockKafkaProducer.On("PublishInvestmentEvent", mock.Anything, mock.AnythingOfType("domain.InvestmentEvent")).Return(nil)

	// Act
	err := investmentService.RequestInvestment(context.Background(), userID, loanID, amount)

	// Assert
	assert.NoError(t, err)

	mockInvestorRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
	mockKafkaProducer.AssertExpectations(t)
}

// Test Investment Processing - Happy Flow
func TestInvestmentService_ProcessInvestment_Success(t *testing.T) {
	// Arrange
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockLoanRepo := new(mockLoanRepository)
	mockInvestorRepo := new(mockInvestorRepository)
	mockKafkaProducer := new(mockKafkaProducer)
	mockNotificationService := new(mockNotificationService)

	investmentService := NewInvestmentService(mockInvestmentRepo, mockLoanRepo, mockInvestorRepo, mockKafkaProducer, mockNotificationService)

	eventID := uuid.New()
	loanID := uuid.New()
	investorID := uuid.New()
	amount := 50000.0

	event := domain.InvestmentEvent{
		ID:         eventID,
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     amount,
		Timestamp:  time.Now(),
	}

	loan := &domain.Loan{
		ID:                  loanID,
		State:               domain.LoanStateApproved,
		InvestedAmount:      0,
		RemainingInvestment: 100000.0,
	}

	mockLoanRepo.On("GetByIDWithLock", mock.Anything, loanID).Return(loan, nil)
	mockInvestmentRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*domain.Investment"), mock.AnythingOfType("*domain.Loan")).Return(nil)

	// Act
	err := investmentService.ProcessInvestment(context.Background(), event)

	// Assert
	assert.NoError(t, err)

	mockLoanRepo.AssertExpectations(t)
	mockInvestmentRepo.AssertExpectations(t)
}

// Test Investment Processing - Fully Funded
func TestInvestmentService_ProcessInvestment_FullyFunded(t *testing.T) {
	// Arrange
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockLoanRepo := new(mockLoanRepository)
	mockInvestorRepo := new(mockInvestorRepository)
	mockKafkaProducer := new(mockKafkaProducer)
	mockNotificationService := new(mockNotificationService)

	investmentService := NewInvestmentService(mockInvestmentRepo, mockLoanRepo, mockInvestorRepo, mockKafkaProducer, mockNotificationService)

	eventID := uuid.New()
	loanID := uuid.New()
	investorID := uuid.New()
	amount := 100000.0 // This will fully fund the loan

	event := domain.InvestmentEvent{
		ID:         eventID,
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     amount,
		Timestamp:  time.Now(),
	}

	loan := &domain.Loan{
		ID:                  loanID,
		State:               domain.LoanStateApproved,
		InvestedAmount:      0,
		RemainingInvestment: 100000.0, // Exactly the investment amount
	}

	mockLoanRepo.On("GetByIDWithLock", mock.Anything, loanID).Return(loan, nil)
	mockInvestmentRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*domain.Investment"), mock.AnythingOfType("*domain.Loan")).Return(nil)
	mockKafkaProducer.On("PublishFullyFundedLoan", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)
	mockNotificationService.On("SendAgreementLetters", mock.Anything, loanID).Return(nil)

	// Act
	err := investmentService.ProcessInvestment(context.Background(), event)

	// Assert
	assert.NoError(t, err)

	mockLoanRepo.AssertExpectations(t)
	mockInvestmentRepo.AssertExpectations(t)
	mockKafkaProducer.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

// Test Get Investor Investments - Happy Flow
func TestInvestmentService_GetInvestorInvestments_Success(t *testing.T) {
	// Arrange
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockLoanRepo := new(mockLoanRepository)
	mockInvestorRepo := new(mockInvestorRepository)
	mockKafkaProducer := new(mockKafkaProducer)
	mockNotificationService := new(mockNotificationService)

	investmentService := NewInvestmentService(mockInvestmentRepo, mockLoanRepo, mockInvestorRepo, mockKafkaProducer, mockNotificationService)

	investorID := uuid.New()

	expectedInvestments := []domain.Investment{
		{ID: uuid.New(), InvestorID: investorID, Amount: 25000},
		{ID: uuid.New(), InvestorID: investorID, Amount: 30000},
	}

	mockInvestmentRepo.On("GetByInvestorID", mock.Anything, investorID).Return(expectedInvestments, nil)

	// Act
	investments, err := investmentService.GetInvestorInvestments(context.Background(), investorID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, investments, 2)
	assert.Equal(t, expectedInvestments, investments)

	mockInvestmentRepo.AssertExpectations(t)
}

// Test Self Investment Prevention
func TestInvestmentService_RequestInvestment_SelfInvestmentError(t *testing.T) {
	// Arrange
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockLoanRepo := new(mockLoanRepository)
	mockInvestorRepo := new(mockInvestorRepository)
	mockKafkaProducer := new(mockKafkaProducer)
	mockNotificationService := new(mockNotificationService)

	investmentService := NewInvestmentService(mockInvestmentRepo, mockLoanRepo, mockInvestorRepo, mockKafkaProducer, mockNotificationService)

	userID := uuid.New() // Same user ID for both investor and borrower
	loanID := uuid.New()
	investorID := uuid.New()
	amount := 50000.0

	investor := &domain.Investor{
		ID:     investorID,
		UserID: userID,
	}

	loan := &domain.Loan{
		ID:                  loanID,
		State:               domain.LoanStateApproved,
		RemainingInvestment: 100000.0,
		Borrower: domain.Borrower{
			UserID: userID, // Same user ID as investor (self-investment)
		},
	}

	mockInvestorRepo.On("GetByUserID", mock.Anything, userID).Return(investor, nil)
	mockLoanRepo.On("GetByID", mock.Anything, loanID).Return(loan, nil)

	// Act
	err := investmentService.RequestInvestment(context.Background(), userID, loanID, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrSelfInvestment, err)

	mockInvestorRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
}
