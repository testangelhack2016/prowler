# Use an official Go runtime as a parent image
FROM golang:1.19-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o /engine

# Expose port 8081 to the outside world
# EXPOSE 8081

# Command to run the executable
CMD ["/engine"]
