all: run

run:
	go run main.go

test:
	go test --cover ./...

start:
	docker compose up --build -d

stop:
	docker compose down
