# 构建阶段保持不变
FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ucaptcha

# 运行阶段优化
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/ucaptcha .
EXPOSE 8080
ENV GIN_MODE=release
CMD ["./ucaptcha"]