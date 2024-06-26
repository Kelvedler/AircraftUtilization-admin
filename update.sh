#!/bin/sh

echo "Running database migrations"

migrate -path=/app/migration -database "$DATABASE_URL?sslmode=disable" up

if [ $? -eq 0 ]
then
  echo "Database migrations succeded"
  exit 0
else
  echo "Failed to run database migrations"
  exit 1
fi
