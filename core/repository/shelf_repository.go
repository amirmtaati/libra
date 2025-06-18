package repository

import (
	"fmt"
	"github.com/amirmtaati/libra/core/models"
	"github.com/amirmtaati/libra/core/storage"
)

type ShelfRepository interface {
		Create(shelf *models.Shelf) error
		GetByID(id uint) (*models.Shelf, error)
		GetAll() ([]*models.Shelf, error)
		Update(shelf *models.Shelf) error
		Delete(id uint) error
		
		Add(shelf *models.Shelf, book *models.Book) error
		GetAllBooks(id uint)
		Remove(id uint) error

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
	var shelfs []*models.Shelf
	if err := r.db.DB.First(&shelfs).Error; err != nil {
		return nil, fmt.Errorf("Failed to get all shelfs: %v", err)
	}
	return shelfs, nil

}

func (r *shelfRepository) Update(shelf*models.Shelf) error {
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


func (r *shelfRepository) Add(shelf *models.Shelf, book *models.Book) error {
		r.db.DB.Model(shelf).Select("books").Update()
}
