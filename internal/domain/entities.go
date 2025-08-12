package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleBorrower       UserRole = "borrower"
	RoleInvestor       UserRole = "investor"
	RoleFieldOfficer   UserRole = "field_officer"
	RoleFieldValidator UserRole = "field_validator"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      UserRole  `json:"role" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Borrower entity for storing borrower-specific information
type Borrower struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID `json:"user_id" gorm:"not null;unique"`
	FullName       string    `json:"full_name" gorm:"not null"`
	PhoneNumber    string    `json:"phone_number" gorm:"not null"`
	Address        string    `json:"address" gorm:"not null"`
	IdentityNumber string    `json:"identity_number" gorm:"not null;unique"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	User  User   `json:"user" gorm:"foreignKey:UserID"`
	Loans []Loan `json:"loans,omitempty" gorm:"foreignKey:BorrowerID"`
}

// Investor entity for storing investor-specific information
type Investor struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID `json:"user_id" gorm:"not null;unique"`
	FullName       string    `json:"full_name" gorm:"not null"`
	PhoneNumber    string    `json:"phone_number" gorm:"not null"`
	Address        string    `json:"address" gorm:"not null"`
	IdentityNumber string    `json:"identity_number" gorm:"not null;unique"`
	TotalInvested  float64   `json:"total_invested" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	User        User         `json:"user" gorm:"foreignKey:UserID"`
	Investments []Investment `json:"investments,omitempty" gorm:"foreignKey:InvestorID"`
}

type LoanState string

const (
	LoanStateProposed  LoanState = "proposed"
	LoanStateApproved  LoanState = "approved"
	LoanStateInvested  LoanState = "invested"
	LoanStateDisbursed LoanState = "disbursed"
)

type Loan struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BorrowerID          uuid.UUID `json:"borrower_id" gorm:"not null"`
	PrincipalAmount     float64   `json:"principal_amount" gorm:"not null"`
	InvestedAmount      float64   `json:"invested_amount" gorm:"default:0"`
	RemainingInvestment float64   `json:"remaining_investment" gorm:"not null"`
	Rate                float64   `json:"rate" gorm:"not null"`           // Interest rate for borrower
	ROI                 float64   `json:"roi" gorm:"not null"`            // Return on investment for investors (calculated)
	TotalInterest       float64   `json:"total_interest" gorm:"not null"` // Total interest borrower must pay
	State               LoanState `json:"state" gorm:"not null;default:'proposed'"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Relations
	Borrower     Borrower      `json:"borrower" gorm:"foreignKey:BorrowerID"`
	Approval     *Approval     `json:"approval,omitempty"`
	Investments  []Investment  `json:"investments,omitempty"`
	Disbursement *Disbursement `json:"disbursement,omitempty"`
}

type Approval struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LoanID        uuid.UUID `json:"loan_id" gorm:"not null"`
	ValidatorID   uuid.UUID `json:"validator_id" gorm:"not null"`
	PhotoProofURL string    `json:"photo_proof_url" gorm:"not null"`
	ApprovalDate  time.Time `json:"approval_date" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`

	// Relations
	Loan      Loan `json:"loan" gorm:"foreignKey:LoanID"`
	Validator User `json:"validator" gorm:"foreignKey:ValidatorID"`
}

type Investment struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LoanID             uuid.UUID `json:"loan_id" gorm:"not null"`
	InvestorID         uuid.UUID `json:"investor_id" gorm:"not null"`
	Amount             float64   `json:"amount" gorm:"not null"`
	Status             string    `json:"status" gorm:"default:'pending'"` // pending, completed, failed
	AgreementLetterURL string    `json:"agreement_letter_url"`            // PDF link for the investor
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Relations
	Loan     Loan     `json:"loan" gorm:"foreignKey:LoanID"`
	Investor Investor `json:"investor" gorm:"foreignKey:InvestorID"`
}

type Disbursement struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LoanID           uuid.UUID `json:"loan_id" gorm:"not null"`
	OfficerID        uuid.UUID `json:"officer_id" gorm:"not null"`
	AgreementFileURL string    `json:"agreement_file_url" gorm:"not null"`
	DisbursementDate time.Time `json:"disbursement_date" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at"`

	// Relations
	Loan    Loan `json:"loan" gorm:"foreignKey:LoanID"`
	Officer User `json:"officer" gorm:"foreignKey:OfficerID"`
}

// Investment event for Kafka
type InvestmentEvent struct {
	ID         uuid.UUID `json:"id"`
	LoanID     uuid.UUID `json:"loan_id"`
	InvestorID uuid.UUID `json:"investor_id"`
	Amount     float64   `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
}
