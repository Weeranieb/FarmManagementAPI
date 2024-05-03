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
mockClientRepo:
	mockery --dir=pkg/repositories --name=IClientRepository --filename=IClientRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockUserRepo:
	mockery --dir=pkg/repositories --name=IUserRepository --filename=IUserRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockActivePondRepo:
	mockery --dir=pkg/repositories --name=IActivePondRepository --filename=IActivePondRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockActivityRepo:
	mockery --dir=pkg/repositories --name=IActivityRepository --filename=IActivityRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockFarmGroupRepo:
	mockery --dir=pkg/repositories --name=IFarmGroupRepository --filename=IFarmGroupRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockFarmRepo:
	mockery --dir=pkg/repositories --name=IFarmRepository --filename=IFarmRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockFarmOnFarmGroupRepositoryRepo:
	mockery --dir=pkg/repositories --name=IFarmOnFarmGroupRepository --filename=IFarmOnFarmGroupRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockPondRepo:
	mockery --dir=pkg/repositories --name=IPondRepository --filename=IPondRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockSellDetailRepo:
	mockery --dir=pkg/repositories --name=ISellDetailRepository --filename=ISellDetailRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockBillRepo:
	mockery --dir=pkg/repositories --name=IBillRepository --filename=IBillRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockDailyFeedRepo:
	mockery --dir=pkg/repositories --name=IDailyFeedRepository --filename=IDailyFeedRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockFeedCollectionRepo:
	mockery --dir=pkg/repositories --name=IFeedCollectionRepository --filename=IFeedCollectionRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockFeedPriceHistoryRepo:
	mockery --dir=pkg/repositories --name=IFeedPriceHistoryRepository --filename=IFeedPriceHistoryRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockMerchantRepo:
	mockery --dir=pkg/repositories --name=IMerchantRepository --filename=IMerchantRepository.go --output=pkg/repositories/mocks  --outpkg=mocks
mockWorkerRepo:
	mockery --dir=pkg/repositories --name=IWorkerRepository --filename=IWorkerRepository.go --output=pkg/repositories/mocks  --outpkg=mocks

.PHONY: postgres createdb dropdb  migrateup migratedown sqlc test server mock migrateup1 migratedown2 mockClientRepo mockActivePondRepo mockUserRepo mockActivityRepo mockFarmGroupRepo mockFarmRepo mockFarmOnFarmGroupRepositoryRepo mockPondRepo mockSellDetailRepo mockBillRepo