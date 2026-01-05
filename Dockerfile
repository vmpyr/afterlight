# ==========================================
# Stage 1: Build Frontend
# ==========================================
FROM node:24-alpine AS web-builder
WORKDIR /build/web

# Install dependencies
COPY web/package*.json ./
RUN npm ci

# Copy source and build
COPY web/ ./
RUN npm run build

# ==========================================
# Stage 2: Build Backend (With Embedding)
# ==========================================
FROM golang:1.25-alpine AS api-builder
WORKDIR /build

# GCC needed for SQLite
RUN apk add --no-cache gcc musl-dev

# 1. Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy source code
COPY . .

# 3. Copy the 'dist' folder from Stage 1 into the expected location
COPY --from=web-builder /build/web/dist ./web/dist

# 4. Build the binary (Static link)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o afterlight main.go

# ==========================================
# Stage 3: Final Production Image
# ==========================================
FROM alpine:latest
WORKDIR /app

# Copy only the compiled binary
COPY --from=api-builder /build/afterlight .

# Create the data volume directory
RUN mkdir -p /data/artifacts

# Set standard paths for your app to read
ENV DB_PATH=/data/afterlight.db
ENV ARTIFACTS_PATH=/data/artifacts

# Expose HTTP port
EXPOSE 8080

# Run
CMD ["./afterlight"]
