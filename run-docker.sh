#!/bin/bash

# Build the Docker image
docker build -t confluent-prometheus-sd .

# Run the container with environment variables from .env
docker run -p 8080:8080 --env-file .env confluent-prometheus-sd