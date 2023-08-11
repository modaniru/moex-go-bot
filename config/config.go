package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Level string
	Token string
}

func LoadConfig() *Config {
	log.Println("loading configuration")
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file was not found")
	} else {
		log.Println(".env file was found successfully")
	}
	log.Println("insertData to structure")
	cfg := Config{
		Level: os.Getenv("LEVEL"),
		Token: os.Getenv("BOT_TOKEN"),
	}
	log.Println("config was successfully loaded")
	return &cfg
}
