package config

import (
	"os"
	"time"
	"log"
	// "log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 		string `yaml:"env" env:"ENV" env-required:"true" env-default:"dev" `
	StoragePath string `yaml:"storage_path" env:"STORAGEPATH" env-required:"true"`
	HTTPServer  `yaml:"http_server" `
}

type HTTPServer struct {
	Address  	 string `yaml:"address" env:"ADDRESS" env-required:"true" env-default:"localhost:8080"`
	Timeout 	 time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLETIMEOUT" env-default:"10s"`
	User 		 string `yaml:"user" env-required:"true"`
	Password     string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalln("config file is not exist:", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalln("failed to read config file:", err)
	}

	return &cfg
}