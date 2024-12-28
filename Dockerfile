# Build stage
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . ./

# Change directory to cmd and build the Go application with static linking
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cmd/carbond-app .

# Final stage with PostgreSQL
FROM postgres:13 AS database
WORKDIR /db

# Environment variables for PostgreSQL
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=pwaiswa
ENV POSTGRES_DB=carbond

# Copy the database dump
COPY carbond_backup.dump /docker-entrypoint-initdb.d/carbond_backup.dump

# Add a script to restore the dump
COPY restore.sh /docker-entrypoint-initdb.d/restore.sh
RUN chmod +x /docker-entrypoint-initdb.d/restore.sh

# Add the app stage
FROM alpine:latest
WORKDIR /root/

# Install necessary dependencies
RUN apk --no-cache add ca-certificates bash

# Copy PostgreSQL data from the database image
COPY --from=database /db /root/db

# Copy the application binary
COPY --from=builder /app/cmd/carbond-app /root/carbond-app

# Copy the .env file
COPY /cmd/.env /root/.env

# Expose ports
EXPOSE 8080 5432

# Set the entrypoint script to run PostgreSQL and the app
CMD ["bash", "-c", "postgres -D /var/lib/postgresql/data & ./carbond-app"]
