package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type notificationService struct {
	loanRepo       domain.LoanRepository
	investmentRepo domain.InvestmentRepository
}

func NewNotificationService(
	loanRepo domain.LoanRepository,
	investmentRepo domain.InvestmentRepository,
) domain.NotificationService {
	return &notificationService{
		loanRepo:       loanRepo,
		investmentRepo: investmentRepo,
	}
}

func (s *notificationService) SendAgreementLetters(ctx context.Context, loanID uuid.UUID) error {
	// Get all investments for this loan
	investments, err := s.investmentRepo.GetByLoanID(ctx, loanID)
	if err != nil {
		return fmt.Errorf("failed to get investments: %w", err)
	}

	// Generate agreement letter URL for each investment and simulate email sending
	for _, investment := range investments {
		// Generate dummy PDF URL for each investor's agreement letter
		agreementURL := s.generateAgreementLetterURL(loanID, investment.InvestorID, investment.ID)

		// Update investment with agreement letter URL
		err := s.investmentRepo.UpdateAgreementLetterURL(ctx, investment.ID, agreementURL)
		if err != nil {
			log.Printf("Failed to update agreement letter URL for investment %s: %v", investment.ID, err)
			continue
		}

		// Simulate sending email
		s.simulateEmailSending(investment.Investor.User.Email, investment.Investor.FullName, agreementURL, loanID)
	}

	return nil
}

// generateAgreementLetterURL generates a dummy PDF URL for the agreement letter
func (s *notificationService) generateAgreementLetterURL(loanID, investorID, investmentID uuid.UUID) string {
	return fmt.Sprintf("https://amf-documents.s3.amazonaws.com/agreements/loan_%s/investor_%s/agreement_%s.pdf",
		loanID.String(), investorID.String(), investmentID.String())
}

// simulateEmailSending logs the email that would be sent to the investor
func (s *notificationService) simulateEmailSending(email, fullName, agreementURL string, loanID uuid.UUID) {
	log.Printf("SIMULATED EMAIL SENT")
	log.Printf("To: %s (%s)", email, fullName)
	log.Printf("Subject: Investment Agreement Letter - Loan %s", loanID.String())
	log.Printf("Body: Dear %s,", fullName)
	log.Printf("Your investment has been successfully processed. Please find your agreement letter at:")
	log.Printf("Agreement Letter: %s", agreementURL)
	log.Printf("Thank you for investing with AMF Loan Service!")
	log.Printf("---")
}
