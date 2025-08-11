package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/sigitisme/amf-loan-service/internal/infrastructure/email"
)

type notificationService struct {
	loanRepo       domain.LoanRepository
	investmentRepo domain.InvestmentRepository
	emailService   *email.Service
}

func NewNotificationService(
	loanRepo domain.LoanRepository,
	investmentRepo domain.InvestmentRepository,
	emailService *email.Service,
) domain.NotificationService {
	return &notificationService{
		loanRepo:       loanRepo,
		investmentRepo: investmentRepo,
		emailService:   emailService,
	}
}

func (s *notificationService) SendAgreementLetters(ctx context.Context, loanID uuid.UUID) error {
	// Get loan details
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return fmt.Errorf("failed to get loan: %w", err)
	}

	// Get all investments for this loan
	investments, err := s.investmentRepo.GetByLoanID(ctx, loanID)
	if err != nil {
		return fmt.Errorf("failed to get investments: %w", err)
	}

	// Send email to each investor
	for _, investment := range investments {
		err = s.emailService.SendAgreementLetter(
			investment.Investor.User.Email,
			loan.Borrower.User.Email, // Using borrower email as name for now
			loanID.String(),
			loan.AgreementLetterURL,
		)
		if err != nil {
			// Log error but continue with other investors
			fmt.Printf("Failed to send agreement letter to %s: %v\n", investment.Investor.User.Email, err)
		}
	}

	return nil
}
