#!/bin/bash

# BrewPy Release Script
# Usage: ./scripts/release.sh 1.0.1

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 1.0.1"
    exit 1
fi

VERSION=$1
TAG="v$VERSION"

echo "Starting release process for $TAG"

# Check if we're in the right directory
if [ ! -f "src/main.go" ]; then
    echo "Error: Please run this script from the project root directory"
    exit 1
fi

# Check if git is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: Git working directory is not clean. Please commit or stash changes."
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^$TAG$"; then
    echo "Tag $TAG already exists"
    exit 1
fi

# Build release artifacts
echo "Building release artifacts..."
make release VERSION=$VERSION

# Check if release artifacts were created
if [ ! -f "dist/brewpy-v$VERSION-darwin-arm64.tar.gz" ]; then
    echo "Error: Release artifacts not found"
    exit 1
fi

# Create and push git tag
echo "Creating git tag $TAG..."
git tag -a $TAG -m "Release $TAG"
git push origin $TAG

# Create GitHub release
echo "Creating GitHub release..."
if command -v gh &> /dev/null; then
    gh release create $TAG \
        --title "BrewPy $TAG" \
        --notes "## BrewPy $TAG

### Installation
\`\`\`bash
brew install landoncrabtree/tap/brewpy
\`\`\`" \
        dist/brewpy-v$VERSION-darwin-arm64.tar.gz \
        dist/brewpy-v$VERSION-darwin-amd64.tar.gz \
        dist/checksums.txt
    
    echo "GitHub release created successfully!"
    echo "GitHub Actions will automatically update the Homebrew formula"
else
    echo "GitHub CLI not found. Please create the release manually at:"
    echo "   https://github.com/landoncrabtree/brewpy/releases/new"
    echo "   Tag: $TAG"
    echo "   Upload files from: dist/"
fi

echo ""
echo "Release $TAG completed successfully!"
echo ""
echo "Next steps:"
echo "1. GitHub Actions will automatically update the Homebrew formula"
echo "2. Users can install with: brew install landoncrabtree/tap/brewpy"
echo "3. Monitor the GitHub Actions workflow for any issues" 
