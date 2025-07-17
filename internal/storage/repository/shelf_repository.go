package repository

import (
	"fmt"
	"github.com/amirmtaati/libra/internal/storage/models"
	"github.com/amirmtaati/libra/internal/storage"
)

type ShelfRepository interface {
	Create(shelf *models.Shelf) error
	GetByID(id uint) (*models.Shelf, error)
	GetAll() ([]*models.Shelf, error)
	Update(shelf *models.Shelf) error
	Delete(id uint) error
	AddBook(shelf *models.Shelf, book *models.Book) error
	GetAllBooks(shelfID uint) ([]*models.Book, error)
	RemoveBook(shelf *models.Shelf, book *models.Book) error
}

type shelfRepository struct {
	db *storage.Database
}

func NewShelfRepository(db *storage.Database) ShelfRepository {
	return &shelfRepository{db: db}
}

func (r *shelfRepository) Create(shelf *models.Shelf) error {
	return r.db.DB.Create(shelf).Error
}

func (r *shelfRepository) GetByID(id uint) (*models.Shelf, error) {
	var shelf models.Shelf
	if err := r.db.DB.First(&shelf, id).Error; err != nil {
		return nil, fmt.Errorf("Failed to get shelf by ID: %v", err)
	}

	return &shelf, nil
}

func (r *shelfRepository) GetAll() ([]*models.Shelf, error) {
	var shelves []*models.Shelf
	if err := r.db.DB.First(&shelves).Error; err != nil {
		return nil, fmt.Errorf("Failed to get all shelfs: %v", err)
	}
	return shelves, nil

}

func (r *shelfRepository) Update(shelf *models.Shelf) error {
	if err := r.db.DB.Save(shelf).Error; err != nil {
		return err
	}
	return nil
}

func (r *shelfRepository) Delete(id uint) error {
	var shelf models.Shelf
	if err := r.db.DB.Delete(shelf, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *shelfRepository) AddBook(shelf *models.Shelf, book *models.Book) error {
	if err := r.db.DB.Model(&shelf).Association("Books").Append(&book); err != nil {
		return fmt.Errorf("Failed to add book to shelf: %v", err)
	}

	return nil
}

func (r *shelfRepository) GetAllBooks(shelfID uint) ([]*models.Book, error) {
	var shelf models.Shelf
	if err := r.db.DB.Preload("Books").First(&shelf, shelfID).Error; err != nil {
		return nil, fmt.Errorf("Failed to get all books: %v", err)
	}

	return shelf.Books, nil
}

func (r *shelfRepository) RemoveBook(shelf *models.Shelf, book *models.Book) error {
	if err := r.db.DB.Model(&shelf).Association("Books").Delete(&book); err != nil {
		return fmt.Errorf("Failed to delete book from shelf: %v", err)
	}

	return nil
}
