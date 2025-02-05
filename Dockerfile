FROM golang:1.23-alpine
WORKDIR /app

RUN apk add --no-cache make gcc g++ curl

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest

RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    chmod +x tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

COPY . .

RUN tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

RUN templ generate

RUN go build -ldflags "-X main.Environment=production" -o ./bin/app ./cmd/main.go

EXPOSE 4000

CMD ["./bin/app"]
