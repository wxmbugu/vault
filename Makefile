postgres:
	docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -v ~/postgres_data:/data/db -d postgres:14-alpine
createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres vault
startdb:
	docker start postgres
create_vault_table:
	migrate create -ext sql -dir . -seq create_vault
migrateup:
	migrate -path . -database "postgresql://postgres:secret@localhost:5432/vault?sslmode=disable" -verbose up
client:
	 npm run --prefix templates/frontend/ serve
.PHONY:postgres createdb create_vault_table migrateup client
