# =====================================================
# Etapa 1: Build do CSS com Tailwind usando Node
# =====================================================
FROM node:alpine AS tailwind-builder
WORKDIR /app

# Copia o arquivo package.json (certifique-se que o arquivo existe em ./package.json)
COPY package.json .

# Instala as dependências do projeto
RUN npm install -g tailwindcss

# Copia o arquivo CSS de entrada (certifique-se que o arquivo existe em ./static/css/input.css)
COPY static/css/input.css ./static/css/input.css

# (Opcional) Se necessário, gere o arquivo de configuração do Tailwind
# RUN npx tailwindcss init

# Executa o build do CSS: gera o arquivo minificado style.min.css
RUN npm run build:tailwind

# =====================================================
# Etapa 2: Build da aplicação Go
# =====================================================
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Instala dependências necessárias para compilar a aplicação
RUN apk add --no-cache make gcc g++ curl

# Instala ferramentas auxiliares (como templ e air)
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest

# Copia o CSS gerado na etapa anterior para o diretório correspondente
COPY --from=tailwind-builder /app/static/css/style.min.css ./static/css/

# Copia o restante dos arquivos do projeto para a build do Go
COPY . .

# (Opcional) Gera os templates se sua aplicação os utilizar
RUN templ generate

# Compila a aplicação Go com a flag que define o ambiente como produção
RUN go build -ldflags "-X main.Environment=production" -o ./bin/app ./cmd/main.go


# =====================================================
# Etapa Final: Imagem mínima para execução
# =====================================================
FROM alpine:latest
WORKDIR /app

# Instala certificados CA para conexões HTTPS (se necessário)
RUN apk add --no-cache ca-certificates

# Copia apenas os artefatos necessários da etapa de build
COPY --from=builder /app/bin/app .
COPY --from=builder /app/static ./static

# Expõe a porta utilizada pela aplicação
EXPOSE 4000

# Comando para iniciar a aplicação
CMD ["./app"]
