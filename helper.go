package goopentelemetry

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvOrDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load file .env: ", err)
		panic(err)
	}

	value := GetEnvOrDefault(key, "").(string)
	return value
}

func StringToBool(value string) bool {
	if value == "true" {
		return true
	}

	return false
}
