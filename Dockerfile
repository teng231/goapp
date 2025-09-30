# Start from the latest golang base image
FROM public.ecr.aws/docker/library/golang:1.24-alpine3.21
RUN apk --no-cache add ca-certificates tzdata curl nano

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o tikcart .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./tikcart", "start"]