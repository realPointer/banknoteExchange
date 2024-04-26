#!/bin/bash

docker-compose up -d --build banknote_exchange

docker-compose up --build integration_tests
integration_tests_exit_code=$(docker inspect -f '{{ .State.ExitCode }}' integration_tests)

if [ $integration_tests_exit_code -ne 0 ]; then
  echo "Integration tests terminated with an error. Stopping and deleting containers..."
  make compose-down
  exit 1
fi

echo "Integration tests were successful. Service continues to run. Switching to logs..."
docker-compose logs -f