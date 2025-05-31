FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

EXPOSE 8080

# Run binary
CMD ["./main"]
