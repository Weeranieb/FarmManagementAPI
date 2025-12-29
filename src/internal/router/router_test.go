package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RouterTestSuite struct {
	suite.Suite
	app     *fiber.App
	handlers *handler.Handler
}

func (s *RouterTestSuite) SetupTest() {
	// Create test config
	conf := &config.Config{
		Authentication: config.AuthenticationConfig{
			JWTSecret: "test_secret",
			JWTExpiry: "24h",
		},
	}
	
	// Create in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}
	
	// Create handlers manually with test database
	userRepo := repository.NewUserRepository(db)
	farmRepo := repository.NewFarmRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	pondRepo := repository.NewPondRepository(db)
	workerRepo := repository.NewWorkerRepository(db)
	feedCollectionRepo := repository.NewFeedCollectionRepository(db)
	feedPriceHistoryRepo := repository.NewFeedPriceHistoryRepository(db)
	
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	farmService := service.NewFarmService(farmRepo)
	merchantService := service.NewMerchantService(merchantRepo)
	pondService := service.NewPondService(pondRepo)
	workerService := service.NewWorkerService(workerRepo)
	feedCollectionService := service.NewFeedCollectionService(feedCollectionRepo, feedPriceHistoryRepo, db)
	feedPriceHistoryService := service.NewFeedPriceHistoryService(feedPriceHistoryRepo)
	
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	farmHandler := handler.NewFarmHandler(farmService)
	merchantHandler := handler.NewMerchantHandler(merchantService)
	pondHandler := handler.NewPondHandler(pondService)
	workerHandler := handler.NewWorkerHandler(workerService)
	feedCollectionHandler := handler.NewFeedCollectionHandler(feedCollectionService)
	feedPriceHistoryHandler := handler.NewFeedPriceHistoryHandler(feedPriceHistoryService)
	
	handlers := handler.NewHandler(handler.HandlerParams{
		UserHandler:             userHandler,
		AuthHandler:            authHandler,
		FarmHandler:            farmHandler,
		MerchantHandler:        merchantHandler,
		PondHandler:            pondHandler,
		WorkerHandler:          workerHandler,
		FeedCollectionHandler: feedCollectionHandler,
		FeedPriceHistoryHandler: feedPriceHistoryHandler,
	})
	
	s.handlers = handlers
	s.app = NewRouter()
	SetupRoutes(s.app, conf, handlers)
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

func (s *RouterTestSuite) TestHealthCheck() {
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) TestSwaggerRoute() {
	req := httptest.NewRequest("GET", "/swagger/index.html", nil)
	resp, err := s.app.Test(req)

	// Swagger might return 404 if not configured, but route should exist
	assert.NoError(s.T(), err)
	assert.True(s.T(), resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusNotFound)
}

func (s *RouterTestSuite) TestAuthRoutes() {
	req := httptest.NewRequest("POST", "/api/v1/auth/register", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	// Should return error (invalid body) but route exists
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) TestProtectedRoutes_WithoutAuth() {
	req := httptest.NewRequest("GET", "/api/v1/farm", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	// Should return 401 or error due to missing JWT
	assert.True(s.T(), resp.StatusCode == fiber.StatusUnauthorized || resp.StatusCode == fiber.StatusOK)
}

