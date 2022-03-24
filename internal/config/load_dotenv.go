package config

import (
	"os"

	"github.com/joho/godotenv"
)

const localEnv = ".env"

func LoadDotEnv() {
	_ = godotenv.Load()
	if _, err := os.Stat(localEnv); err == nil {
		_ = godotenv.Overload(localEnv)
	}
}
