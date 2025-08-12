package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type loanService struct {
	loanRepo         domain.LoanRepository
	approvalRepo     domain.ApprovalRepository
	disbursementRepo domain.DisbursementRepository
	investmentRepo   domain.InvestmentRepository
	borrowerRepo     domain.BorrowerRepository
}

func NewLoanService(
	loanRepo domain.LoanRepository,
	approvalRepo domain.ApprovalRepository,
	disbursementRepo domain.DisbursementRepository,
	investmentRepo domain.InvestmentRepository,
	borrowerRepo domain.BorrowerRepository,
) domain.LoanService {
	return &loanService{
		loanRepo:         loanRepo,
		approvalRepo:     approvalRepo,
		disbursementRepo: disbursementRepo,
		investmentRepo:   investmentRepo,
		borrowerRepo:     borrowerRepo,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, userID uuid.UUID, principalAmount, rate float64) (*domain.Loan, error) {
	// Get borrower by user ID
	borrower, err := s.borrowerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate total interest that borrower must pay
	totalInterest := principalAmount * rate

	// Calculate ROI for investors (80% of borrower's rate, platform keeps 20% margin)
	roi := rate * 0.8

	loan := &domain.Loan{
		ID:                  uuid.New(),
		BorrowerID:          borrower.ID, // Use borrower's ID, not user's ID
		PrincipalAmount:     principalAmount,
		InvestedAmount:      0,
		RemainingInvestment: principalAmount,
		Rate:                rate,
		ROI:                 roi,
		TotalInterest:       totalInterest,
		State:               domain.LoanStateProposed,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	err = s.loanRepo.Create(ctx, loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *loanService) ApproveLoan(ctx context.Context, loanID uuid.UUID, validatorID uuid.UUID, photoProofURL string, approvalDate time.Time) error {
	// Get loan
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrLoanNotFound
		}
		return err
	}

	// Check if loan is in proposed state
	if loan.State != domain.LoanStateProposed {
		return domain.ErrLoanAlreadyApproved
	}

	// Create approval record
	approval := &domain.Approval{
		ID:            uuid.New(),
		LoanID:        loanID,
		ValidatorID:   validatorID,
		PhotoProofURL: photoProofURL,
		ApprovalDate:  approvalDate,
		CreatedAt:     time.Now(),
	}

	err = s.approvalRepo.Create(ctx, approval)
	if err != nil {
		return err
	}

	// Update loan state
	loan.State = domain.LoanStateApproved
	loan.UpdatedAt = time.Now()

	return s.loanRepo.Update(ctx, loan)
}

func (s *loanService) GetLoansByState(ctx context.Context, state domain.LoanState) ([]domain.Loan, error) {
	return s.loanRepo.GetByState(ctx, state)
}

func (s *loanService) GetLoanByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	loan, err := s.loanRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLoanNotFound
		}
		return nil, err
	}
	return loan, nil
}

func (s *loanService) GetBorrowerLoans(ctx context.Context, borrowerID uuid.UUID) ([]domain.Loan, error) {
	return s.loanRepo.GetByBorrowerID(ctx, borrowerID)
}

func (s *loanService) GetBorrowerLoansByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Loan, error) {
	// Get borrower by user ID first
	borrower, err := s.borrowerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get loans by borrower ID
	return s.loanRepo.GetByBorrowerID(ctx, borrower.ID)
}

func (s *loanService) DisburseLoan(ctx context.Context, loanID uuid.UUID, officerID uuid.UUID, agreementFileURL string, disbursementDate time.Time) error {
	// Get loan
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrLoanNotFound
		}
		return err
	}

	// Check if loan is in invested state
	if loan.State != domain.LoanStateInvested {
		return domain.ErrLoanNotInvested
	}

	// Create disbursement record
	disbursement := &domain.Disbursement{
		ID:               uuid.New(),
		LoanID:           loanID,
		OfficerID:        officerID,
		AgreementFileURL: agreementFileURL,
		DisbursementDate: disbursementDate,
		CreatedAt:        time.Now(),
	}

	err = s.disbursementRepo.Create(ctx, disbursement)
	if err != nil {
		return err
	}

	// Update loan state
	loan.State = domain.LoanStateDisbursed
	loan.UpdatedAt = time.Now()

	return s.loanRepo.Update(ctx, loan)
}
