# Dockerfile.build
FROM golang:1.23-alpine AS builder

# Install curl and git for the swag installation
RUN apk add --no-cache curl git
# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./

# Download necessary Go modules
RUN go mod download

# Copy the entire application code
COPY . .

# Generate Swagger documentation
RUN swag init

# Build the Go application binary
RUN go build -o user_db .

# Dockerfile.run
FROM alpine:3.20

WORKDIR /app

# Copy the binary and Swagger documentation
COPY --from=builder /app/user_db /app/
COPY --from=builder /app/docs /app/docs

EXPOSE 8080

CMD ["./user_db"]
