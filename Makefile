postgres:
	docker run --name myfarm -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it myfarm createdb --username=root --owner=root boonma_farm
dropdb:
	docker exec -it myfarm dropdb boonma_farm
migrateup:
	migrate -path pkg/db/migration -database "postgresql://root:secret@localhost:5432/boonma_farm?sslmode=disable" -verbose up
migrateup1:
	migrate -path pkg/db/migration -database "postgresql://root:secret@localhost:5432/boonma_farm?sslmode=disable" -verbose up 1
migratedown:
	migrate -path pkg/db/migration -database "postgresql://root:secret@localhost:5432/boonma_farm?sslmode=disable" -verbose down
migratedown1:
	migrate -path pkg/db/migration -database "postgresql://root:secret@localhost:5432/boonma_farm?sslmode=disable" -verbose down 1	
sqlc:
	sqlc generate
test:
	go test -v -cover ./...	
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Weeranieb/simplebank/db/sqlc Store 

.PHONY: postgres createdb dropdb  migrateup migratedown sqlc test server mock migrateup1 migratedown2