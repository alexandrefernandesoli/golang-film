# Stage 1: Build the frontend assets using Node.js
FROM node:lts AS frontend-builder

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json (if exists) for npm dependencies
COPY package*.json ./

# Install npm dependencies globally (TailwindCSS)
RUN npm install -g tailwindcss

# Install TailwindCSS globally
RUN npm install -g tailwindcss postcss autoprefixer

# Copy the rest of the frontend source files
COPY . .

# Build TailwindCSS assets
CMD ["npx", "tailwindcss", "-i", "./static/css/input.css", "-o", "./static/css/style.min.css", "--minify"]

# Stage 2: Build the Go application
FROM golang:1.23 AS go-builder

# Set the working directory inside the container
WORKDIR /app

# Install the 'templ' tool
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy Go module files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod tidy

# Copy the rest of the Go source files
COPY . .

# Generate templ files (if required)
RUN templ generate

# Build the Go application with production flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Environment=production" \
    -o ./bin/app ./cmd/main.go

# Stage 3: Create the final runtime image
FROM alpine:latest AS runtime

# Install necessary runtime dependencies (e.g., ca-certificates for HTTPS)
RUN apk add --no-cache bash

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the Go builder stage
COPY --from=go-builder /app/bin/app /app/bin/app

# Copy the built frontend assets from the frontend builder stage
COPY --from=frontend-builder /app/dist /app/dist

# Expose the port the app will run on (adjust as needed)
EXPOSE 4000

# Command to run the application
CMD ["/app/bin/app"]