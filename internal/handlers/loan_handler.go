package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type LoanHandler struct {
	loanService domain.LoanService
}

func NewLoanHandler(loanService domain.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "validation_failed",
			Message: err.Error(),
		})
		return
	}

	// Get user from context (set by auth middleware)
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

	// Only borrowers can create loans
	if userObj.Role != domain.RoleBorrower {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Success: false,
			Error:   "forbidden",
			Message: "Only borrowers can create loans",
		})
		return
	}

	// Convert handler DTO to service parameters
	loan, err := h.loanService.CreateLoan(c.Request.Context(), userObj.ID, req.PrincipalAmount, req.Rate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "creation_failed",
			Message: "Failed to create loan",
		})
		return
	}

	// Convert domain entity to handler response (don't include borrower for creation)
	response := MapLoanToResponse(loan, false, false)
	c.JSON(http.StatusCreated, response)
}

func (h *LoanHandler) ApproveLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := uuid.Parse(loanIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	var req ApproveLoanRequest
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

	// Only field validators can approve loans
	if userObj.Role != domain.RoleFieldValidator {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Success: false,
			Error:   "forbidden",
			Message: "Only field validators can approve loans",
		})
		return
	}

	// Convert handler DTO to service parameters
	err = h.loanService.ApproveLoan(c.Request.Context(), loanID, userObj.ID, req.PhotoProofURL, req.ApprovalDate)
	if err != nil {
		switch err {
		case domain.ErrLoanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrLoanAlreadyApproved:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve loan"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan approved successfully"})
}

func (h *LoanHandler) GetLoans(c *gin.Context) {
	stateStr := c.Query("state")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userObj, ok := user.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	var loans []domain.Loan
	var err error

	if userObj.Role == domain.RoleBorrower {
		// Borrowers can only see their own loans
		loans, err = h.loanService.GetBorrowerLoans(c.Request.Context(), userObj.ID)
	} else if stateStr != "" {
		// Staff members can filter by state
		state := domain.LoanState(stateStr)
		loans, err = h.loanService.GetLoansByState(c.Request.Context(), state)
	} else {
		// For staff without state filter, get approved loans
		loans, err = h.loanService.GetLoansByState(c.Request.Context(), domain.LoanStateApproved)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loans"})
		return
	}

	c.JSON(http.StatusOK, loans)
}

func (h *LoanHandler) GetLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := uuid.Parse(loanIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	loan, err := h.loanService.GetLoanByID(c.Request.Context(), loanID)
	if err != nil {
		switch err {
		case domain.ErrLoanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loan"})
		}
		return
	}

	c.JSON(http.StatusOK, loan)
}

func (h *LoanHandler) DisburseLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := uuid.Parse(loanIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	var req DisburseLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userObj, ok := user.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	// Only field officers can disburse loans
	if userObj.Role != domain.RoleFieldOfficer {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only field officers can disburse loans"})
		return
	}

	err = h.loanService.DisburseLoan(c.Request.Context(), loanID, userObj.ID, req.AgreementFileURL, req.DisbursementDate)
	if err != nil {
		switch err {
		case domain.ErrLoanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrLoanNotInvested:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disburse loan"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan disbursed successfully"})
}
