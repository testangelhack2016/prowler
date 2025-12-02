# Use an official Go runtime as a parent image
FROM golang:1.19-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Download dependencies and update go.sum
RUN go mod tidy

# Build the Go app
RUN go build -o /ingestor

# Expose port 8080 to the outside world
# EXPOSE 8080

# Command to run the executable
CMD ["/ingestor"]
