package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"gorm.io/gorm"
)

type investorRepository struct {
	db *gorm.DB
}

func NewInvestorRepository(db *gorm.DB) domain.InvestorRepository {
	return &investorRepository{db: db}
}

func (r *investorRepository) Create(ctx context.Context, investor *domain.Investor) error {
	return r.db.WithContext(ctx).Create(investor).Error
}

func (r *investorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Investor, error) {
	var investor domain.Investor
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&investor).Error
	if err != nil {
		return nil, err
	}
	return &investor, nil
}

func (r *investorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Investor, error) {
	var investor domain.Investor
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&investor).Error
	if err != nil {
		return nil, err
	}
	return &investor, nil
}

func (r *investorRepository) Update(ctx context.Context, investor *domain.Investor) error {
	return r.db.WithContext(ctx).Save(investor).Error
}
