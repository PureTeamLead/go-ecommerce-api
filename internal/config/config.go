package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type database struct {
	Username string `yaml:"db_username" env-required:"true"`
	Password string `yaml:"db_password" env-required:"true"`
	Host     string `yaml:"db_host" env-required:"true"`
	Port     string `yaml:"db_port" env-required:"true"`
	Name     string `yaml:"db_name" env-required:"true"`
}

type app struct {
	Host        string        `yaml:"serv_host" env-required:"true"`
	Port        string        `yaml:"serv_port" env-required:"true"`
	Username    string        `yaml:"serv_username" env-required:"true"`
	Password    string        `yaml:"serv_password" env-required:"true" env:"HTTP_USER_PASSWORD"`
	RWTimeout   time.Duration `yaml:"rw_timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

// Config TODO: refactor with all config types in cleanenv
// Config TODO: change types of config env, make different configs
type Config struct {
	Env string   `yaml:"env" env-default:"prod"`
	DB  database `yaml:"db"`
	App app      `yaml:"app"`
}

func LoadConfig(path string) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatal("error loading config: " + err.Error())
	}

	return &cfg
}

// MakeUrlDB TODO:ssl mode to config
func (c *Config) MakeUrlDB() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=public", c.DB.Username, c.DB.Password,
		c.DB.Host, c.DB.Port, c.DB.Name)
}
