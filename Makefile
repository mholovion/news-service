.PHONY: build run test docker-build docker-run

build:
	go build -o main ./cmd/main.go

run:
	go run ./cmd/main.go

docker-build:
	docker-compose build

docker-run:
	docker-compose up --build
