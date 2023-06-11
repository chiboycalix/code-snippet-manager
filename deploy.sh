#!/bin/bash

# Stop the existing Go server
systemctl stop code-snippet-manager

# Navigate to your Go server directory
cd /projects/code-snippet-manager

# Pull the latest changes from the GitHub repository
git pull

# Build the Go server application
go build

# Restart the Go server
systemctl start /projects/code-snippet-manager
