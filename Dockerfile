# Etapa de construção
FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o weather-service .

# Etapa final
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/weather-service .
ENTRYPOINT ["./weather-service"]
