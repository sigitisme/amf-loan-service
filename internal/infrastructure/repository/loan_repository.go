package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"gorm.io/gorm"
)

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) domain.LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

func (r *loanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	var loan domain.Loan
	err := r.db.WithContext(ctx).
		Preload("Borrower").
		Preload("Borrower.User").
		Preload("Approval").
		Preload("Approval.Validator").
		Preload("Investments").
		Preload("Investments.Investor").
		Preload("Investments.Investor.User").
		Preload("Disbursement").
		Preload("Disbursement.Officer").
		Where("id = ?", id).
		First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetByIDWithLock(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	var loan domain.Loan
	err := r.db.WithContext(ctx).
		Preload("Borrower").
		Preload("Borrower.User").
		Preload("Approval").
		Preload("Approval.Validator").
		Preload("Investments").
		Preload("Investments.Investor").
		Preload("Investments.Investor.User").
		Preload("Disbursement").
		Preload("Disbursement.Officer").
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", id).
		First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetByBorrowerID(ctx context.Context, borrowerID uuid.UUID) ([]domain.Loan, error) {
	var loans []domain.Loan
	err := r.db.WithContext(ctx).
		Preload("Borrower").
		Preload("Borrower.User").
		Preload("Approval").
		Preload("Investments").
		Preload("Investments.Investor").
		Preload("Disbursement").
		Where("borrower_id = ?", borrowerID).
		Find(&loans).Error
	return loans, err
}

func (r *loanRepository) GetByState(ctx context.Context, state domain.LoanState) ([]domain.Loan, error) {
	var loans []domain.Loan
	err := r.db.WithContext(ctx).
		Preload("Borrower").
		Preload("Borrower.User").
		Preload("Approval").
		Preload("Investments").
		Preload("Investments.Investor").
		Preload("Disbursement").
		Where("state = ?", state).
		Find(&loans).Error
	return loans, err
}

func (r *loanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	return r.db.WithContext(ctx).Save(loan).Error
}

func (r *loanRepository) List(ctx context.Context, limit, offset int) ([]domain.Loan, error) {
	var loans []domain.Loan
	err := r.db.WithContext(ctx).
		Preload("Borrower").
		Preload("Approval").
		Preload("Investments").
		Preload("Disbursement").
		Limit(limit).
		Offset(offset).
		Find(&loans).Error
	return loans, err
}
