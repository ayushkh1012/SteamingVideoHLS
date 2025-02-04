#!/bin/bash

# Build Docker image
docker build -t video-streaming:latest .

# Create Kind cluster if it doesn't exist
if ! kind get clusters | grep -q "^kind$"; then
  kind create cluster --config kind-config.yaml
fi

# Load image into Kind cluster
kind load docker-image video-streaming:latest

# Install/Upgrade Helm chart
helm upgrade --install video-streaming ./helm/video-streaming \
  -f ./helm/video-streaming/values-local.yaml \
  --create-namespace \
  --namespace video-streaming

# Wait for deployment to be ready
kubectl wait --namespace video-streaming \
  --for=condition=available \
  --timeout=300s \
  deployment/video-streaming

echo "Video streaming service is running!"
echo "Access it at http://localhost:8080" 