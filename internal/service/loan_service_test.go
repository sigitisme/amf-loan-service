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

// Mock repositories for loan service
type mockLoanRepository struct {
	mock.Mock
}

func (m *mockLoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *mockLoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func (m *mockLoanRepository) GetByIDWithLock(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func (m *mockLoanRepository) GetByBorrowerID(ctx context.Context, borrowerID uuid.UUID) ([]domain.Loan, error) {
	args := m.Called(ctx, borrowerID)
	return args.Get(0).([]domain.Loan), args.Error(1)
}

func (m *mockLoanRepository) GetByState(ctx context.Context, state domain.LoanState) ([]domain.Loan, error) {
	args := m.Called(ctx, state)
	return args.Get(0).([]domain.Loan), args.Error(1)
}

func (m *mockLoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *mockLoanRepository) List(ctx context.Context, limit, offset int) ([]domain.Loan, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Loan), args.Error(1)
}

type mockApprovalRepository struct {
	mock.Mock
}

func (m *mockApprovalRepository) Create(ctx context.Context, approval *domain.Approval) error {
	args := m.Called(ctx, approval)
	return args.Error(0)
}

func (m *mockApprovalRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Approval, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Approval), args.Error(1)
}

type mockDisbursementRepository struct {
	mock.Mock
}

func (m *mockDisbursementRepository) Create(ctx context.Context, disbursement *domain.Disbursement) error {
	args := m.Called(ctx, disbursement)
	return args.Error(0)
}

func (m *mockDisbursementRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Disbursement, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Disbursement), args.Error(1)
}

type mockInvestmentRepository struct {
	mock.Mock
}

func (m *mockInvestmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *mockInvestmentRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) ([]domain.Investment, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).([]domain.Investment), args.Error(1)
}

func (m *mockInvestmentRepository) GetByInvestorID(ctx context.Context, investorID uuid.UUID) ([]domain.Investment, error) {
	args := m.Called(ctx, investorID)
	return args.Get(0).([]domain.Investment), args.Error(1)
}

func (m *mockInvestmentRepository) GetTotalInvestedAmount(ctx context.Context, loanID uuid.UUID) (float64, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockInvestmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *mockInvestmentRepository) UpdateAgreementLetterURL(ctx context.Context, id uuid.UUID, url string) error {
	args := m.Called(ctx, id, url)
	return args.Error(0)
}

func (m *mockInvestmentRepository) CreateWithTx(ctx context.Context, investment *domain.Investment, loan *domain.Loan) error {
	args := m.Called(ctx, investment, loan)
	return args.Error(0)
}

func (m *mockInvestmentRepository) CreateInvestmentWithLoanLock(ctx context.Context, investment *domain.Investment, loanID uuid.UUID) (*domain.Loan, error) {
	args := m.Called(ctx, investment, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Loan), args.Error(1)
}

// Test Loan Creation - Happy Flow
func TestLoanService_CreateLoan_Success(t *testing.T) {
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
		FullName: "John Doe",
	}

	mockBorrowerRepo.On("GetByUserID", mock.Anything, userID).Return(borrower, nil)
	mockLoanRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)

	// Act
	loan, err := loanService.CreateLoan(context.Background(), userID, principalAmount, rate)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, borrowerID, loan.BorrowerID)
	assert.Equal(t, principalAmount, loan.PrincipalAmount)
	assert.Equal(t, rate, loan.Rate)
	assert.Equal(t, domain.LoanStateProposed, loan.State)
	assert.Equal(t, rate*0.8, loan.ROI) // 80% of borrower rate
	assert.Equal(t, principalAmount*rate, loan.TotalInterest)

	mockBorrowerRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
}

// Test Loan Approval - Happy Flow
func TestLoanService_ApproveLoan_Success(t *testing.T) {
	// Arrange
	mockLoanRepo := new(mockLoanRepository)
	mockApprovalRepo := new(mockApprovalRepository)
	mockDisbursementRepo := new(mockDisbursementRepository)
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockBorrowerRepo := new(mockBorrowerRepository)

	loanService := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockBorrowerRepo)

	loanID := uuid.New()
	validatorID := uuid.New()
	photoProofURL := "https://example.com/proof.jpg"
	approvalDate := time.Now()

	existingLoan := &domain.Loan{
		ID:    loanID,
		State: domain.LoanStateProposed,
	}

	mockLoanRepo.On("GetByID", mock.Anything, loanID).Return(existingLoan, nil)
	mockApprovalRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Approval")).Return(nil)
	mockLoanRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)

	// Act
	err := loanService.ApproveLoan(context.Background(), loanID, validatorID, photoProofURL, approvalDate)

	// Assert
	assert.NoError(t, err)

	mockLoanRepo.AssertExpectations(t)
	mockApprovalRepo.AssertExpectations(t)
}

// Test Get Loans by State - Happy Flow
func TestLoanService_GetLoansByState_Success(t *testing.T) {
	// Arrange
	mockLoanRepo := new(mockLoanRepository)
	mockApprovalRepo := new(mockApprovalRepository)
	mockDisbursementRepo := new(mockDisbursementRepository)
	mockInvestmentRepo := new(mockInvestmentRepository)
	mockBorrowerRepo := new(mockBorrowerRepository)

	loanService := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockBorrowerRepo)

	expectedLoans := []domain.Loan{
		{ID: uuid.New(), State: domain.LoanStateProposed},
		{ID: uuid.New(), State: domain.LoanStateProposed},
	}

	mockLoanRepo.On("GetByState", mock.Anything, domain.LoanStateProposed).Return(expectedLoans, nil)

	// Act
	loans, err := loanService.GetLoansByState(context.Background(), domain.LoanStateProposed)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, loans, 2)
	assert.Equal(t, expectedLoans, loans)

	mockLoanRepo.AssertExpectations(t)
}
