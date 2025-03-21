.PHONY: up down migrate swag

up:
	docker-compose up --build

down:
	docker-compose down -v

migrate:
	docker-compose run --rm app /app/main migrate

swag:
	swag init -g cmd/main/main.go -o docs