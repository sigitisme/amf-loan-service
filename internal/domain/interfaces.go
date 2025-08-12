package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LoginResponse represents the response returned after a successful login
type LoginResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Repository interfaces for clean architecture

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BorrowerRepository interface {
	Create(ctx context.Context, borrower *Borrower) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Borrower, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Borrower, error)
	Update(ctx context.Context, borrower *Borrower) error
}

type InvestorRepository interface {
	Create(ctx context.Context, investor *Investor) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Investor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Investor, error)
	Update(ctx context.Context, investor *Investor) error
}

type LoanRepository interface {
	Create(ctx context.Context, loan *Loan) error
	GetByID(ctx context.Context, id uuid.UUID) (*Loan, error)
	GetByIDWithLock(ctx context.Context, id uuid.UUID) (*Loan, error) // For pessimistic locking
	GetByBorrowerID(ctx context.Context, borrowerID uuid.UUID) ([]Loan, error)
	GetByState(ctx context.Context, state LoanState) ([]Loan, error)
	Update(ctx context.Context, loan *Loan) error
	List(ctx context.Context, limit, offset int) ([]Loan, error)
}

type ApprovalRepository interface {
	Create(ctx context.Context, approval *Approval) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) (*Approval, error)
}

type InvestmentRepository interface {
	Create(ctx context.Context, investment *Investment) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) ([]Investment, error)
	GetByInvestorID(ctx context.Context, investorID uuid.UUID) ([]Investment, error)
	GetTotalInvestedAmount(ctx context.Context, loanID uuid.UUID) (float64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	CreateWithTx(ctx context.Context, investment *Investment, loan *Loan) error // Transaction method
}

type DisbursementRepository interface {
	Create(ctx context.Context, disbursement *Disbursement) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) (*Disbursement, error)
}

// Service interfaces

type AuthService interface {
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	ValidateToken(tokenString string) (*User, error)
}

type LoanService interface {
	CreateLoan(ctx context.Context, borrowerID uuid.UUID, principalAmount, rate float64) (*Loan, error)
	ApproveLoan(ctx context.Context, loanID uuid.UUID, validatorID uuid.UUID, photoProofURL string, approvalDate time.Time) error
	GetLoansByState(ctx context.Context, state LoanState) ([]Loan, error)
	GetLoanByID(ctx context.Context, id uuid.UUID) (*Loan, error)
	GetBorrowerLoans(ctx context.Context, borrowerID uuid.UUID) ([]Loan, error)
	GetBorrowerLoansByUserID(ctx context.Context, userID uuid.UUID) ([]Loan, error)
	DisburseLoan(ctx context.Context, loanID uuid.UUID, officerID uuid.UUID, agreementFileURL string, disbursementDate time.Time) error
}

type InvestmentService interface {
	RequestInvestment(ctx context.Context, investorID uuid.UUID, loanID uuid.UUID, amount float64) error // Just validate and publish
	ProcessInvestment(ctx context.Context, event InvestmentEvent) error                                  // Consumer logic
	GetInvestorInvestments(ctx context.Context, investorID uuid.UUID) ([]Investment, error)
	GetInvestorInvestmentsByUserID(ctx context.Context, userID uuid.UUID) ([]Investment, error)
	GetLoanInvestments(ctx context.Context, loanID uuid.UUID) ([]Investment, error)
}

type NotificationService interface {
	SendAgreementLetters(ctx context.Context, loanID uuid.UUID) error
}

type KafkaProducer interface {
	PublishInvestmentEvent(ctx context.Context, event InvestmentEvent) error
	PublishFullyFundedLoan(ctx context.Context, loan *Loan) error
}

type InvestmentConsumer interface {
	StartConsumer(ctx context.Context) error
	StopConsumer() error
}
