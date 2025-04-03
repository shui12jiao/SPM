.PHONY: postgres createdb dropdb sqlc swag

postgres:
	docker run --name postgres --network host -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:17-alpine

createdb:
	docker exec -it postgres createdb --username=admin --owner=admin 
	
dropdb:
	docker exec -it postgres dropdb 

sqlc:
	sqlc generate

swag:
	swag init -g main.go --parseDependency --parseInternal --parseDepth 1
	@echo "Swagger documentation generated in the docs directory."