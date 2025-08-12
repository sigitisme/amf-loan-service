#!/bin/bash

# AMF Loan Service - Infrastructure Cleanup Script
# This script removes all Docker containers, volumes, and data

set -e

echo "🧹 AMF Loan Service Infrastructure Cleanup"
echo "========================================="

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "📍 Project directory: $PROJECT_DIR"

# Change to project directory
cd "$PROJECT_DIR"

echo ""
echo "🛑 Stopping all containers..."
docker-compose down --remove-orphans || true

echo ""
echo "🗑️  Removing containers..."
docker rm -f loan-postgres loan-kafka loan-zookeeper 2>/dev/null || true

echo ""
echo "💾 Removing volumes..."
# Remove named volumes
docker volume rm -f amf-loan-service_postgres_data 2>/dev/null || true
docker volume rm -f postgres_data 2>/dev/null || true

# Remove any dangling volumes
echo "🔍 Removing dangling volumes..."
docker volume prune -f

echo ""
echo "🌐 Removing networks..."
docker network rm -f amf-loan-service_loan-network 2>/dev/null || true
docker network rm -f loan-network 2>/dev/null || true

echo ""
echo "🧽 Cleaning up system..."
docker system prune -f

echo ""
echo "✅ Cleanup completed successfully!"
echo ""
echo "📋 Summary:"
echo "   - All containers stopped and removed"
echo "   - All volumes and data removed"
echo "   - Project-specific Docker images removed"
echo "   - Networks cleaned up"
echo ""
echo "💡 To restart the infrastructure, run:"
echo "   ./scripts/start-infra.sh"
echo ""
