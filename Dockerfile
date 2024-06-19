FROM golang:latest as builder
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY ./ ./

# Build
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main .

FROM alpine:latest


RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY swagger.yaml /app/swagger.yaml

# Expose the port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
