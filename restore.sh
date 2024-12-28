#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
until pg_isready -h localhost; do
  echo "Waiting for PostgreSQL to start..."
  sleep 1
done

# Restore the database dump
pg_restore -U postgres -d carbond /docker-entrypoint-initdb.d/carbond_backup.dump