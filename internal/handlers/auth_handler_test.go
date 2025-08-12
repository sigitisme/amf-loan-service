package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Auth Service
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginResponse), args.Error(1)
}

func (m *mockAuthService) ValidateToken(tokenString string) (*domain.User, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Test Auth Handler Login - Happy Flow
func TestAuthHandler_Login_Success(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Arrange
	mockAuthService := new(mockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResponse := &domain.LoginResponse{
		UserID:    uuid.New(),
		Email:     "test@example.com",
		Token:     "jwt-token-here",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockAuthService.On("Login", mock.Anything, loginReq.Email, loginReq.Password).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	authHandler.Login(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.Email, response.Email)
	assert.Equal(t, expectedResponse.Token, response.Token)

	mockAuthService.AssertExpectations(t)
}

// Test Auth Handler Login - Invalid Credentials
func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Arrange
	mockAuthService := new(mockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockAuthService.On("Login", mock.Anything, loginReq.Email, loginReq.Password).Return(nil, domain.ErrInvalidCredentials)

	// Create HTTP request
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	authHandler.Login(c)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid_credentials", response.Error)

	mockAuthService.AssertExpectations(t)
}

// Test Auth Handler Login - Invalid JSON
func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Arrange
	mockAuthService := new(mockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	authHandler.Login(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "validation_failed", response.Error)
}
