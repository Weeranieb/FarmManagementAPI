postgres:
	docker run --name FarmManagement -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it FarmManagement createdb --username=root --owner=root FarmManagementDB
dropdb:
	docker exec -it FarmManagement dropdb FarmManagementDB
migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/FarmManagementDB?sslmode=disable" -verbose up
migrateup1:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/FarmManagementDB?sslmode=disable" -verbose up 1
migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/FarmManagementDB?sslmode=disable" -verbose down
migratedown1:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/FarmManagementDB?sslmode=disable" -verbose down 1	
test:
	go test -v -cover ./...	
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Weeranieb/simplebank/db/sqlc Store 

.PHONY: postgres createdb dropdb  migrateup migratedown sqlc test server mock migrateup1 migratedown2