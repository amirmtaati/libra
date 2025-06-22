package repository

import (
	"fmt"
	"github.com/amirmtaati/libra/core/models"
	"github.com/amirmtaati/libra/core/storage"
)

type BookRepository interface {
	Create(book *models.Book) error
	GetByID(id uint) (*models.Book, error)
	GetAll() ([]*models.Book, error)
	Update(book *models.Book) error
	Delete(id uint) error
}

type bookRepository struct {
	db *storage.Database
}

func NewBookRepository(db *storage.Database) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *models.Book) error {
	return r.db.DB.Create(book).Error
}

func (r *bookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	if err := r.db.DB.First(&book, id).Error; err != nil {
		return nil, fmt.Errorf("Failed to get book by ID: %v", err)
	}

	return &book, nil
}

func (r *bookRepository) GetAll() ([]*models.Book, error) {
	var books []*models.Book
	if err := r.db.DB.First(&books).Error; err != nil {
		return nil, fmt.Errorf("Failed to get all books: %v", err)
	}
	return books, nil

}

func (r *bookRepository) Update(book *models.Book) error {
	if err := r.db.DB.Save(book).Error; err != nil {
		return err
	}
	return nil
}

func (r *bookRepository) Delete(id uint) error {
	var book models.Book
	if err := r.db.DB.Delete(book, id).Error; err != nil {
		return err
	}
	return nil
}
