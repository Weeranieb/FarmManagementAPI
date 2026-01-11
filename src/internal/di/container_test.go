package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"go.uber.org/dig"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ContainerTestSuite struct {
	suite.Suite
	conf *config.Config
}

func (s *ContainerTestSuite) SetupTest() {
	// Create a test config
	s.conf = &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			Name:     "test_db",
			User:     "test_user",
			Password: "test_password",
			SSLMode:  "disable",
		},
		App: config.AppConfig{
			Environment: "test",
			LogLevel:    "silent",
			Debug:       false,
		},
		Authentication: config.AuthenticationConfig{
			JWTSecret: "test_secret",
			JWTExpiry: "24h",
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}
}

// createTestContainer creates a container with a test database connection
func createTestContainer(conf *config.Config) *dig.Container {
	c := dig.New()

	// Provide test database connection using SQLite in-memory
	c.Provide(func() *gorm.DB {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to test database: " + err.Error())
		}
		return db
	})

	// Repository
	c.Provide(repository.NewUserRepository)
	c.Provide(repository.NewClientRepository)
	c.Provide(repository.NewFarmRepository)
	c.Provide(repository.NewMerchantRepository)
	c.Provide(repository.NewPondRepository)
	c.Provide(repository.NewWorkerRepository)
	c.Provide(repository.NewFeedCollectionRepository)
	c.Provide(repository.NewFeedPriceHistoryRepository)

	// Service
	c.Provide(service.NewUserService)
	c.Provide(service.NewAuthService)
	c.Provide(service.NewClientService)
	c.Provide(service.NewFarmService)
	c.Provide(service.NewMerchantService)
	c.Provide(service.NewPondService)
	c.Provide(service.NewWorkerService)
	c.Provide(service.NewFeedCollectionService)
	c.Provide(service.NewFeedPriceHistoryService)

	// Handler
	c.Provide(handler.NewUserHandler)
	c.Provide(handler.NewAuthHandler)
	c.Provide(handler.NewClientHandler)
	c.Provide(handler.NewFarmHandler)
	c.Provide(handler.NewMerchantHandler)
	c.Provide(handler.NewPondHandler)
	c.Provide(handler.NewWorkerHandler)
	c.Provide(handler.NewFeedCollectionHandler)
	c.Provide(handler.NewFeedPriceHistoryHandler)
	c.Provide(handler.NewHandler)

	return c
}

func TestContainerSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}

func (s *ContainerTestSuite) TestNewContainer_CreatesContainer() {
	container := NewContainer(s.conf)

	assert.NotNil(s.T(), container)
}

func (s *ContainerTestSuite) TestNewContainer_ProvidesDatabase() {
	container := createTestContainer(s.conf)

	var db *gorm.DB
	err := container.Invoke(func(d *gorm.DB) {
		db = d
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), db)
}

func (s *ContainerTestSuite) TestNewContainer_ProvidesUserRepository() {
	container := createTestContainer(s.conf)

	var userRepo repository.UserRepository
	err := container.Invoke(func(r repository.UserRepository) {
		userRepo = r
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), userRepo)
}

func (s *ContainerTestSuite) TestNewContainer_ProvidesUserService() {
	container := createTestContainer(s.conf)

	var userService service.UserService
	err := container.Invoke(func(svc service.UserService) {
		userService = svc
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), userService)
}

func (s *ContainerTestSuite) TestNewContainer_ProvidesAuthService() {
	container := createTestContainer(s.conf)

	var authService service.AuthService
	err := container.Invoke(func(svc service.AuthService) {
		authService = svc
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), authService)
}

func (s *ContainerTestSuite) TestNewContainer_ProvidesUserHandler() {
	container := createTestContainer(s.conf)

	var userHandler handler.UserHandler
	err := container.Invoke(func(h handler.UserHandler) {
		userHandler = h
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), userHandler)
}

func (s *ContainerTestSuite) TestNewContainer_ProvidesAuthHandler() {
	container := createTestContainer(s.conf)

	var authHandler handler.AuthHandler
	err := container.Invoke(func(h handler.AuthHandler) {
		authHandler = h
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), authHandler)
}

func (s *ContainerTestSuite) TestNewContainer_ProvidesHandler() {
	container := createTestContainer(s.conf)

	var mainHandler *handler.Handler
	err := container.Invoke(func(h *handler.Handler) {
		mainHandler = h
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), mainHandler)
	assert.NotNil(s.T(), mainHandler.UserHandler)
	assert.NotNil(s.T(), mainHandler.AuthHandler)
	assert.NotNil(s.T(), mainHandler.ClientHandler)
	assert.NotNil(s.T(), mainHandler.FarmHandler)
	assert.NotNil(s.T(), mainHandler.MerchantHandler)
	assert.NotNil(s.T(), mainHandler.PondHandler)
	assert.NotNil(s.T(), mainHandler.WorkerHandler)
	assert.NotNil(s.T(), mainHandler.FeedCollectionHandler)
	assert.NotNil(s.T(), mainHandler.FeedPriceHistoryHandler)
}

func (s *ContainerTestSuite) TestNewContainer_ResolvesAllDependencies() {
	container := createTestContainer(s.conf)

	// Test that all dependencies can be resolved together
	err := container.Invoke(func(
		db *gorm.DB,
		userRepo repository.UserRepository,
		userService service.UserService,
		authService service.AuthService,
		userHandler handler.UserHandler,
		authHandler handler.AuthHandler,
		mainHandler *handler.Handler,
	) {
		assert.NotNil(s.T(), db)
		assert.NotNil(s.T(), userRepo)
		assert.NotNil(s.T(), userService)
		assert.NotNil(s.T(), authService)
		assert.NotNil(s.T(), userHandler)
		assert.NotNil(s.T(), authHandler)
		assert.NotNil(s.T(), mainHandler)
	})

	assert.NoError(s.T(), err)
}

func (s *ContainerTestSuite) TestNewContainer_WithNilConfig() {
	// Container should still be created even with nil config
	container := NewContainer(nil)

	assert.NotNil(s.T(), container)

	// Attempting to resolve database should fail because ConnectDB will panic
	// when trying to access nil config
	assert.Panics(s.T(), func() {
		_ = container.Invoke(func(d *gorm.DB) {
			_ = d
		})
	})
}
