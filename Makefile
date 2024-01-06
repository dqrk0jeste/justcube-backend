DB_URL=postgresql://root:secret@localhost:5432/letscube?sslmode=disable

#dont know what it is
network:
	docker network create bank-network

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root letscube

dropdb:
	docker exec -it postgres dropdb letscube

terminaldb:
	docker exec -it postgres bin/sh

migrateup:
	migrate -path database/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path database/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path database/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path database/migrations -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir database/migrations -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis