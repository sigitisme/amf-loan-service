package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "validation_failed",
			Message: err.Error(),
		})
		return
	}

	// Convert handler DTO to service parameters
	domainResponse, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "login_failed",
				Message: "An error occurred during login",
			})
		}
		return
	}

	// Convert domain response to handler response
	response := domain.LoginResponse{
		Token:     domainResponse.Token,
		UserID:    domainResponse.UserID,
		Email:     domainResponse.Email,
		ExpiresAt: domainResponse.ExpiresAt,
	}
	c.JSON(http.StatusOK, response)
}
