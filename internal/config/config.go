package config

import (
"log"
"os"

"github.com/joho/godotenv"
)

type Config struct {
ENV        string
ServerPort string
DBHost     string
DBPort     string
DBUser     string
DBPassword string
DBName     string
}

var AppConfig Config

func LoadConfig() {
err := godotenv.Load()
if err != nil {
log.Println("No .env file found, using environment variables")
}

AppConfig = Config{
ENV:        os.Getenv("ENV"),
ServerPort: os.Getenv("PORT"),
DBHost:     os.Getenv("DB_HOST"),
DBPort:     os.Getenv("DB_PORT"),
DBUser:     os.Getenv("DB_USER"),
DBPassword: os.Getenv("DB_PASSWORD"),
DBName:     os.Getenv("DB_NAME"),
}
}
