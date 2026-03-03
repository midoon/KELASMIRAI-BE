package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigApp struct {
	Server   Server
	Database Database
	JWT      JWT
	Redis    Redis
	Midtrans Midtrans
}

func GetConfig() *ConfigApp {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &ConfigApp{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},

		JWT: JWT{
			Key:    os.Getenv("JWT_KEY"),
			Issuer: os.Getenv("JWT_ISSUER"),
		},

		Redis: Redis{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},

		Midtrans: Midtrans{
			Key:    os.Getenv("MIDTRANS_KEY"),
			IsProd: os.Getenv("MIDTRANS_ENV") == "Production",
		},
	}
}
