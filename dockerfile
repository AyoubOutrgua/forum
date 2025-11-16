# Build stage
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY handlers ./handlers
COPY helpers ./helpers
COPY middleware ./middleware
COPY routes ./routing
COPY tools ./tools
COPY database ./database

RUN go build -o main ./cmd/main.go


# Final stage
FROM ubuntu:24.04
WORKDIR /app

COPY --from=builder /app/main ./main
COPY templates/ ./templates/
COPY static/ ./static/
COPY db/ ./db/
# Metadata
LABEL maintainers="mhilli, boulhaj, melghama, aoutrgua, btoumana"
LABEL version="1.0"
LABEL description="forum"
# Expose port 8080
EXPOSE 8080
# Command to run the application
CMD ["./main"]