package di

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/transaction"

	"go.uber.org/dig"
)

func NewContainer(conf *config.Config) *dig.Container {
	c := dig.New()

	c.Provide(func() *config.Config { return conf })

	c.Provide(conf.ConnectDB)

	// Repository
	c.Provide(repository.NewUserRepository)
	c.Provide(repository.NewClientRepository)
	c.Provide(repository.NewFarmRepository)
	c.Provide(repository.NewMerchantRepository)
	c.Provide(repository.NewPondRepository)
	c.Provide(repository.NewActivePondRepository)
	c.Provide(repository.NewActivityRepository)
	c.Provide(repository.NewAdditionalCostRepository)
	c.Provide(repository.NewSellDetailRepository)
	c.Provide(repository.NewFishSizeGradeRepository)
	c.Provide(repository.NewWorkerRepository)
	c.Provide(repository.NewFarmGroupRepository)
	c.Provide(repository.NewFeedCollectionRepository)
	c.Provide(repository.NewFeedPriceHistoryRepository)
	c.Provide(repository.NewDailyLogRepository)

	// Transaction
	c.Provide(transaction.NewManager)

	// Service
	c.Provide(service.NewUserService)
	c.Provide(service.NewAuthService)
	c.Provide(service.NewClientService)
	c.Provide(service.NewFarmService)
	c.Provide(service.NewMerchantService)
	c.Provide(service.NewPondService)
	c.Provide(service.NewWorkerService)
	c.Provide(service.NewFarmGroupService)
	c.Provide(service.NewFeedCollectionService)
	c.Provide(service.NewFeedPriceHistoryService)
	c.Provide(service.NewFishSizeGradeService)
	c.Provide(service.NewDailyLogService)

	// Handler
	c.Provide(handler.NewUserHandler)
	c.Provide(handler.NewAuthHandler)
	c.Provide(handler.NewClientHandler)
	c.Provide(handler.NewFarmHandler)
	c.Provide(handler.NewMerchantHandler)
	c.Provide(handler.NewPondHandler)
	c.Provide(handler.NewWorkerHandler)
	c.Provide(handler.NewFarmGroupHandler)
	c.Provide(handler.NewFeedCollectionHandler)
	c.Provide(handler.NewFeedPriceHistoryHandler)
	c.Provide(handler.NewFishSizeGradeHandler)
	c.Provide(handler.NewDailyLogHandler)
	c.Provide(handler.NewHandler)

	return c
}
