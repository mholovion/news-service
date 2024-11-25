FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o news_service ./cmd/main.go

EXPOSE 8080

CMD ["/app/news_service"]
