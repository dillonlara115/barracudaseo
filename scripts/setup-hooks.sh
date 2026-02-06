
#!/bin/bash

# Setup git hooks for Barracuda SEO development
# Run this script after cloning the repository

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Setting up git hooks..."

# Create hooks directory if it doesn't exist
mkdir -p "$PROJECT_ROOT/.git/hooks"

# Create pre-commit hook
cat > "$PROJECT_ROOT/.git/hooks/pre-commit" << 'EOF'
#!/bin/bash

# Pre-commit hook for Barracuda SEO
# Runs formatting and type checking before each commit

set -e

echo "Running pre-commit checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we have staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

# Check if we have staged frontend files
STAGED_FRONTEND_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep -E '^web/.*\.(js|svelte|css)$' || true)

# Go formatting
if [ -n "$STAGED_GO_FILES" ]; then
    echo -e "${YELLOW}Checking Go formatting...${NC}"

    # Run go fmt on staged files
    for file in $STAGED_GO_FILES; do
        if [ -f "$file" ]; then
            gofmt -l -w "$file"
            git add "$file"
        fi
    done

    echo -e "${GREEN}Go formatting complete${NC}"
fi

# Frontend checks
if [ -n "$STAGED_FRONTEND_FILES" ]; then
    echo -e "${YELLOW}Checking frontend formatting...${NC}"

    # Check if node_modules exists
    if [ ! -d "web/node_modules" ]; then
        echo -e "${RED}Error: web/node_modules not found. Run 'cd web && npm install' first.${NC}"
        exit 1
    fi

    # Run Prettier on staged frontend files
    cd web
    for file in $STAGED_FRONTEND_FILES; do
        # Remove 'web/' prefix for the file path
        relative_file="${file#web/}"
        if [ -f "$relative_file" ]; then
            npx prettier --write "$relative_file" 2>/dev/null || true
            cd ..
            git add "$file"
            cd web
        fi
    done
    cd ..

    echo -e "${GREEN}Frontend formatting complete${NC}"

    # Run svelte-check (type checking)
    echo -e "${YELLOW}Running Svelte type checking...${NC}"
    cd web
    if npx svelte-check --tsconfig ./jsconfig.json 2>&1 | grep -E "^(Error|error)"; then
        echo -e "${RED}Svelte type check failed. Please fix errors before committing.${NC}"
        cd ..
        exit 1
    fi
    cd ..
    echo -e "${GREEN}Svelte type check passed${NC}"
fi

echo -e "${GREEN}All pre-commit checks passed!${NC}"
exit 0
EOF

# Make hook executable
chmod +x "$PROJECT_ROOT/.git/hooks/pre-commit"

echo "Git hooks installed successfully!"
echo ""
echo "The pre-commit hook will:"
echo "  - Format Go files with gofmt"
echo "  - Format frontend files with Prettier"
echo "  - Run Svelte type checking"
echo ""
echo "Make sure to run 'cd web && npm install' to install frontend dependencies."
