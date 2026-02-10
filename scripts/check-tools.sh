#!/usr/bin/env bash
#
# Prerequisites check for go-testing-tools skill.
# Verifies: Go (>= 1.21), bubbletea dependency, git.
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

ERRORS=0

check() {
    local name="$1"
    local cmd="$2"
    local install_hint="$3"

    if eval "$cmd" &>/dev/null; then
        echo -e "${GREEN}OK${NC}  $name"
        return 0
    else
        echo -e "${RED}FAIL${NC}  $name"
        echo -e "  ${YELLOW}-> $install_hint${NC}"
        ERRORS=$((ERRORS + 1))
        return 1
    fi
}

echo "Checking prerequisites for go-testing-tools..."
echo ""

# --- Git ---
check "git" \
    "command -v git" \
    "Install git: https://git-scm.com/downloads"

# --- Go installed ---
check "Go installed" \
    "command -v go" \
    "Install Go: brew install go (or https://go.dev/dl/)"

# --- Go version >= 1.21 ---
if command -v go &>/dev/null; then
    GO_VERSION=$(go version | grep -oE '[0-9]+\.[0-9]+' | head -1)
    GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
    GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)

    if [[ "$GO_MAJOR" -gt 1 ]] || { [[ "$GO_MAJOR" -eq 1 ]] && [[ "$GO_MINOR" -ge 21 ]]; }; then
        echo -e "${GREEN}OK${NC}  Go version >= 1.21 (found: $GO_VERSION)"
    else
        echo -e "${RED}FAIL${NC}  Go version >= 1.21 (found: $GO_VERSION)"
        echo -e "  ${YELLOW}-> Update Go: brew upgrade go${NC}"
        ERRORS=$((ERRORS + 1))
    fi
fi

# --- bubbletea dependency ---
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
GOMOD="$REPO_DIR/tuitestkit/go.mod"

if [[ -f "$GOMOD" ]]; then
    if grep -q 'charmbracelet/bubbletea' "$GOMOD"; then
        echo -e "${GREEN}OK${NC}  bubbletea in go.mod"
    else
        echo -e "${RED}FAIL${NC}  bubbletea not found in go.mod"
        echo -e "  ${YELLOW}-> Run: cd tuitestkit && go get github.com/charmbracelet/bubbletea${NC}"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo -e "${YELLOW}WARN${NC}  go.mod not found at $GOMOD (tuitestkit not yet created?)"
fi

# --- Summary ---
echo ""
if [[ $ERRORS -eq 0 ]]; then
    echo -e "${GREEN}All prerequisites satisfied.${NC}"
    exit 0
else
    echo -e "${RED}Missing $ERRORS prerequisite(s). Fix them and re-run.${NC}"
    exit 1
fi
