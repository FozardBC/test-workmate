package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Log        string `env:"LOG_MODE" env-default:"debug"`
	ServerHost string `env:"SRV_HOST"`
	ServerPort string `env:"SRV_PORT" env-default:"8080"`
}

func MustRead() *Config {

	if err := godotenv.Load(); err != nil { // DEBUG:
		panic(err)
	}

	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		help, _ := cleanenv.GetDescription(cfg, nil)
		log.Print(help)
		log.Fatal(err)
	}

	return &cfg
}
