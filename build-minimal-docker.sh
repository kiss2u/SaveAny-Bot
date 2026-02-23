#!/bin/bash
# Build script for SaveAny-Bot minimal Docker image

# This script builds the most minimal version of SaveAny-Bot using the pico Dockerfile
# To use this script, ensure you have Docker installed and running

echo "Building SaveAny-Bot minimal Docker image..."
echo "Using Dockerfile.pico (most minimal configuration)"

# Build the image with pico tags
docker build -f dist/Dockerfile.pico -t saveany-bot:pico . --no-cache

echo ""
echo "Built image saveany-bot:pico with minimal features:"
echo "- No JS parser support"
echo "- No MinIO support" 
echo "- Minimal SQLite driver"
echo "- No bubbletea TUI components"
echo ""
echo "To run the container, use the following example command:"
echo "docker run -d --name saveany-bot -v /path/to/config.toml:/app/config.toml saveany-bot:pico"

