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
    curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
    chmod +x tailwindcss-linux-x64
    tailwindcss-linux-x64 -i ./static/css/input.css -o ./static/css/style.min.css --minify
    templ generate
    go build -ldflags "-X main.Environment=production" -o ./bin/app ./cmd/main.go
    ./bin/app
fi
