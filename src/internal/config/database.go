package config

import (
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB connects to the database using the config
func (c *Config) ConnectDB() *gorm.DB {
	dsn := c.GetDSN()

	var logLevel logger.LogLevel
	switch strings.ToLower(c.App.LogLevel) {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
	return db
}
