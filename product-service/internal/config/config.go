package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string           `yaml:"env" envDefault:"local"`
	StorageConn string           `yaml:"storage_conn" env-required:"true"`
	HttpServer  HttpServerConfig `yaml:"http_server" env-required:"true"`
	Jwt         JwtConfig        `yaml:"jwt" env-required:"true"`
}

type HttpServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JwtConfig struct {
	SecretKey string `yaml:"secret_key"`
}

func MustLoad() *Config {
	configPath := "product-service/config/config-product.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file doesn't exists: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config: ", err)
	}

	return &cfg
}
