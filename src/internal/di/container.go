package di

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/transaction"

	"go.uber.org/dig"
)

func mustProvide(c *dig.Container, constructor interface{}) {
	if err := c.Provide(constructor); err != nil {
		panic("DI provide: " + err.Error())
	}
}

func NewContainer(conf *config.Config) *dig.Container {
	c := dig.New()

	mustProvide(c, func() *config.Config { return conf })

	mustProvide(c, conf.ConnectDB)

	// Repository
	mustProvide(c, repository.NewUserRepository)
	mustProvide(c, repository.NewClientRepository)
	mustProvide(c, repository.NewFarmRepository)
	mustProvide(c, repository.NewMerchantRepository)
	mustProvide(c, repository.NewPondRepository)
	mustProvide(c, repository.NewActivePondRepository)
	mustProvide(c, repository.NewActivityRepository)
	mustProvide(c, repository.NewAdditionalCostRepository)
	mustProvide(c, repository.NewSellDetailRepository)
	mustProvide(c, repository.NewFishSizeGradeRepository)
	mustProvide(c, repository.NewWorkerRepository)
	mustProvide(c, repository.NewFarmGroupRepository)
	mustProvide(c, repository.NewFeedCollectionRepository)
	mustProvide(c, repository.NewFeedPriceHistoryRepository)
	mustProvide(c, repository.NewDailyLogRepository)

	// Transaction
	mustProvide(c, transaction.NewManager)

	// Service
	mustProvide(c, service.NewUserService)
	mustProvide(c, service.NewAuthService)
	mustProvide(c, service.NewClientService)
	mustProvide(c, service.NewFarmService)
	mustProvide(c, service.NewMerchantService)
	mustProvide(c, service.NewPondService)
	mustProvide(c, service.NewWorkerService)
	mustProvide(c, service.NewFarmGroupService)
	mustProvide(c, service.NewFeedCollectionService)
	mustProvide(c, service.NewFeedPriceHistoryService)
	mustProvide(c, service.NewFishSizeGradeService)
	mustProvide(c, service.NewDailyLogService)

	// Handler
	mustProvide(c, handler.NewUserHandler)
	mustProvide(c, handler.NewAuthHandler)
	mustProvide(c, handler.NewClientHandler)
	mustProvide(c, handler.NewFarmHandler)
	mustProvide(c, handler.NewMerchantHandler)
	mustProvide(c, handler.NewPondHandler)
	mustProvide(c, handler.NewWorkerHandler)
	mustProvide(c, handler.NewFarmGroupHandler)
	mustProvide(c, handler.NewFeedCollectionHandler)
	mustProvide(c, handler.NewFeedPriceHistoryHandler)
	mustProvide(c, handler.NewFishSizeGradeHandler)
	mustProvide(c, handler.NewDailyLogHandler)
	mustProvide(c, handler.NewHandler)

	return c
}
