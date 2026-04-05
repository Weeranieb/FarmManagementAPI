// src/internal/config/config.go
package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	Database       DatabaseConfig       `mapstructure:"database"`
	App            AppConfig            `mapstructure:"app"`
	Authentication AuthenticationConfig `mapstructure:"authentication"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type AppConfig struct {
	Environment         string `mapstructure:"environment"`
	LogLevel            string `mapstructure:"log_level"`
	Debug               bool   `mapstructure:"debug"`
	DailyLogUploadPath  string `mapstructure:"daily_log_upload_path"`
	DailyFeedUploadPath string `mapstructure:"daily_feed_upload_path"` // legacy alias, used if daily_log_upload_path is empty
}

type AuthenticationConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	JWTExpiry string `mapstructure:"jwt_expiry"`
}

// LoadConfig loads configuration using viper
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./configuration")

	// Set default values
	setDefaults()

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Failed to unmarshal config:", err)
	}
	if strings.TrimSpace(config.Authentication.JWTSecret) == "" {
		log.Fatal("authentication.jwt_secret must not be empty")
	}

	return &config
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "go_kit_base")
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.ssl_mode", "disable")

	// App defaults
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.debug", false)
	viper.SetDefault("app.daily_log_upload_path", "./data/uploads/daily-log")
	viper.SetDefault("app.daily_feed_upload_path", "./data/uploads/daily-feed")

	// Authentication defaults
	viper.SetDefault("authentication.jwt_secret", "")
	viper.SetDefault("authentication.jwt_expiry", "24h")
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
	)
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.App.Environment) == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.App.Environment) == "production"
}
