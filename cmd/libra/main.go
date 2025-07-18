package main

import (
	"fmt"
	"github.com/amirmtaati/libra/internal/service"
	"github.com/amirmtaati/libra/internal/storage"
	"github.com/amirmtaati/libra/internal/storage/repository"
	"github.com/amirmtaati/libra/pkg/core"
)

func main() {
	db, err := storage.NewDatabase(".libra")
	if err != nil {
		fmt.Errorf("Error creating DB: %v", err)
	}

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	scanner := core.NewScanner("/home/amirmamad/Documents/library")
	files, err := scanner.ScanDir()
	if err != nil {
		fmt.Errorf("Error scanning directory: %v", err)
	}

	for _, file := range files {
		book, err := bookService.Create(file.Title, "author", "tags", file.Ext, file.Path)
		if err != nil {
			fmt.Errorf("Error creating a book: %v", err)
		}
		fmt.Println(book.Title)
	}
}
