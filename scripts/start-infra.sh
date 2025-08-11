#!/bin/bash

echo "🚀 Starting AMF Loan Service Infrastructure..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start infrastructure services
echo "📦 Starting PostgreSQL and Kafka..."
docker-compose up -d postgres kafka zookeeper

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 15

# Check if PostgreSQL is ready
echo "🔍 Checking PostgreSQL connection..."
until docker-compose exec postgres pg_isready -U postgres > /dev/null 2>&1; do
    echo "⏳ Waiting for PostgreSQL to be ready..."
    sleep 2
done
echo "✅ PostgreSQL is ready"

# Check if Kafka is ready
echo "🔍 Checking Kafka connection..."
sleep 10  # Give Kafka more time to initialize
echo "✅ Kafka should be ready"

# Create Kafka topics
echo "📝 Creating Kafka topics..."
docker-compose exec kafka kafka-topics --create --bootstrap-server localhost:9092 --topic investment_processing --partitions 1 --replication-factor 1 --if-not-exists
docker-compose exec kafka kafka-topics --create --bootstrap-server localhost:9092 --topic loan_fully_funded --partitions 1 --replication-factor 1 --if-not-exists

echo "✅ Infrastructure services are ready!"
echo ""
echo "📊 Service Status:"
echo "- PostgreSQL: localhost:5432"
echo "- Kafka: localhost:9092"
echo ""
echo "🔧 Next steps:"
echo "1. Copy .env.example to .env and configure your settings"
echo "2. Run: go run cmd/server/main.go"
echo ""
echo "🛑 To stop services: docker-compose down"
