.PHONY: postgres createdb dropdb sqlc swag test migrate migrateup migratedown

DB_URL=postgres://admin:admin@192.168.31.128/man?sslmode=disable

postgres:
	docker run --name postgres --network host -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:17-alpine

createdb:
	docker exec -it postgres createdb --username=admin --owner=admin man
	
dropdb:
	docker exec -it postgres dropdb --username=admin man

migrate:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

mock:
	mockery

sqlc:
	sqlc generate

swag:
	swag init -g main.go --parseDependency --parseInternal --parseDepth 1
	@echo "Swagger documentation generated in the docs directory."

test:
	go test -v -cover -short ./...