package main

import (
	"boonmafarm/api/controllers"
	"boonmafarm/api/middlewares"
	"boonmafarm/api/pkg/repositories"
	"boonmafarm/api/pkg/services"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Entrypoint for app.
func main() {
	viper.SetConfigName("config") // get config filename
	viper.AddConfigPath(".")      // set path file config
	viper.AutomaticEnv()          // set ENV variable

	// read config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// connection db to gorm
	cfg := viper.GetString("postgres.connection")
	dsn := cfg
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return
	}

	// gin
	router := gin.Default()

	// cors
	router.Use(middlewares.Cors())

	// jwt authentication
	router.Use(middlewares.JWTAuthMiddleware())

	// repositories
	userRepo := repositories.NewUserRepository(db)
	clientRepo := repositories.NewClientRepository(db)
	farmRepo := repositories.NewFarmRepository(db)
	farmGroupRepo := repositories.NewFarmGroupRepository(db)
	farmOnFarmGroupRepo := repositories.NewFarmOnFarmGroupRepository(db)
	pondRepo := repositories.NewPondRepository(db)
	activePondRepo := repositories.NewActivePondRepository(db)
	activityRepo := repositories.NewActivityRepository(db)
	sellDetailRepo := repositories.NewSellDetailRepository(db)

	// services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	clientService := services.NewClientService(clientRepo)
	farmService := services.NewFarmService(farmRepo)
	farmGroupService := services.NewFarmGroupService(farmGroupRepo)
	farmOnFarmGroupService := services.NewFarmOnFarmGroupService(farmOnFarmGroupRepo)
	pondService := services.NewPondService(pondRepo)
	activePondService := services.NewActivePondService(activePondRepo)
	activityService := services.NewActivityService(activityRepo, sellDetailRepo)

	// controllers
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	clientController := controllers.NewClientController(clientService)
	farmController := controllers.NewFarmController(farmService)
	farmGroupController := controllers.NewFarmGroupController(farmGroupService)
	farmOnFarmGroupController := controllers.NewFarmOnFarmGroupController(farmOnFarmGroupService)
	pondController := controllers.NewPondController(pondService)
	activePondController := controllers.NewActivePondController(activePondService)
	activityController := controllers.NewActivityController(activityService)

	// apply route
	userController.ApplyRoute(router)
	authController.ApplyRoute(router)
	clientController.ApplyRoute(router)
	farmController.ApplyRoute(router)
	farmGroupController.ApplyRoute(router)
	farmOnFarmGroupController.ApplyRoute(router)
	pondController.ApplyRoute(router)
	activePondController.ApplyRoute(router)
	activityController.ApplyRoute(router)

	// run server
	router.Run(":8080")
}
