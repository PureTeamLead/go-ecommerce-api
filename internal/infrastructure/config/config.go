package config

import (
	httpServer "eshop/internal/transport/http-server"
	"eshop/pkg/postgre"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env string               `yaml:"env" env-default:"prod"`
	DB  postgre.DBconfig     `yaml:"db" env-required:"true"`
	App httpServer.AppConfig `yaml:"app" env-required:"true"`
}

func LoadConfig(path string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal("error loading config file: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("error loading env variables: " + err.Error())
	}

	return &cfg
}
