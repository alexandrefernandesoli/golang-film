# =====================================================
# Etapa 1: Build do CSS com Tailwind (sem package.json)
# =====================================================
FROM alpine:latest AS tailwind-builder
WORKDIR /app

# Instala o curl para baixar o binário do Tailwind
RUN apk add --no-cache curl

# Copia os arquivos do projeto (garanta que o arquivo input.css esteja em ./static/css/)
COPY . .

# Baixa o binário do Tailwind, torna-o executável e move para o PATH
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    chmod +x tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# Executa o build do CSS com Tailwind
RUN tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify


# =====================================================
# Etapa 2: Build da aplicação Go
# =====================================================
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Instala dependências de sistema necessárias para compilar a aplicação
RUN apk add --no-cache make gcc g++ curl

# Instala ferramentas auxiliares (como templ e air)
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest

# Copia o CSS gerado na etapa anterior para este estágio
COPY --from=tailwind-builder /app/static/css/style.min.css ./static/css/

# Copia o restante dos arquivos do projeto
COPY . .

# (Opcional) Gera os templates, se sua aplicação utilizar o 'templ'
RUN templ generate

# Compila a aplicação Go
RUN go build -ldflags "-X main.Environment=production" -o ./bin/app ./cmd/main.go


# =====================================================
# Etapa Final: Imagem mínima para execução
# =====================================================
FROM alpine:latest
WORKDIR /app

# Instala certificados de CA para conexões seguras (se necessário)
RUN apk add --no-cache ca-certificates

# Copia o binário e os arquivos estáticos (incluindo o CSS gerado) da etapa builder
COPY --from=builder /app/bin/app .
COPY --from=builder /app/static ./static

# Expõe a porta utilizada pela aplicação
EXPOSE 4000

# Comando de inicialização
CMD ["./app"]
