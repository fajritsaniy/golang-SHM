package common

import (
	"errors"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return errors.New("failed to load .env file")
	}
	return nil
}
