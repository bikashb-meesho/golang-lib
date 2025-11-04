#!/bin/bash

# Pre-commit validation script for Library
# Runs all checks that CI/CD will run

set -e  # Exit on first error

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}🔍 Running Library Pre-Commit Checks${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

FAILED=0

run_check() {
    local check_name=$1
    local check_cmd=$2
    
    echo -e "${YELLOW}▶ ${check_name}...${NC}"
    if eval "$check_cmd"; then
        echo -e "${GREEN}✅ ${check_name} passed${NC}"
        echo ""
        return 0
    else
        echo -e "${RED}❌ ${check_name} failed${NC}"
        echo ""
        FAILED=1
        return 1
    fi
}

# 1. Go Format Check
run_check "Go Format Check" '
    UNFORMATTED=$(gofmt -l .)
    if [ -n "$UNFORMATTED" ]; then
        echo "The following files are not formatted:"
        echo "$UNFORMATTED"
        echo ""
        echo "Run: go fmt ./..."
        exit 1
    fi
'

# 2. Go Vet
run_check "Go Vet" "go vet ./..."

# 2b. golangci-lint (if installed)
if command -v golangci-lint &> /dev/null; then
    run_check "golangci-lint" "golangci-lint run ./..."
else
    echo -e "${YELLOW}⚠️  golangci-lint not installed (skipping)${NC}"
    echo "   Install: brew install golangci-lint"
    echo ""
fi

# 3. Go Mod Verify
run_check "Go Mod Verify" "go mod verify"

# 4. Go Mod Tidy Check
run_check "Go Mod Tidy Check" '
    cp go.mod go.mod.backup
    cp go.sum go.sum.backup
    go mod tidy
    if ! diff -q go.mod go.mod.backup >/dev/null 2>&1 || ! diff -q go.sum go.sum.backup >/dev/null 2>&1; then
        echo "go.mod or go.sum is not tidy"
        mv go.mod.backup go.mod
        mv go.sum.backup go.sum
        exit 1
    fi
    rm go.mod.backup go.sum.backup
'

# 5. Unit Tests
run_check "Unit Tests" "go test -v ./..."

# 6. Race Detection
run_check "Race Detection" "go test -race ./..."

# 7. Coverage Check
run_check "Coverage Check" '
    go test -coverprofile=coverage.out ./... >/dev/null
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk "{print \$3}" | sed "s/%//")
    echo "Total coverage: ${COVERAGE}%"
    
    THRESHOLD=20  # Adjusted for sample project
    if [ $(echo "$COVERAGE < $THRESHOLD" | bc 2>/dev/null || echo 0) -eq 1 ]; then
        echo "Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
        rm coverage.out
        exit 1
    fi
    rm coverage.out
'

# 8. Build Verification
run_check "Build Verification" "go build ./..."

# Summary
echo ""
echo -e "${BLUE}================================================${NC}"
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All library checks passed!${NC}"
    echo -e "${BLUE}================================================${NC}"
    exit 0
else
    echo -e "${RED}❌ Some checks failed.${NC}"
    echo -e "${BLUE}================================================${NC}"
    exit 1
fi

