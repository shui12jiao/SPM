.PHONY: postgres createdb dropdb sqlc swag

postgres:
	docker run --name postgres --network host -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:17-alpine

createdb:
	docker exec -it postgres createdb --username=admin --owner=admin man
	
dropdb:
	docker exec -it postgres dropdb --username=admin man

migrate:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

swag:
	swag init -g main.go --parseDependency --parseInternal --parseDepth 1
	@echo "Swagger documentation generated in the docs directory."