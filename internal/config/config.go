package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Kafka    KafkaConfig
	SMTP     SMTPConfig
	API      APIConfig
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type KafkaConfig struct {
	Brokers          []string
	InvestmentTopic  string
	FullyFundedTopic string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type APIConfig struct {
	Port string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	expiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		expiry = 24 * time.Hour
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "loan_service"),
			Port:     getEnv("DB_PORT", "5432"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-super-secret-key"),
			Expiry: expiry,
		},
		Kafka: KafkaConfig{
			Brokers:          []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			InvestmentTopic:  getEnv("KAFKA_INVESTMENT_TOPIC", "investment_processing"),
			FullyFundedTopic: getEnv("KAFKA_FULLY_FUNDED_TOPIC", "loan_fully_funded"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
		API: APIConfig{
			Port: getEnv("API_PORT", "8080"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
