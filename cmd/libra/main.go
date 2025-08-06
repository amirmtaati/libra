package main

import (
	"fmt"
	"log"

	"github.com/amirmtaati/libra/internal/service"
	"github.com/amirmtaati/libra/internal/storage"
	"github.com/amirmtaati/libra/internal/storage/repository"
	"github.com/amirmtaati/libra/pkg/metadata"
	"github.com/amirmtaati/libra/pkg/scanner"
)

func main() {
	db, err := storage.NewDatabase(".libra")
	if err != nil {
		fmt.Errorf("Error creating DB: %v", err)
	}

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	scanner := scanner.NewScanner("/home/amirmhmd/Documents/library")
	files, err := scanner.ScanDir()
	if err != nil {
		fmt.Errorf("Error scanning directory: %v", err)
	}

	for _, file := range files {
		pdfMetadata, err := metadata.ExtractPDFMetadata(file.Path)
		if err != nil {
			log.Fatal("ERROR")
		}
		book, err := bookService.Create(pdfMetadata.Title, pdfMetadata.Author, "", file.Ext, file.Path)
		if err != nil {
			log.Fatal("ERROR")
		}
		fmt.Println(pdfMetadata.Title, book.Title, book.Format)
	}
}
