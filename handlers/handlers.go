package handlers

import (
	"net/http"

	"bookshelf-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookHandler struct {
	DB *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{DB: db}
}

// ShowBooks renders the book list
func (h *BookHandler) ShowBooks(c *gin.Context) {
	var books []models.Book
	if err := h.DB.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"books": books,
	})
}

// AddBook handles adding a new book
func (h *BookHandler) AddBook(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	genre := c.PostForm("genre")
	status := c.PostForm("status")

	book := models.Book{Title: title, Author: author, Genre: genre, Status: status}
	if err := h.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}

// DeleteBook handles deleting a book
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}
