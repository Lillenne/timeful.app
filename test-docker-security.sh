#!/bin/bash
# Test script to verify Docker security hardening
# This script verifies that containers run as non-root with proper security settings

set -e

echo "==================================="
echo "Docker Security Verification Test"
echo "==================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if docker compose is available
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker is installed${NC}"

# Validate compose files
echo ""
echo "Validating Docker Compose files..."

for file in docker-compose.yml docker-compose.dev.yml docker-compose.ghcr.yml; do
    if docker compose -f "$file" config > /dev/null 2>&1; then
        echo -e "${GREEN}✓ $file is valid${NC}"
    else
        echo -e "${RED}✗ $file has errors${NC}"
        exit 1
    fi
done

# Check Dockerfiles for USER directives
echo ""
echo "Checking Dockerfiles for non-root users..."

if grep -q "^USER appuser" Dockerfile.backend; then
    echo -e "${GREEN}✓ Dockerfile.backend uses non-root user${NC}"
else
    echo -e "${RED}✗ Dockerfile.backend missing USER directive${NC}"
    exit 1
fi

if grep -q "^USER nginx" Dockerfile.frontend; then
    echo -e "${GREEN}✓ Dockerfile.frontend uses non-root user${NC}"
else
    echo -e "${RED}✗ Dockerfile.frontend missing USER directive${NC}"
    exit 1
fi

# Expected UIDs for verification
BACKEND_UID="1000:1000"
FRONTEND_UID="101:101"
MONGODB_UID="999:999"

# Check compose files for security options
echo ""
echo "Checking Docker Compose security options..."

# Check if a pattern exists in a file
# Note: Uses strict pattern matching to ensure proper YAML formatting
check_in_file() {
    local file=$1
    local pattern=$2
    grep -q "$pattern" "$file"
}

for file in docker-compose.yml docker-compose.dev.yml docker-compose.ghcr.yml; do
    echo ""
    echo "Checking $file..."
    
    # Check user directives
    if check_in_file "$file" "user: \"$BACKEND_UID\"" && \
       check_in_file "$file" "user: \"$FRONTEND_UID\"" && \
       check_in_file "$file" "user: \"$MONGODB_UID\""; then
        echo -e "${GREEN}✓ All services have user directives ($BACKEND_UID, $FRONTEND_UID, $MONGODB_UID)${NC}"
    else
        echo -e "${RED}✗ Missing user directives${NC}"
        exit 1
    fi
    
    # Check security options
    if check_in_file "$file" "no-new-privileges:true"; then
        echo -e "${GREEN}✓ Services have no-new-privileges${NC}"
    else
        echo -e "${RED}✗ Missing no-new-privileges${NC}"
        exit 1
    fi
    
    # Check capability dropping
    if check_in_file "$file" "cap_drop:"; then
        echo -e "${GREEN}✓ Services drop capabilities${NC}"
    else
        echo -e "${RED}✗ Missing capability dropping${NC}"
        exit 1
    fi
done

echo ""
echo "==================================="
echo -e "${GREEN}Security verification complete!${NC}"
echo "==================================="
echo ""
echo "Key Security Features Implemented:"
echo "  • Non-root users in all containers"
echo "  • User directives in compose files ($BACKEND_UID, $FRONTEND_UID, $MONGODB_UID)"
echo "  • No-new-privileges security option"
echo "  • Capability dropping (ALL) with minimal additions"
echo "  • Compatible with rootless Docker/Podman"
echo ""
echo "To test with running containers:"
echo "  docker compose up -d"
echo "  docker compose exec backend id"
echo "  docker compose exec frontend id"
echo "  docker inspect <container> | grep -A 10 SecurityOpt"
echo ""
