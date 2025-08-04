package main

import (
	"log"
	"os"

	"github.com/avakili/data-ingestion/backend/routes"
	"github.com/avakili/data-ingestion/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
	dataStorageService := services.NewDataPointStorageServiceImpl(db)
	r := gin.Default()
	routes.DataPointRoutes(r, dataStorageService)
	r.Run()
}
