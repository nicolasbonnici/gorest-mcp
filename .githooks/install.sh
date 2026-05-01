#!/bin/bash
#
# Install git hooks for gorest-mcp

set -e

HOOKS_DIR="$(cd "$(dirname "$0")" && pwd)"
GIT_DIR="$(git rev-parse --git-dir)"

echo "Installing git hooks..."

# Make hooks executable
chmod +x "$HOOKS_DIR/pre-commit"

# Create symlink
ln -sf "$HOOKS_DIR/pre-commit" "$GIT_DIR/hooks/pre-commit"

echo "✓ Git hooks installed successfully!"
echo ""
echo "Installed hooks:"
echo "  - pre-commit (format, vet, test)"
echo ""
echo "To skip hooks during commit, use: git commit --no-verify"
