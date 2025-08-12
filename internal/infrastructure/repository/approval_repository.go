package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"gorm.io/gorm"
)

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
