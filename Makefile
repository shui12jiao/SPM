postgres:
	docker run --name postgres --network host -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:17-alpine

createdb:
	docker exec -it postgres createdb --username=admin --owner=admin 
	
dropdb:
	docker exec -it postgres dropdb 

sqlc:
	sqlc generate