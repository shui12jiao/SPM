# 构建
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app main.go

# 运行
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY db/migration ./db/migration

EXPOSE 7077
CMD ["./app"]