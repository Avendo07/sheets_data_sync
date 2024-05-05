# Use the official Golang image as the base
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Download Go modules (if using modules)
#RUN go mod download

# Copy your Go source code into the container
COPY . .

# Build the Go binary (replace "your-script-name" with your actual script name)
RUN go build -o script

# Expose the port the script might use (optional, modify if needed)
EXPOSE 8080

# Run the built Go binary as the entrypoint
CMD ["./script"]