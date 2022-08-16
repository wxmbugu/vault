postgres:
	docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -v ~/postgres_data:/data/db -d postgres:14-alpine
createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres vault
startdb:
	docker start postgres
create_vault_table:
	migrate create -ext sql -dir . -seq create_vault
migrateup:
	migrate -path . -database "postgres://mpxaqnjcqewmve:849658d51852ea38573b12b5d2cb5973760507f4beb29638707a29071418771f@ec2-44-205-112-253.compute-1.amazonaws.com:5432/d58gq8lh7l5ru2" -verbose up
client:
	 npm run --prefix templates/frontend/ serve
.PHONY:postgres createdb create_vault_table migrateup client
