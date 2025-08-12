package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"gorm.io/gorm"
)

type investmentRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) domain.InvestmentRepository {
	return &investmentRepository{db: db}
}

func (r *investmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	return r.db.WithContext(ctx).Create(investment).Error
}

func (r *investmentRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) ([]domain.Investment, error) {
	var investments []domain.Investment
	err := r.db.WithContext(ctx).
		Preload("Investor").
		Preload("Investor.User").
		Where("loan_id = ?", loanID).
		Find(&investments).Error
	return investments, err
}

func (r *investmentRepository) GetByInvestorID(ctx context.Context, investorID uuid.UUID) ([]domain.Investment, error) {
	var investments []domain.Investment
	err := r.db.WithContext(ctx).
		Preload("Loan").
		Preload("Loan.Borrower").
		Preload("Loan.Borrower.User").
		Where("investor_id = ?", investorID).
		Find(&investments).Error
	return investments, err
}

func (r *investmentRepository) GetTotalInvestedAmount(ctx context.Context, loanID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&domain.Investment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("loan_id = ? AND status = ?", loanID, "completed").
		Scan(&total).Error
	return total, err
}

func (r *investmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Investment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *investmentRepository) UpdateAgreementLetterURL(ctx context.Context, id uuid.UUID, url string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Investment{}).
		Where("id = ?", id).
		Update("agreement_letter_url", url).Error
}

func (r *investmentRepository) CreateWithTx(ctx context.Context, investment *domain.Investment, loan *domain.Loan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the investment
		if err := tx.Create(investment).Error; err != nil {
			return err
		}

		// Update loan amounts and state
		if err := tx.Save(loan).Error; err != nil {
			return err
		}

		// Update investor total invested
		if err := tx.Model(&domain.Investor{}).
			Where("id = ?", investment.InvestorID).
			Update("total_invested", gorm.Expr("total_invested + ?", investment.Amount)).Error; err != nil {
			return err
		}

		return nil
	})
}

// CreateInvestmentWithLoanLock atomically locks the loan and creates investment in the same transaction
func (r *investmentRepository) CreateInvestmentWithLoanLock(ctx context.Context, investment *domain.Investment, loanID uuid.UUID) (*domain.Loan, error) {
	var loan domain.Loan

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock the loan within the transaction (SELECT FOR UPDATE)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Preload("Borrower").
			Where("id = ?", loanID).
			First(&loan).Error; err != nil {
			return err
		}

		// 2. Verify loan is still in approved state
		if loan.State != domain.LoanStateApproved {
			return domain.ErrInvalidLoanState
		}

		// 3. Check if investment still fits within remaining amount
		if investment.Amount > loan.RemainingInvestment {
			return domain.ErrInvestmentExceedsLimit
		}

		// 4. Update loan amounts
		loan.InvestedAmount += investment.Amount
		loan.RemainingInvestment -= investment.Amount

		// 5. Check if loan is fully funded
		if loan.RemainingInvestment <= 0 {
			loan.State = domain.LoanStateInvested
			loan.RemainingInvestment = 0 // Ensure it's exactly 0
		}

		// 6. Create the investment
		if err := tx.Create(investment).Error; err != nil {
			return err
		}

		// 7. Update the loan
		if err := tx.Save(&loan).Error; err != nil {
			return err
		}

		// 8. Update investor total invested
		if err := tx.Model(&domain.Investor{}).
			Where("id = ?", investment.InvestorID).
			Update("total_invested", gorm.Expr("total_invested + ?", investment.Amount)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &loan, nil
}

type approvalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) domain.ApprovalRepository {
	return &approvalRepository{db: db}
}

func (r *approvalRepository) Create(ctx context.Context, approval *domain.Approval) error {
	return r.db.WithContext(ctx).Create(approval).Error
}

func (r *approvalRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Approval, error) {
	var approval domain.Approval
	err := r.db.WithContext(ctx).
		Preload("Validator").
		Where("loan_id = ?", loanID).
		First(&approval).Error
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

type disbursementRepository struct {
	db *gorm.DB
}

func NewDisbursementRepository(db *gorm.DB) domain.DisbursementRepository {
	return &disbursementRepository{db: db}
}

func (r *disbursementRepository) Create(ctx context.Context, disbursement *domain.Disbursement) error {
	return r.db.WithContext(ctx).Create(disbursement).Error
}

func (r *disbursementRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Disbursement, error) {
	var disbursement domain.Disbursement
	err := r.db.WithContext(ctx).
		Preload("Officer").
		Where("loan_id = ?", loanID).
		First(&disbursement).Error
	if err != nil {
		return nil, err
	}
	return &disbursement, nil
}
