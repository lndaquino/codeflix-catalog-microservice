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
	uuid "github.com/satori/go.uuid"
)

var server = controllers.Server{}
var categoryInstance = models.Category{}

func TestMain(m *testing.M) {
	if _, err := os.Stat("./../.env"); os.IsNotExist(err) {
		log.Fatalf(".env not found: %v\n", err)
	} else {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../.env"))
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

func refreshCategoryTable() error {
	err := server.DB.DropTableIfExists(&models.Category{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Category{}).Error
	if err != nil {
		return err
	}
	log.Printf("Sucessfully refreshed table")
	return nil
}

func seedCategories() ([]models.Category, error) {
	var err error
	x := true
	categories := []models.Category{
		models.Category{
			ID:          uuid.NewV4().String(),
			Name:        "category 1",
			Description: "description 1",
			IsActive:    &x,
		},
		models.Category{
			ID:          uuid.NewV4().String(),
			Name:        "category 2",
			Description: "description 2",
			IsActive:    &x,
		},
	}

	for i := range categories {
		err = server.DB.Model(&models.Category{}).Create(&categories[i]).Error
		if err != nil {
			return []models.Category{}, err
		}
	}
	return categories, nil
}

func seedOneCategory() (models.Category, error) {
	x := true
	category := models.Category{
		ID:          uuid.NewV4().String(),
		Name:        "category name",
		Description: "describing category",
		IsActive:    &x,
	}

	err := server.DB.Model(&models.Category{}).Create(&category).Error
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}
