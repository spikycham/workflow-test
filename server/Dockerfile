# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build a statically linked binary for Alpine compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -o output .

# Stage 2: Final lightweight image
FROM alpine:3.20
# Install the Docker CLI quickly via apk
RUN apk add --no-cache docker-cli
WORKDIR /app
COPY --from=builder /app/output .

CMD ["./output"]
