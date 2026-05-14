APP_NAME=subscriptions-service

DB_URL=postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable

.PHONY: run build test migrate-up migrate-down migrate-force \
docker-up docker-down logs swagger clean

run:
	go run ./cmd/app

build:
	go build -o bin/$(APP_NAME) ./cmd/app

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-force:
	@read -p "Version: " version; \
	migrate -path migrations -database "$(DB_URL)" force $$version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

swagger:
	swag init -g cmd/app/main.go

docker-up:
	docker compose up --build

docker-down:
	docker compose down

logs:
	docker compose logs -f

clean:
	rm -rf bin

swagger:
	swag init -g cmd/app/main.go
