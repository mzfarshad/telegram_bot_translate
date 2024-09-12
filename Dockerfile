# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

ENV CGO_ENABLED=1

RUN go mod tidy
RUN go build -o translate-bot ./cmd

# Stage 2: Final
FROM alpine:3.19
LABEL MAINTAINER=" Author <farshad.akbari.arzati@gmail.com>"

WORKDIR /app

COPY --from=builder /app/translate-bot /app/translate-bot
COPY .env /app/.env
# Copy the specific directories containing .txt files
COPY internal/key /app/internal/key
COPY internal/help /app/internal/help

CMD ["./translate-bot"]
