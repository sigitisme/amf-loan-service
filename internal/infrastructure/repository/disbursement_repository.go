package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"gorm.io/gorm"
)

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
