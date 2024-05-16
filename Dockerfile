FROM golang:1.19 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY ./ ./

# Build
RUN CGO_ENABLED=0 go build -o /main

#Run stage
FROM alpine:latest

# Set destination for the binary
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /main /main

# Expose
EXPOSE 8000

# Run
CMD ["/main"]
