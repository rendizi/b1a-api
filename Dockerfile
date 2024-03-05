# Use an official Golang runtime as a parent image
FROM golang:1.22.0 AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Install necessary tools
RUN apt-get update && \
    apt-get install -y postgresql-client-12

# Create the postgres user and database
RUN su postgres -c "createuser -s postgres" && \
    su postgres -c "createdb -O postgres b1a" && \
    su postgres -c "psql -c \"alter user postgres with password '1\"\"' b1a"

# Build the Go binary
RUN go build -o main ./cmd/main.go

# Switch to a non-root user for security reasons
USER postgres

# Copy the binary to the final image and set the working directory
FROM postgres:latest
COPY --from=builder /app/main /usr/local/bin/

# Set the entrypoint to run the Go binary
ENTRYPOINT ["main"]
