#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -e

# Wait for the database to be ready
while ! nc -z $DB_HOST $DB_PORT; do 
  echo "Waiting for database at $DB_HOST:$DB_PORT..."
  sleep 1 
done

# Run database migrations
make migrate-up

# Start the application
make run
