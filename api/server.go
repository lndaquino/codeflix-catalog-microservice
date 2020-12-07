package api

import (
	"fmt"
	"log"
	"os"
	"video-catalog/controllers"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

// Run starts api services
func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error getting env, %v", err)
	}
	fmt.Println("Getting env values...")

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s\n", apiPort)
	server.Run(apiPort)
}
