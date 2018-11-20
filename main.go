package main // import "github.com/petaki/probe"

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Starting Probe.")

	redisPrefix := os.Getenv("REDIS_PREFIX")
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")
	redisDb := os.Getenv("REDIS_DB")

	fmt.Println(redisPrefix, redisHost, redisPassword, redisPort, redisDb)
}
