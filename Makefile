all: build run

build:
	go build -o bin/flight-itinerary ./cmd/main.go

run:
	go run ./cmd/main.go

deps:
	go mod download
	go mod tidy

docker-build:
	docker build -t flight-itinerary-go:latest .


docker-run:
	docker run -p 8080:8080 flight-itinerary-go:latest