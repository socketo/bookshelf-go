package models

type Book struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"size:100"`
	Author string `gorm:"size:100"`
	Genre  string `gorm:"size:50"`
	Status string `gorm:"size:50"` // e.g., "Read" or "Unread"
}
