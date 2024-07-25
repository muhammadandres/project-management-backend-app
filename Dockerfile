FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
RUN chmod +x main

EXPOSE 4040

CMD ["./main"]