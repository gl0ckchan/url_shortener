package config

import (
	"os"
	"time"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 	   string `yaml:"env" env:"ENV" env-required:"true" env-default:"dev"`
	Postgres   `yaml:"postgres"` 
	HTTPServer `yaml:"http_server" `
}

type Postgres struct {
	Host            string        `yaml:"host" env:"PG_HOST" env-default:"localhost"`
	Port            int           `yaml:"port" env:"PG_PORT" env-default:"5432"`
	User            string        `yaml:"user" env:"PG_USER" env-required:"true"`
	Password        string        `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
	DBName          string        `yaml:"dbname" env:"PG_DBNAME" env-required:"true"`
	SSLMode         string        `yaml:"sslmode" env:"PG_SSLMODE" env-default:"disable"`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"PG_MAX_OPEN_CONNS" env-default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"PG_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"PG_CONN_MAX_LIFETIME" env-default:"5m"`
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