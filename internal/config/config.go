package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" default:"local"`
	HTTPServer HTTPServer `yaml:"http_server" env:"HTTP_SERVER"`
	Database   Database   `yaml:"database" env:"DATABASE"`
}

type Database struct {
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"USER" env-default:"user"`
	Password string `yaml:"password" env:"PASSWORD"`
	Address  string
}

type HTTPServer struct {
	Address     string
	Host        string        `yaml:"host" env:"HTTP_HOST" default:"localhost"`
	Port        string        `yaml:"port" env:"HTTP_PORT" default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" default:"5"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" default:"5"`
}

const PathKey = "CONFIG_PATH"

const Local = "local"
const Prod = "prod"

func New() (*Config, error) {
	configPath := os.Getenv(PathKey)

	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", configPath)
		}
		return nil, fmt.Errorf("error checking config file: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	cfg.HTTPServer.Address = fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	cfg.Database.Address = fmt.Sprintf(
		"postgres://%s:@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	return &cfg, nil
}
