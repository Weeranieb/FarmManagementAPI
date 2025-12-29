package di

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"

	"go.uber.org/dig"
)

func NewContainer(conf *config.Config) *dig.Container {
	c := dig.New()

	c.Provide(conf.ConnectDB)

	// Repository
	c.Provide(repository.NewUserRepository)
	c.Provide(repository.NewFarmRepository)
	c.Provide(repository.NewMerchantRepository)
	c.Provide(repository.NewPondRepository)
	c.Provide(repository.NewWorkerRepository)
	c.Provide(repository.NewFeedCollectionRepository)
	c.Provide(repository.NewFeedPriceHistoryRepository)

	// Service
	c.Provide(service.NewUserService)
	c.Provide(service.NewAuthService)
	c.Provide(service.NewFarmService)
	c.Provide(service.NewMerchantService)
	c.Provide(service.NewPondService)
	c.Provide(service.NewWorkerService)
	c.Provide(service.NewFeedCollectionService)
	c.Provide(service.NewFeedPriceHistoryService)

	// Handler
	c.Provide(handler.NewUserHandler)
	c.Provide(handler.NewAuthHandler)
	c.Provide(handler.NewFarmHandler)
	c.Provide(handler.NewMerchantHandler)
	c.Provide(handler.NewPondHandler)
	c.Provide(handler.NewWorkerHandler)
	c.Provide(handler.NewFeedCollectionHandler)
	c.Provide(handler.NewFeedPriceHistoryHandler)
	c.Provide(handler.NewHandler)

	return c
}
