#!/bin/bash

# Prompt user to enter GitHub access token
read -p "Enter GitHub access token: " GITHUB_TOKEN

# Prompt user to enter GitHub username
read -p "Enter GitHub username: " GITHUB_USERNAME

# Set environment variables
export GITHUB_TOKEN
export GITHUB_USERNAME

# Run the automation script
go run main.go
