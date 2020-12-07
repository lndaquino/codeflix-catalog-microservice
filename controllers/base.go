package controllers

import (
	"fmt"
	"log"
	"net/http"
	"video-catalog/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver
)

// Server models server structure
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

//Initialize inits db connection and system routes
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("Error connecting to database: ", err)
	} else {
		fmt.Printf("Connected to %s database\n", Dbdriver)
	}

	//database migration
	server.DB.Debug().AutoMigrate(
		&models.Category{},
	)

	server.Router = gin.Default()
	server.initializeRoutes()
}

// Run starts server
func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
