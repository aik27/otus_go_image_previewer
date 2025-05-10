package config

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	App struct {
		Name     string `env:"APP_NAME" validate:"required"`
		Debug    string `env:"APP_DEBUG" validate:"required,oneof=true false"`
		Env      string `env:"APP_ENV" validate:"required,oneof=local production preproduction qa2"`
		LogLevel string `env:"APP_LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	}

	HTTPServer struct {
		ListenAddr   string `env:"HTTP_SERVER_LISTEN_ADDR" validate:"required,ip"`
		ListenPort   int    `env:"HTTP_SERVER_LISTEN_PORT" validate:"required,numeric"`
		ReadTimeout  int    `env:"HTTP_SERVER_READ_TIMEOUT" validate:"required,numeric"`
		WriteTimeout int    `env:"HTTP_SERVER_WRITE_TIMEOUT" validate:"required,numeric"`
		IdleTimeout  int    `env:"HTTP_SERVER_IDLE_TIMEOUT" validate:"required,numeric"`
	}

	Cache struct {
		Dir      string `env:"CACHE_DIR" validate:"required"`
		Capacity int    `env:"CACHE_CAPACITY" validate:"required,numeric"`
	}
}

func Load() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		panic("Config read error: " + err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		panic("Config validation error: " + err.Error())
	}

	// @TODO - read from vault

	return &cfg
}

func GetConfig() *Config {
	once.Do(func() {
		instance = Load()
	})

	return instance
}
