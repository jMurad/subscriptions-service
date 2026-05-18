include .env
export

APP_NAME=subscriptions-service

SERVICE=./internal/service/subscriptions
HANDLER=./internal/transport/http/handler/subscriptions
MIDDLEWARE=./internal/middleware
REPO=./internal/repository/postgres/subscriptions

.PHONY: \
	run \
	build \
	test \
	migrate-up \
	migrate-down \
	migrate-force \
	migrate-create \
	swagger \
	docker-up \
	docker-down \
	logs \
	test-service \
	test-handler \
	test-midlleware \
	test-repo \
	migrate-run \
	migrate-stop \
	clean

run:
	go run ./cmd/app

build:
	go build -o bin/$(APP_NAME) ./cmd/app

test:
	test-service
	test-handler
	test-middleware
	test-repo

test-service:
	go test -v $(SERVICE)

test-handler:
	go test -v $(HANDLER)

test-midlleware:
	go test -v $(MIDDLEWARE)

test-repo:
	go test -v $(REPO)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-force:
	@read -p "Version: " version; \
	migrate -path migrations -database "$(DATABASE_URL)" force $$version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

swagger:
	swag init -g main.go -d cmd/app,internal

docker-up:
	docker compose up --build

docker-down:
	docker compose down

postgres-up:
	docker compose up -d postgres

postgres-down:
	docker compose stop postgres

migrate-run:
	docker compose up -d migrate

migrate-stop:
	docker compose stop migrate	

docker-clean:
	docker compose down --rmi all --volumes	

logs:
	docker compose logs -f

clean:
	rm -rf bin

