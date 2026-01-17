package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bookshelf-go/handlers"
	"bookshelf-go/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	// テスト用のSQLiteインメモリデータベースをセットアップ
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// マイグレーションを実行
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		return nil, err
	}
	return db, nil
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	bookHandler := handlers.NewBookHandler(db)

	// ルートを設定
	router.GET("/", bookHandler.ShowBooks)
	router.POST("/add", bookHandler.AddBook)
	router.POST("/delete/:id", bookHandler.DeleteBook)

	return router
}

func TestShowBooks(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	router := setupRouter(db)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddBook(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	router := setupRouter(db)

	// 書籍を追加
	payload := strings.NewReader("title=TestBook&author=TestAuthor&genre=TestGenre&status=未読")
	req, _ := http.NewRequest("POST", "/add", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code) // 成功時のリダイレクト

	// 書籍が追加されたかを確認
	var book models.Book
	result := db.First(&book)
	assert.NoError(t, result.Error)
	assert.Equal(t, "TestBook", book.Title)
	assert.Equal(t, "TestAuthor", book.Author)
}

func TestDeleteBook(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	router := setupRouter(db)

	// テストデータを追加
	book := models.Book{Title: "DeleteBook", Author: "TestAuthor", Genre: "TestGenre", Status: "未読"}
	db.Create(&book)

	// 書籍を削除
	req, _ := http.NewRequest("POST", fmt.Sprintf("/delete/%d", book.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code) // 成功時のリダイレクト

	// 書籍が削除されたかを確認
	var result models.Book
	err = db.First(&result, book.ID).Error
	assert.Error(t, err) // 見つからないはず
}