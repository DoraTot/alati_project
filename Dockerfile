FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY ./ ./

# Build
RUN CGO_ENABLED=0 go build -o /main

# Expose
EXPOSE 8080

# Run
CMD ["/main"]
