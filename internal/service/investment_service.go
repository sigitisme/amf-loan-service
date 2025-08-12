package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type investmentService struct {
	investmentRepo      domain.InvestmentRepository
	loanRepo            domain.LoanRepository
	investorRepo        domain.InvestorRepository
	kafkaProducer       domain.KafkaProducer
	notificationService domain.NotificationService
}

func NewInvestmentService(
	investmentRepo domain.InvestmentRepository,
	loanRepo domain.LoanRepository,
	investorRepo domain.InvestorRepository,
	kafkaProducer domain.KafkaProducer,
	notificationService domain.NotificationService,
) domain.InvestmentService {
	return &investmentService{
		investmentRepo:      investmentRepo,
		loanRepo:            loanRepo,
		investorRepo:        investorRepo,
		kafkaProducer:       kafkaProducer,
		notificationService: notificationService,
	}
}

// RequestInvestment validates the request and publishes to Kafka
func (s *investmentService) RequestInvestment(ctx context.Context, userID uuid.UUID, loanID uuid.UUID, amount float64) error {
	// Get investor to validate existence (userID is actually userID from the handler)
	investor, err := s.investorRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return err
	}

	// Get loan to validate (without lock, just for validation)
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrLoanNotFound
		}
		return err
	}

	// Check if loan is approved (allow investment only for approved)
	if loan.State != domain.LoanStateApproved {
		return domain.ErrLoanNotApproved
	}

	// Check if investor is trying to invest in their own loan (compare user IDs)
	if loan.Borrower.UserID == userID {
		return domain.ErrSelfInvestment
	}

	// Validate investment amount
	if amount <= 0 {
		return domain.ErrInvalidInvestmentAmount
	}

	// Check if investment would exceed remaining amount (basic check, final check in consumer)
	if amount > loan.RemainingInvestment {
		return domain.ErrInvestmentExceedsLimit
	}

	// Create investment event using the actual investor ID
	event := domain.InvestmentEvent{
		ID:         uuid.New(),
		LoanID:     loanID,
		InvestorID: investor.ID, // Use the actual investor ID, not user ID
		Amount:     amount,
		Timestamp:  time.Now(),
	}

	// Publish to Kafka for processing
	if err := s.kafkaProducer.PublishInvestmentEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to publish investment event: %w", err)
	}

	return nil
}

// ProcessInvestment handles the actual investment processing with transaction and locking
func (s *investmentService) ProcessInvestment(ctx context.Context, event domain.InvestmentEvent) error {
	// Get loan with pessimistic lock
	loan, err := s.loanRepo.GetByIDWithLock(ctx, event.LoanID)
	if err != nil {
		return fmt.Errorf("failed to get loan with lock: %w", err)
	}

	// Verify loan is still in approved state
	if loan.State != domain.LoanStateApproved {
		return fmt.Errorf("loan is no longer in approved state: %s", loan.State)
	}

	// Check if investment still fits within remaining amount
	if event.Amount > loan.RemainingInvestment {
		return domain.ErrInvestmentExceedsLimit
	}

	// Create investment record
	investment := &domain.Investment{
		ID:         event.ID,
		LoanID:     event.LoanID,
		InvestorID: event.InvestorID,
		Amount:     event.Amount,
		Status:     "completed",
		CreatedAt:  event.Timestamp,
		UpdatedAt:  time.Now(),
	}

	// Update loan amounts
	loan.InvestedAmount += event.Amount
	loan.RemainingInvestment -= event.Amount
	loan.UpdatedAt = time.Now()

	// Check if loan is fully funded
	if loan.RemainingInvestment <= 0 {
		loan.State = domain.LoanStateInvested
		loan.RemainingInvestment = 0 // Ensure it's exactly 0
	}

	// Execute transaction with both investment creation and loan update
	if err := s.investmentRepo.CreateWithTx(ctx, investment, loan); err != nil {
		return fmt.Errorf("failed to create investment with transaction: %w", err)
	}

	// If loan is fully funded, publish fully funded event and send agreement letters
	if loan.State == domain.LoanStateInvested {
		if s.kafkaProducer != nil {
			if err := s.kafkaProducer.PublishFullyFundedLoan(ctx, loan); err != nil {
				// Log error but don't fail the investment
				fmt.Printf("Failed to publish fully funded loan event: %v\n", err)
			}
		}

		// Send agreement letters to all investors
		if s.notificationService != nil {
			if err := s.notificationService.SendAgreementLetters(ctx, loan.ID); err != nil {
				// Log error but don't fail the investment
				fmt.Printf("Failed to send agreement letters: %v\n", err)
			}
		}
	}

	return nil
}

func (s *investmentService) GetInvestorInvestments(ctx context.Context, investorID uuid.UUID) ([]domain.Investment, error) {
	return s.investmentRepo.GetByInvestorID(ctx, investorID)
}

func (s *investmentService) GetInvestorInvestmentsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Investment, error) {
	// Get investor by user ID first
	investor, err := s.investorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get investments by investor ID
	return s.investmentRepo.GetByInvestorID(ctx, investor.ID)
}

func (s *investmentService) GetLoanInvestments(ctx context.Context, loanID uuid.UUID) ([]domain.Investment, error) {
	return s.investmentRepo.GetByLoanID(ctx, loanID)
}
