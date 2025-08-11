package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

// ============================================================================
// AUTH DTOs
// ============================================================================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token    string            `json:"token"`
	User     UserResponse      `json:"user"`
	Borrower *BorrowerResponse `json:"borrower,omitempty"`
	Investor *InvestorResponse `json:"investor,omitempty"`
}

type UserResponse struct {
	ID    uuid.UUID       `json:"id"`
	Email string          `json:"email"`
	Role  domain.UserRole `json:"role"`
}

// ============================================================================
// BORROWER DTOs
// ============================================================================

type BorrowerResponse struct {
	ID             uuid.UUID     `json:"id"`
	UserID         uuid.UUID     `json:"user_id"`
	FullName       string        `json:"full_name"`
	PhoneNumber    string        `json:"phone_number"`
	Address        string        `json:"address"`
	IdentityNumber string        `json:"identity_number"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	User           *UserResponse `json:"user,omitempty"`
}

// ============================================================================
// INVESTOR DTOs
// ============================================================================

type InvestorResponse struct {
	ID             uuid.UUID     `json:"id"`
	UserID         uuid.UUID     `json:"user_id"`
	FullName       string        `json:"full_name"`
	PhoneNumber    string        `json:"phone_number"`
	Address        string        `json:"address"`
	IdentityNumber string        `json:"identity_number"`
	TotalInvested  float64       `json:"total_invested"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	User           *UserResponse `json:"user,omitempty"`
}

// ============================================================================
// LOAN DTOs
// ============================================================================

type CreateLoanRequest struct {
	PrincipalAmount float64 `json:"principal_amount" binding:"required,min=1000"`
	Rate            float64 `json:"rate" binding:"required,min=0.01,max=1"`
}

type LoanResponse struct {
	ID                  uuid.UUID        `json:"id"`
	BorrowerID          uuid.UUID        `json:"borrower_id"`
	PrincipalAmount     float64          `json:"principal_amount"`
	InvestedAmount      float64          `json:"invested_amount"`
	RemainingInvestment float64          `json:"remaining_investment"`
	Rate                float64          `json:"rate"`
	ROI                 float64          `json:"roi"`
	TotalInterest       float64          `json:"total_interest"`
	State               domain.LoanState `json:"state"`
	AgreementLetterURL  string           `json:"agreement_letter_url,omitempty"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
	// Related data - only included when requested
	Borrower    *BorrowerResponse    `json:"borrower,omitempty"`
	Investments []InvestmentResponse `json:"investments,omitempty"`
}

type ApproveLoanRequest struct {
	PhotoProofURL string    `json:"photo_proof_url" binding:"required,url"`
	ApprovalDate  time.Time `json:"approval_date" binding:"required"`
}

type DisburseLoanRequest struct {
	AgreementFileURL string    `json:"agreement_file_url" binding:"required,url"`
	DisbursementDate time.Time `json:"disbursement_date" binding:"required"`
}

// ============================================================================
// INVESTMENT DTOs
// ============================================================================

type InvestRequest struct {
	LoanID uuid.UUID `json:"loan_id" binding:"required"`
	Amount float64   `json:"amount" binding:"required,min=1000"`
}

type InvestmentResponse struct {
	ID         uuid.UUID `json:"id"`
	LoanID     uuid.UUID `json:"loan_id"`
	InvestorID uuid.UUID `json:"investor_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// Related data - only included when requested
	Loan     *LoanResponse     `json:"loan,omitempty"`
	Investor *InvestorResponse `json:"investor,omitempty"`
}

// ============================================================================
// PAGINATION & FILTERING DTOs
// ============================================================================

type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type LoansFilter struct {
	PaginationRequest
	State      domain.LoanState `form:"state"`
	BorrowerID uuid.UUID        `form:"borrower_id"`
	MinAmount  float64          `form:"min_amount"`
	MaxAmount  float64          `form:"max_amount"`
}

// ============================================================================
// API RESPONSE WRAPPERS
// ============================================================================

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool               `json:"success"`
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
