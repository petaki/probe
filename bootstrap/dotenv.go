package bootstrap

import (
	"log"

	"github.com/joho/godotenv"
)

// Dotenv function.
func Dotenv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
