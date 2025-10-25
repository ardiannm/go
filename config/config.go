package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MONGODB_URI        string
	DATABASE_NAME      string
	OPEN_AI_API_KEY    string
	PROMPT_TEMPLATE    string
	SECRET_ACCESS_KEY  string
	SECRET_REFRESH_KEY string
)

func init() {
	// Load .env once

	err := godotenv.Load(".env")

	if err != nil {
		log.Println("⚠️ .env file not found, using system environment")
	}

	// Read required env vars

	MONGODB_URI = os.Getenv("MONGODB_URI")
	DATABASE_NAME = os.Getenv("DATABASE_NAME")
	OPEN_AI_API_KEY = os.Getenv("OPEN_AI_API_KEY")
	PROMPT_TEMPLATE = os.Getenv("PROMPT_TEMPLATE")
	SECRET_ACCESS_KEY = os.Getenv("SECRET_ACCESS_KEY")
	SECRET_REFRESH_KEY = os.Getenv("SECRET_REFRESH_KEY")

	// Validate critical vars

	if MONGODB_URI == "" {
		log.Fatal("❌ Missing MONGODB_URI in environment")
	}

	if DATABASE_NAME == "" {
		log.Fatal("❌ Missing DATABASE_NAME in environment")
	}

	if OPEN_AI_API_KEY == "" {
		log.Fatal("❌ Missing OPEN_AI_API_KEY in environment")
	}

	if SECRET_ACCESS_KEY == "" {
		log.Fatal("❌ Missing SECRET_ACCESS_KEY in environment")
	}

	if SECRET_ACCESS_KEY == "" {
		log.Fatal("❌ Missing SECRET_ACCESS_KEY in environment")
	}

	if PROMPT_TEMPLATE == "" {
		log.Fatal("❌ Missing PROMPT_TEMPLATE in environment")
	}
}
