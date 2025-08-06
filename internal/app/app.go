package app

import (
	"fmt"
	
	"github.com/amirmtaati/libra/internal/service"
	"github.com/amirmtaati/libra/internal/storage"
	"github.com/amirmtaati/libra/internal/storage/repository"
	"github.com/amirmtaati/libra/pkg/scanner"
	"github.com/amirmtaati/libra/pkg/metadata"
	"path/filepath"
	"strings"

)

type App struct {
	DB           *storage.Database
	BookService  *service.BookService
	ShelfService *service.ShelfService
	Scanner      *scanner.Scanner
	Config       *Config
}

type Config struct {
	DBPath  string
	LibPath string
	Port    string
}

func NewApp(cfg *Config) *App {
	return &App{
		Config: cfg,
	}
}

func (a *App) Init() error {
	db, err := storage.NewDatabase(a.Config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	a.DB = db

	bookRepo := repository.NewBookRepository(db)
	shelfRepo := repository.NewShelfRepository(db)

	a.BookService = service.NewBookService(bookRepo)
	a.ShelfService = service.NewShelfService(shelfRepo, bookRepo)
	a.Scanner = scanner.NewScanner(a.Config.LibPath)

	return nil
}

func (a *App) Shutdown() error {
    if a.DB != nil {
        return a.DB.Close()
    }
    return nil
}

func (a *App) ScanLibrary() error {
    files, err := a.Scanner.ScanDir()
    if err != nil {
        return fmt.Errorf("error scanning directory: %v", err)
    }

    for _, file := range files {
        if strings.ToLower(filepath.Ext(file.Path)) != ".pdf" {
            continue
        }
        
        pdfMeta, err := metadata.ExtractPDFMetadata(file.Path)
        if err != nil {
            fmt.Printf("Warning: Failed to extract metadata for %s: %v\n", file.Path, err)
            pdfMeta = &metadata.PDFMetadata{
                Title: file.Title,
                PageCount: 0,
            }
        }
        
        _, err = a.BookService.Create(
            pdfMeta.Title,
            pdfMeta.Author, 
            pdfMeta.Keywords,
            file.Ext,
            file.Path,
        )
        if err != nil {
            fmt.Printf("Error creating book: %v\n", err)
            continue
        }
        
        fmt.Printf("✓ Imported: %s\n", pdfMeta.Title)
    }
    
    return nil
}
