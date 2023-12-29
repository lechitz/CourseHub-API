# build stage
FROM golang:1.20 AS builder

# working directory
WORKDIR /app

# Copy only the necessary files to download dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

# Rebuild internal libraries and disable cgo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/coursehub-api ./cmd/coursehub-api

# final stage
FROM alpine:latest

# working directory
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=builder /app/coursehub-api .

# Expose the port on which the application will listen
EXPOSE 5001

# Run the coursehub-api command when the container starts
CMD ["./coursehub-api"]
