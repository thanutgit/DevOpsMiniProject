FROM golang:1.26.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/main ./cmd

EXPOSE 3010
CMD ["./main"]