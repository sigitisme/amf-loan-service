package domain

import "errors"

var (
	// Auth errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidToken       = errors.New("invalid token")

	// Loan errors
	ErrLoanNotFound         = errors.New("loan not found")
	ErrLoanAlreadyApproved  = errors.New("loan is already approved")
	ErrLoanNotApproved      = errors.New("loan is not approved yet")
	ErrLoanAlreadyInvested  = errors.New("loan is already fully invested")
	ErrLoanNotInvested      = errors.New("loan is not fully invested yet")
	ErrLoanAlreadyDisbursed = errors.New("loan is already disbursed")
	ErrInvalidLoanState     = errors.New("invalid loan state for this operation")

	// Investment errors
	ErrInvestmentExceedsLimit  = errors.New("investment amount exceeds remaining loan amount")
	ErrInvalidInvestmentAmount = errors.New("investment amount must be greater than 0")
	ErrSelfInvestment          = errors.New("borrower cannot invest in their own loan")

	// Permission errors
	ErrInsufficientPermission = errors.New("insufficient permission for this operation")
	ErrInvalidRole            = errors.New("invalid role for this operation")
)
