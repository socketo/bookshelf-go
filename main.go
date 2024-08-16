package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"size:100"`
	Author string `gorm:"size:100"`
	Genre  string `gorm:"size:50"`
	Status string `gorm:"size:50"` // e.g., "Read" or "Unread"
}

var db *gorm.DB
var err error

func main() {
	// Initialize Gin & setup DB
	router := gin.Default()
	db, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Book{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Routes
	router.LoadHTMLGlob("templates/*")
	router.GET("/", showBooks)
	router.POST("/add", addBook)
	router.POST("/delete/:id", deleteBook)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// showBooks renders the book list
func showBooks(c *gin.Context) {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"books": books,
	})
}

// addBook handles adding a new book
func addBook(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	genre := c.PostForm("genre")
	status := c.PostForm("status")

	book := Book{Title: title, Author: author, Genre: genre, Status: status}
	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}

// deleteBook handles deleting a book
func deleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}
