package config

import (
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all environment variables
type Config struct {
	MONGODB_URI        string
	DATABASE_NAME      string
	OPEN_AI_API_KEY    string
	PROMPT_TEMPLATE    string
	SECRET_ACCESS_KEY  string
	SECRET_REFRESH_KEY string
}

// Env is the global config instance
var Env Config

func init() {
	// Load .env (ignore if not found)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ .env file not found, using system environment")
	}

	v := reflect.ValueOf(&Env).Elem()
	t := v.Type()

	var missingEnvVariables string

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		envVariableName := field.Name
		value := os.Getenv(envVariableName)

		if value == "" {
			missingEnvVariables += "❌ " + envVariableName + "\n"
		} else {
			if fieldValue.Kind() == reflect.String {
				fieldValue.SetString(value)
			} else if fieldValue.Kind() == reflect.Int64 {
				parsed, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					log.Fatalf("❌ Invalid value for %s: %s", envVariableName, value)
				}
				fieldValue.SetInt(parsed)
			} else {
				log.Fatalf("Unsupported field type for %s", field.Name)
			}
		}
	}

	if missingEnvVariables != "" {
		log.Fatalf("Missing required environment variables:\n%s", missingEnvVariables)
	}
}
