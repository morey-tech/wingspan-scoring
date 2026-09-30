#!/bin/bash
set -e

echo "Running post-create setup..."

# Install Go dependencies
echo "Downloading Go dependencies..."
go mod download
echo "Go dependencies installed"

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "Installing GitHub CLI..."
    curl -fsSL https://cli.github.com/packages/rpm/gh-cli.repo | sudo tee /etc/yum.repos.d/gh-cli.repo
    sudo dnf install -y gh
    echo "GitHub CLI installed: $(gh --version)"
else
    echo "GitHub CLI already available: $(gh --version)"
fi

echo "Post-create setup complete!"
