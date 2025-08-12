package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/config"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockBorrowerRepository struct {
	mock.Mock
}

func (m *mockBorrowerRepository) Create(ctx context.Context, borrower *domain.Borrower) error {
	args := m.Called(ctx, borrower)
	return args.Error(0)
}

func (m *mockBorrowerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Borrower, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Borrower), args.Error(1)
}

func (m *mockBorrowerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Borrower, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Borrower), args.Error(1)
}

func (m *mockBorrowerRepository) Update(ctx context.Context, borrower *domain.Borrower) error {
	args := m.Called(ctx, borrower)
	return args.Error(0)
}

type mockInvestorRepository struct {
	mock.Mock
}

func (m *mockInvestorRepository) Create(ctx context.Context, investor *domain.Investor) error {
	args := m.Called(ctx, investor)
	return args.Error(0)
}

func (m *mockInvestorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Investor, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Investor), args.Error(1)
}

func (m *mockInvestorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Investor, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Investor), args.Error(1)
}

func (m *mockInvestorRepository) Update(ctx context.Context, investor *domain.Investor) error {
	args := m.Called(ctx, investor)
	return args.Error(0)
}

// Test AuthService Login - Happy Flow
func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockUserRepository)
	mockBorrowerRepo := new(mockBorrowerRepository)
	mockInvestorRepo := new(mockInvestorRepository)

	jwtConfig := &config.JWTConfig{
		Secret: "test-secret",
		Expiry: time.Hour,
	}

	authService := NewAuthService(mockUserRepo, mockBorrowerRepo, mockInvestorRepo, jwtConfig)

	userID := uuid.New()
	email := "test@example.com"
	hashedPassword := "$2a$14$hashedpassword" // Mock bcrypt hash

	user := &domain.User{
		ID:       userID,
		Email:    email,
		Password: hashedPassword,
		Role:     domain.RoleInvestor,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)

	// Act
	response, err := authService.Login(context.Background(), email, "password")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, email, response.Email)
	assert.NotEmpty(t, response.Token)
	assert.True(t, response.ExpiresAt.After(time.Now()))

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(mockUserRepository)
	mockBorrowerRepo := new(mockBorrowerRepository)
	mockInvestorRepo := new(mockInvestorRepository)

	jwtConfig := &config.JWTConfig{
		Secret: "test-secret",
		Expiry: time.Hour,
	}

	authService := NewAuthService(mockUserRepo, mockBorrowerRepo, mockInvestorRepo, jwtConfig)

	email := "nonexistent@example.com"

	mockUserRepo.On("GetByEmail", mock.Anything, email).Return(nil, domain.ErrInvalidCredentials)

	// Act
	response, err := authService.Login(context.Background(), email, "password")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, domain.ErrInvalidCredentials, err)

	mockUserRepo.AssertExpectations(t)
}
