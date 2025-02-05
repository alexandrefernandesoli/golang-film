#!/bin/sh
set -e  # Exit on error

cd /app

if [ "$ENVIRONMENT" = "development" ]; then
    echo "Starting in development mode..."

    echo "Starting dev server..."
    air &

    echo "Starting templ watch..."
    make templ-watch &

    wait
else
    echo "Starting in production mode..."
    tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify
    templ generate
    go build -ldflags "-X main.Environment=production" -o ./bin/app ./cmd/main.go
    ./bin/app
fi
