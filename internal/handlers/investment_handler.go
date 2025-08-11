package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type InvestmentHandler struct {
	investmentService domain.InvestmentService
}

func NewInvestmentHandler(investmentService domain.InvestmentService) *InvestmentHandler {
	return &InvestmentHandler{
		investmentService: investmentService,
	}
}

func (h *InvestmentHandler) Invest(c *gin.Context) {
	var req InvestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "validation_failed",
			Message: err.Error(),
		})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "unauthorized",
			Message: "User not found in context",
		})
		return
	}

	userObj, ok := user.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "internal_error",
			Message: "Invalid user type",
		})
		return
	}

	// Only investors can invest
	if userObj.Role != domain.RoleInvestor {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Success: false,
			Error:   "forbidden",
			Message: "Only investors can invest in loans",
		})
		return
	}

	// Convert handler DTO to service parameters
	err := h.investmentService.RequestInvestment(c.Request.Context(), userObj.ID, req.LoanID, req.Amount)
	if err != nil {
		switch err {
		case domain.ErrLoanNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error:   "loan_not_found",
				Message: "The specified loan was not found",
			})
		case domain.ErrInvalidLoanState:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "invalid_loan_state",
				Message: "Loan is not available for investment",
			})
		case domain.ErrInvestmentExceedsLimit:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "insufficient_remaining",
				Message: "Investment amount exceeds remaining loan amount",
			})
		case domain.ErrInvalidInvestmentAmount:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "invalid_amount",
				Message: err.Error(),
			})
		case domain.ErrSelfInvestment:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "self_investment",
				Message: "Borrowers cannot invest in their own loans",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "investment_failed",
				Message: "Failed to process investment request",
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, SuccessResponseWithMessage("Investment request submitted for processing", nil))
}

func (h *InvestmentHandler) GetMyInvestments(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "unauthorized",
			Message: "User not found in context",
		})
		return
	}

	userObj, ok := user.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "internal_error",
			Message: "Invalid user type",
		})
		return
	}

	// Only investors can view investments
	if userObj.Role != domain.RoleInvestor {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Success: false,
			Error:   "forbidden",
			Message: "Only investors can view investments",
		})
		return
	}

	investments, err := h.investmentService.GetInvestorInvestments(c.Request.Context(), userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "fetch_failed",
			Message: "Failed to fetch investments",
		})
		return
	}

	// Convert domain entities to handler responses
	responses := MapInvestmentsToResponse(investments, true, false)
	c.JSON(http.StatusOK, responses)
}

func (h *InvestmentHandler) GetLoanInvestments(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := uuid.Parse(loanIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "invalid_id",
			Message: "Invalid loan ID format",
		})
		return
	}

	investments, err := h.investmentService.GetLoanInvestments(c.Request.Context(), loanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "fetch_failed",
			Message: "Failed to fetch loan investments",
		})
		return
	}

	// Convert domain entities to handler responses
	responses := MapInvestmentsToResponse(investments, false, true)
	c.JSON(http.StatusOK, responses)
}
