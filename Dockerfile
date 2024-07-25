FROM golang:1.22.2-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main ./main.go
RUN chmod +x main
EXPOSE 4040
CMD [ "./main" ]