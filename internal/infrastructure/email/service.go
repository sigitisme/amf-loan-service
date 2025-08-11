package email

import (
	"fmt"
	"log"

	"github.com/sigitisme/amf-loan-service/internal/config"
	"gopkg.in/gomail.v2"
)

type Service struct {
	smtpConfig *config.SMTPConfig
}

func NewService(cfg *config.SMTPConfig) *Service {
	return &Service{
		smtpConfig: cfg,
	}
}

func (s *Service) SendAgreementLetter(to, borrowerName, loanID, agreementURL string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.smtpConfig.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Loan Agreement Letter - Loan ID: "+loanID)

	body := fmt.Sprintf(`
Dear Investor,

Thank you for your investment in loan ID: %s for borrower %s.

The loan has been fully funded and is ready for disbursement. 
Please find your agreement letter at the following link:

%s

Best regards,
AMF Loan Service Team
`, loanID, borrowerName, agreementURL)

	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		s.smtpConfig.Host,
		587, // Convert string to int if needed
		s.smtpConfig.Username,
		s.smtpConfig.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Agreement letter sent to %s for loan %s", to, loanID)
	return nil
}
