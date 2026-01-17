package main

import (
	"log"

	"bookshelf-go/database"
	"bookshelf-go/handlers"
	"bookshelf-go/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin & setup DB
	router := gin.Default()
	db := database.InitDB("books.db")

	// Migrate the schema
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	bookHandler := handlers.NewBookHandler(db)

	// Routes
	router.LoadHTMLGlob("templates/*")
	router.GET("/", bookHandler.ShowBooks)
	router.POST("/add", bookHandler.AddBook)
	router.POST("/delete/:id", bookHandler.DeleteBook)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}