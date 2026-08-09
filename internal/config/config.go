package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ThirdPartyURL string
	Port          string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	return &Config{
		ThirdPartyURL: os.Getenv("URL"),
		Port:          os.Getenv("PORT"),
	}
}
