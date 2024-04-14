package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port              string `envconfig:"PORT" default:"4000"`
	DatabaseHost      string `envconfig:"DATABASE_HOST" default:"localhost"`
	DatabasePort      string `envconfig:"DATABASE_PORT" default:"5432"`
	DatabaseName      string `envconfig:"DATABASE_NAME" default:"todo.db"`
	DatabaseUser      string `envconfig:"DATABASE_USER" default:"admin"`
	DatabasePassword  string `envconfig:"DATABASE_PASSWORD" default:"password"`
	RedisHost         string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort         string `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword     string `envconfig:"REDIS_PASSWORD" default:"password"`
	SessionCookieName string `envconfig:"SESSION_COOKIE_NAME" default:"session_cookie"`
}

func loadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
