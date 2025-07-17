package service

import (
	"fmt"
	"github.com/amirmtaati/libra/internal/storage/models"
	"github.com/amirmtaati/libra/internal/storage/repository"
	"strings"
)

type BookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(r repository.BookRepository) *BookService {
	return &BookService{bookRepo: r}
}

func (s *BookService) Create(title, author, tags, format, filePath string) (*models.Book, error) {
	book := &models.Book{
		Title:    strings.TrimSpace(title),
		Author:   strings.TrimSpace(author),
		Tags:     strings.TrimSpace(tags),
		Format:   strings.ToUpper(format),
		FilePath: filePath,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, fmt.Errorf("Failed to create book: %v", err)
	}

	return book, nil
}

func (s *BookService) GetByID(id uint) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get book by ID: %v", err)
	}

	return book, nil
}

func (s *BookService) GetAll() ([]*models.Book, error) {
	books, err := s.bookRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to get all books: %v", err)
	}

	return books, nil
}

func (s *BookService) Delete(id uint) error {
	err := s.bookRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("Failed to delete book: %v", err)
	}

	return nil
}
