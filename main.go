package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/gorm"

	"expense-app-backend/config"
	"expense-app-backend/models"
	"expense-app-backend/routes"
)

var db *gorm.DB

func initDatabase() {
	var err error
	db, err = config.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&models.Category{}, &models.SubCategory{}, &models.User{})
	fmt.Println("Database connected and table migrated")
}

func main() {
	initDatabase()

	r := routes.SetupRouter(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
