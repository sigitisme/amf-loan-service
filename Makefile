build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

mock-users:
	go run ./cmd/create-mock-users/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t amf-loan-service .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

migrate-up:
	# Add migration tool commands here
	echo "Migrations will be handled by the application"

lint:
	golangci-lint run

deps:
	go mod download
	go mod tidy

.PHONY: build run mock-users test clean docker-build docker-run docker-stop migrate-up lint deps
