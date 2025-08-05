FROM golang:1.24-bullseye

WORKDIR /app

RUN apt-get update && apt-get install -y build-essential

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o pong cmd/pong/main.go

EXPOSE 8080

CMD ["./pong"]
