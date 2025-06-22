package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title    string
	Author   string
	Tags     string
	Format   string
	FilePath string
		Shelves *[]Shelf `gorm:"many2many:book_shelves;"`
}


