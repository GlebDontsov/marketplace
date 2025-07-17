FROM golang:1.23-alpine AS builder

ENV GOPROXY=https://goproxy.io,direct
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app
COPY . .

RUN swag init -g ./main.go
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o marketplace ./main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/marketplace .
COPY --from=builder /app/.env* ./

EXPOSE 8080

CMD ["./marketplace"]