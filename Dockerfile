# Builder stage
FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache curl
# Create working directories
RUN mkdir /app/bin -p
# Set home directory
WORKDIR /app
# Copy go.mod
COPY go.mod go.sum /app/
# Download go dependencies
RUN go mod download
# Copy all local files
COPY . /app
# Build the Go application
RUN GOOS=linux go build -o /app/bin/app ./cmd/app




# Start a new stage using the Alpine Linux image
FROM alpine:latest as dev

# Install packages
RUN apk --no-cache add ca-certificates && apk add --no-cache bash
# Create home directory
WORKDIR /app
# Copy the built executable from the builder stage
COPY --from=builder /app/bin/app /app/app
# Print the contents of /app after the COPY command
RUN ls
# Copy config file
COPY /config/local.yaml ./local.yaml
EXPOSE 8080
# Define the command to run the application
CMD ["./app","--config=./local.yaml"]