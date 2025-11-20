package config

import (
	"path/filepath"
	"runtime"
	"time"

	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds the entire configuration for the application
type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	Logger         LoggerConfig         `mapstructure:"logger"`
	Redis          RedisConfig          `mapstructure:"redis"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port            int                `mapstructure:"port"`
	Host            string             `mapstructure:"host"`
	ReadTimeout     time.Duration      `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration      `mapstructure:"write_timeout"`
	AppName         string             `mapstructure:"app_name"`
	ShutdownTimeout time.Duration      `mapstructure:"shutdown_timeout"`
	RateLimiter     RateLimiterConfig  `mapstructure:"rate_limiter"`
}

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Enabled    bool          `mapstructure:"enabled"`
	Max        int           `mapstructure:"max"`
	Expiration time.Duration `mapstructure:"expiration"`
}


func Load() (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	// Set default values
	setDefaults(v)

	// Config File Settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(GetBaseDir())

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Environment Variables
	v.AutomaticEnv()

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetBaseDir() string {
	// Get the current file's directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}

	dir := filepath.Dir(filename)
	projectRoot := filepath.Dir(filepath.Dir(dir))
	return projectRoot
}

func GetServerAddress(cfg *Config) string {
	return cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.app_name", "Microservice Starter")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.shutdown_timeout", 10)
	v.SetDefault("server.rate_limiter.enabled", true)
	v.SetDefault("server.rate_limiter.max", 100)
	v.SetDefault("server.rate_limiter.expiration", 1*time.Minute)

	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")

	// Redis defaults
	v.SetDefault("redis.host", "redis")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
}
