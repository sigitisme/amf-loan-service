package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"gorm.io/gorm"
)

type borrowerRepository struct {
	db *gorm.DB
}

func NewBorrowerRepository(db *gorm.DB) domain.BorrowerRepository {
	return &borrowerRepository{db: db}
}

func (r *borrowerRepository) Create(ctx context.Context, borrower *domain.Borrower) error {
	return r.db.WithContext(ctx).Create(borrower).Error
}

func (r *borrowerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Borrower, error) {
	var borrower domain.Borrower
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&borrower).Error
	if err != nil {
		return nil, err
	}
	return &borrower, nil
}

func (r *borrowerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Borrower, error) {
	var borrower domain.Borrower
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&borrower).Error
	if err != nil {
		return nil, err
	}
	return &borrower, nil
}

func (r *borrowerRepository) Update(ctx context.Context, borrower *domain.Borrower) error {
	return r.db.WithContext(ctx).Save(borrower).Error
}
