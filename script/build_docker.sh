#!/bin/bash

# Script to build Docker images for Gateway and Quiz services

# Define service names and Dockerfile paths
GATEWAY_SERVICE="gateway"
QUIZ_SERVICE="quiz"

# Docker image names
GATEWAY_IMAGE="gateway-service:latest"
QUIZ_IMAGE="quiz-service:latest"

# Build the Gateway service Docker image
echo "Building Docker image for Gateway service..."
docker build -t $GATEWAY_IMAGE -f ./cmd/$GATEWAY_SERVICE/Dockerfile .
if [ $? -ne 0 ]; then
  echo "Failed to build Gateway service image."
  exit 1
fi

# Build the Quiz service Docker image
echo "Building Docker image for Quiz service..."
docker build -t $QUIZ_IMAGE -f ./cmd/$QUIZ_SERVICE/Dockerfile .
if [ $? -ne 0 ]; then
  echo "Failed to build Quiz service image."
  exit 1
fi

echo "Docker images built successfully!"
echo "Gateway Image: $GATEWAY_IMAGE"
echo "Quiz Image: $QUIZ_IMAGE"
