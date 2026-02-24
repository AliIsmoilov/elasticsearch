package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresData
}

type PostgresData struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       string
}

func NewConfig(path string) Config {
	godotenv.Load(path + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()
	return Config{
		Postgres: PostgresData{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			Username: conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			DB:       conf.GetString("POSTGRES_DB"),
		},
	}
}
