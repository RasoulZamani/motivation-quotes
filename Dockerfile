# Stage 1: Builder (with CGO)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy only module files first (for layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Install SQLite build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /motivational-api

# Stage 2: Production
FROM alpine:3.19

# Install SQLite runtime dependencies
RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY --from=builder /motivational-api /app/
COPY --from=builder /app/storage/quotes.json /app/storage/

VOLUME /app/storage
EXPOSE 8080
CMD ["/app/motivational-api"]