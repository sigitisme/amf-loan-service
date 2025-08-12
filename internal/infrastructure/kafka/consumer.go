package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sigitisme/amf-loan-service/internal/config"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type Consumer struct {
	reader            *kafka.Reader
	investmentService domain.InvestmentService
	running           bool
}

func NewConsumer(cfg *config.KafkaConfig, investmentService domain.InvestmentService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.InvestmentTopic,
		GroupID:        "investment-processor",
		MinBytes:       1,                      // Process messages immediately, don't wait to batch
		MaxBytes:       10e6,                   // 10MB max
		CommitInterval: 100 * time.Millisecond, // Commit more frequently
		MaxWait:        100 * time.Millisecond, // Don't wait long for batching
		StartOffset:    kafka.LastOffset,       // Start from latest
	})

	return &Consumer{
		reader:            reader,
		investmentService: investmentService,
		running:           false,
	}
}

func (c *Consumer) StartConsumer(ctx context.Context) error {
	c.running = true
	log.Println("Starting investment event consumer...")

	for c.running {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled, shutting down...")
			return ctx.Err()
		default:
			// Add timeout for message fetching
			fetchCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			message, err := c.reader.FetchMessage(fetchCtx)
			cancel()

			if err != nil {
				if !c.running {
					return nil
				}
				// Don't log timeout errors, they're expected
				if err != context.DeadlineExceeded {
					log.Printf("Error fetching message: %v", err)
				}
				continue // Continue immediately without sleeping
			}

			if err := c.processMessage(ctx, message); err != nil {
				log.Printf("Error processing message: %v", err)
				// Don't commit the message if processing failed
				continue
			}

			if err := c.reader.CommitMessages(ctx, message); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}

	return nil
}

func (c *Consumer) StopConsumer() error {
	c.running = false
	log.Println("Stopping investment event consumer...")
	return c.reader.Close()
}

func (c *Consumer) processMessage(ctx context.Context, message kafka.Message) error {
	var event domain.InvestmentEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		log.Printf("Error unmarshaling investment event: %v", err)
		return err
	}

	log.Printf("Processing investment event: Loan %s, Investor %s, Amount %.2f",
		event.LoanID, event.InvestorID, event.Amount)

	// Process the investment with transaction and locking
	if err := c.investmentService.ProcessInvestment(ctx, event); err != nil {
		log.Printf("Error processing investment: %v", err)
		return err
	}

	log.Printf("Successfully processed investment event: %s", event.ID)
	return nil
}
