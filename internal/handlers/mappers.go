package handlers

import (
	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

// ============================================================================
// USER MAPPERS
// ============================================================================

func MapUserToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
}

// ============================================================================
// BORROWER MAPPERS
// ============================================================================

func MapBorrowerToResponse(borrower *domain.Borrower) BorrowerResponse {
	response := BorrowerResponse{
		ID:             borrower.ID,
		UserID:         borrower.UserID,
		FullName:       borrower.FullName,
		PhoneNumber:    borrower.PhoneNumber,
		Address:        borrower.Address,
		IdentityNumber: borrower.IdentityNumber,
		CreatedAt:      borrower.CreatedAt,
		UpdatedAt:      borrower.UpdatedAt,
	}

	// Include user if loaded (has valid ID)
	if borrower.User.ID != uuid.Nil {
		userResp := MapUserToResponse(&borrower.User)
		response.User = &userResp
	}

	return response
}

// ============================================================================
// INVESTOR MAPPERS
// ============================================================================

func MapInvestorToResponse(investor *domain.Investor) InvestorResponse {
	response := InvestorResponse{
		ID:             investor.ID,
		UserID:         investor.UserID,
		FullName:       investor.FullName,
		PhoneNumber:    investor.PhoneNumber,
		Address:        investor.Address,
		IdentityNumber: investor.IdentityNumber,
		TotalInvested:  investor.TotalInvested,
		CreatedAt:      investor.CreatedAt,
		UpdatedAt:      investor.UpdatedAt,
	}

	// Include user if loaded (has valid ID)
	if investor.User.ID != uuid.Nil {
		userResp := MapUserToResponse(&investor.User)
		response.User = &userResp
	}

	return response
}

// ============================================================================
// LOAN MAPPERS
// ============================================================================

func MapLoanToResponse(loan *domain.Loan, includeBorrower, includeInvestments bool) LoanResponse {
	response := LoanResponse{
		ID:                  loan.ID,
		BorrowerID:          loan.BorrowerID,
		PrincipalAmount:     loan.PrincipalAmount,
		InvestedAmount:      loan.InvestedAmount,
		RemainingInvestment: loan.RemainingInvestment,
		Rate:                loan.Rate,
		ROI:                 loan.ROI,
		TotalInterest:       loan.TotalInterest,
		State:               loan.State,
		CreatedAt:           loan.CreatedAt,
		UpdatedAt:           loan.UpdatedAt,
	}

	// Include borrower if requested and actually loaded (has valid ID)
	if includeBorrower && loan.Borrower.ID != uuid.Nil {
		borrowerResp := MapBorrowerToResponse(&loan.Borrower)
		response.Borrower = &borrowerResp
	}

	// Include investments if requested and loaded
	if includeInvestments && len(loan.Investments) > 0 {
		response.Investments = MapInvestmentsToResponse(loan.Investments, false, false)
	}

	return response
}

// ============================================================================
// INVESTMENT MAPPERS
// ============================================================================

func MapInvestmentToResponse(investment *domain.Investment, includeLoan, includeInvestor bool) InvestmentResponse {
	response := InvestmentResponse{
		ID:         investment.ID,
		LoanID:     investment.LoanID,
		InvestorID: investment.InvestorID,
		Amount:     investment.Amount,
		CreatedAt:  investment.CreatedAt,
		UpdatedAt:  investment.UpdatedAt,
	}

	// Include loan if requested and loaded (has valid ID)
	if includeLoan && investment.Loan.ID != uuid.Nil {
		loanResp := MapLoanToResponse(&investment.Loan, false, false)
		response.Loan = &loanResp
	}

	// Include investor if requested and loaded (has valid ID)
	if includeInvestor && investment.Investor.ID != uuid.Nil {
		investorResp := MapInvestorToResponse(&investment.Investor)
		response.Investor = &investorResp
	}

	return response
}

// ============================================================================
// COLLECTION MAPPERS
// ============================================================================

func MapLoansToResponse(loans []domain.Loan, includeBorrower, includeInvestments bool) []LoanResponse {
	responses := make([]LoanResponse, len(loans))
	for i, loan := range loans {
		responses[i] = MapLoanToResponse(&loan, includeBorrower, includeInvestments)
	}
	return responses
}

func MapInvestmentsToResponse(investments []domain.Investment, includeLoan, includeInvestor bool) []InvestmentResponse {
	responses := make([]InvestmentResponse, len(investments))
	for i, investment := range investments {
		responses[i] = MapInvestmentToResponse(&investment, includeLoan, includeInvestor)
	}
	return responses
}

// ============================================================================
// HELPER FUNCTIONS FOR API RESPONSES
// ============================================================================

func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func SuccessResponseWithMessage(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponseFromString(err string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
	}
}

func ErrorResponseWithMessage(err, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
		Message: message,
	}
}

func PaginatedSuccessResponse(data interface{}, pagination PaginationResponse) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	}
}

// ============================================================================
// PAGINATION HELPERS
// ============================================================================

func CalculatePagination(page, pageSize int, totalItems int64) PaginationResponse {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

func GetOffsetAndLimit(page, pageSize int) (offset int, limit int) {
	offset = (page - 1) * pageSize
	limit = pageSize
	return offset, limit
}
