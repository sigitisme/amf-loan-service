package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/sigitisme/amf-loan-service/internal/config"
	"github.com/sigitisme/amf-loan-service/internal/domain"
)

type Producer struct {
	investmentWriter  *kafka.Writer
	fullyFundedWriter *kafka.Writer
}

func NewProducer(cfg *config.KafkaConfig) *Producer {
	investmentWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.InvestmentTopic,
		Balancer: &kafka.LeastBytes{},
	}

	fullyFundedWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.FullyFundedTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		investmentWriter:  investmentWriter,
		fullyFundedWriter: fullyFundedWriter,
	}
}

func (p *Producer) PublishInvestmentEvent(ctx context.Context, event domain.InvestmentEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.LoanID.String()),
		Value: data,
	}

	err = p.investmentWriter.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Error publishing investment event: %v", err)
		return err
	}

	log.Printf("Investment event published for loan: %s, amount: %.2f", event.LoanID, event.Amount)
	return nil
}

func (p *Producer) PublishInvestment(ctx context.Context, investment *domain.Investment) error {
	data, err := json.Marshal(investment)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(investment.LoanID.String()),
		Value: data,
	}

	err = p.investmentWriter.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Error publishing investment message: %v", err)
		return err
	}

	log.Printf("Investment message published for loan: %s", investment.LoanID)
	return nil
}

func (p *Producer) PublishFullyFundedLoan(ctx context.Context, loan *domain.Loan) error {
	data, err := json.Marshal(loan)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(loan.ID.String()),
		Value: data,
	}

	err = p.fullyFundedWriter.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Error publishing fully funded loan message: %v", err)
		return err
	}

	log.Printf("Fully funded loan message published for loan: %s", loan.ID)
	return nil
}

func (p *Producer) Close() {
	if p.investmentWriter != nil {
		p.investmentWriter.Close()
	}
	if p.fullyFundedWriter != nil {
		p.fullyFundedWriter.Close()
	}
}
