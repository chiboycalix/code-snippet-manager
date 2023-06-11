#!/bin/bash

# Stop the existing Go server
docker stop code-snippet-manager

# Navigate to your Go server directory
cd "$(dirname "$0")"

# Pull the latest changes from the GitHub repository
git pull

# Build the Go server application
go build

# Restart the Go server
docker start code-snippet-manager
