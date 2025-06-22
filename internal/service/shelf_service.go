package service

import (
	"fmt"
	"github.com/amirmtaati/libra/core/models"
	"github.com/amirmtaati/libra/core/repository"
	"strings"
)

type ShelfService struct {
	shelfRepo repository.ShelfRepository
	bookRepo  repository.BookRepository
}

func NewShelfService(shelfRepo repository.ShelfRepository, bookRepo repository.BookRepository) *ShelfService {
	return &ShelfService{
		shelfRepo: shelfRepo,
		bookRepo:  bookRepo,
	}
}

func (s *ShelfService) Create(name string) (*models.Shelf, error) {
	shelf := &models.Shelf{
		Name: strings.TrimSpace(name),
	}

	if err := s.shelfRepo.Create(shelf); err != nil {
		return nil, fmt.Errorf("Failed to create shelf: %v", err)
	}

	return shelf, nil
}

func (s *ShelfService) GetByID(id uint) (*models.Shelf, error) {
	shelf, err := s.shelfRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get shelf by id: %v", err)
	}

	return shelf, nil
}

func (s *ShelfService) GetAll() ([]*models.Shelf, error) {
	shelves, err := s.shelfRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to get all shelves: %v", err)
	}

	return shelves, nil
}

func (s *ShelfService) Delete(id uint) error {
	if err := s.shelfRepo.Delete(id); err != nil {
		return fmt.Errorf("Failed to delete shelf: %v", err)
	}
	return nil
}

func (s *ShelfService) AddBook(shelfID uint, bookID uint) error {
	shelf, err := s.shelfRepo.GetByID(shelfID)
	if err != nil {
		return fmt.Errorf("Failed to get shelf: %v", err)
	}
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return fmt.Errorf("Failed to get book: %v", err)
	}

	if err := s.shelfRepo.AddBook(shelf, book); err != nil {
		return fmt.Errorf("Failed to add book to shelf: %v", err)
	}

	return nil
}

func (s *ShelfService) GetAllBooks(shelfID uint) ([]*models.Book, error) {
	books, err := s.shelfRepo.GetAllBooks(shelfID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get all books: %v", err)
	}

	return books, nil
}

func (s *ShelfService) RemoveBook(shelfID uint, bookID uint) error {
	shelf, err := s.shelfRepo.GetByID(shelfID)
	if err != nil {
		return fmt.Errorf("Failed to get shelf: %v", err)
	}
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return fmt.Errorf("Failed to get book: %v", err)
	}

	if err := s.shelfRepo.RemoveBook(shelf, book); err != nil {
		return fmt.Errorf("Failed to add book to shelf: %v", err)
	}

	return nil

}
