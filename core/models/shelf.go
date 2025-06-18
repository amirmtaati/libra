package models

import (
	"gorm.io/gorm"
)

type Shelf struct {
	gorm.Model
	Name  string
	Books []*Book `gorm:"many2many:book_shelves;"`
}
