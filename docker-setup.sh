#!/bin/bash

# Timeful Docker Quick Start Script
# This script helps you set up Timeful with Docker Compose

set -e

echo "======================================"
echo "Timeful Docker Setup"
echo "======================================"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker is not installed."
    echo "Please install Docker from: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is available
if ! docker compose version &> /dev/null; then
    echo "❌ Error: Docker Compose is not available."
    echo "Please install Docker Compose from: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created"
    echo ""
    echo "⚠️  IMPORTANT: You need to configure your .env file before starting the application."
    echo ""
    echo "Required configuration:"
    echo "1. Google OAuth credentials (CLIENT_ID and CLIENT_SECRET)"
    echo "   - Get them from: https://console.cloud.google.com/"
    echo ""
    echo "2. Encryption key (ENCRYPTION_KEY)"
    echo "   - Generate with: openssl rand -base64 32"
    echo ""
    read -p "Press Enter to edit .env file now, or Ctrl+C to exit and edit it manually..."
    
    # Try to open in an editor
    if command -v nano &> /dev/null; then
        nano .env
    elif command -v vim &> /dev/null; then
        vim .env
    elif command -v vi &> /dev/null; then
        vi .env
    else
        echo "Please edit .env file manually with your preferred editor."
        exit 0
    fi
else
    echo "✅ .env file already exists"
    echo ""
fi

# Ask user if they want to start the application
echo "======================================"
echo "Ready to start Timeful!"
echo "======================================"
echo ""
read -p "Do you want to start the application now? (y/n) " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🚀 Starting Timeful with Docker Compose..."
    echo ""
    
    docker compose up -d
    
    echo ""
    echo "======================================"
    echo "✅ Timeful is starting!"
    echo "======================================"
    echo ""
    echo "Access the application at: http://localhost:3002"
    echo ""
    echo "Useful commands:"
    echo "  - View logs: docker compose logs -f"
    echo "  - Stop: docker compose down"
    echo "  - Restart: docker compose restart"
    echo ""
    echo "For more information, see DOCKER.md"
else
    echo ""
    echo "You can start the application later with:"
    echo "  docker compose up -d"
    echo ""
fi
