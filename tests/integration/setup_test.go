package tests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"video-catalog/controllers"
	"video-catalog/models"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var categoryInstance = models.Category{}

func TestMain(m *testing.M) {
	if _, err := os.Stat("./../../.env"); os.IsNotExist(err) {
		log.Fatalf(".env not found: %v\n", err)
	} else {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		Database()
	}
	os.Exit(m.Run())
}

func Database() {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	server.DB, err = gorm.Open(os.Getenv("DB_DRIVER"), DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", os.Getenv("DB_DRIVER"))
		log.Fatal("Error connecting to database: ", err)
	} else {
		fmt.Printf("Connected to %s database\n", os.Getenv("DB_DRIVER"))
	}
}
