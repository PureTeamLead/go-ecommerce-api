package config

import (
	httpServer "eshop/internal/transport/http-server"
	"eshop/pkg/postgre"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

//TODO: make env-required and env-default, validation

type Config struct {
	Env string               `yaml:"env" env-default:"prod"`
	DB  postgre.DBconfig     `yaml:"db"`
	App httpServer.AppConfig `yaml:"app"`
}

func LoadConfig(path string) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatal("error loading config: " + err.Error())
	}

	return &cfg
}
