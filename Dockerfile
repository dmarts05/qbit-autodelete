# Build stage
FROM golang:1.24-alpine AS builder

# Install git (needed for go modules)
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with flags to reduce size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o qbit-autodelete ./cmd

# Final stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/qbit-autodelete .

# Default command
CMD ["./qbit-autodelete"]