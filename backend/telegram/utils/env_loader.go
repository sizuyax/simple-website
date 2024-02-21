package utils

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	err := godotenv.Load("backend/config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
