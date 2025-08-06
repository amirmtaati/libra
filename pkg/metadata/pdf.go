package metadata

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"os"
	"strings"
	"time"
    "path/filepath"
)

type PDFMetadata struct {
	Title        string
	Author       string
	Subject      string
	Keywords     string
	CreationDate time.Time
	PageCount    int
}

func ExtractPDFMetadata(filePath string) (*PDFMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF file: %v", err)
	}
	defer file.Close()

	pdfInfo, err := api.PDFInfo(file, filePath, nil, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to extract PDF info: %v", err)
	}

	file.Seek(0, 0) // Reset file pointer
	properties, err := api.Properties(file, nil)
	if err != nil {
		properties = make(map[string]string)
	}

	// Create metadata with fallback to filename if title is empty
	title := strings.TrimSpace(properties["Title"])
	if title == "" {
		// Fallback to filename without extension as title
		title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	return &PDFMetadata{
		Title:     title,
		Author:    strings.TrimSpace(properties["Author"]),
		Subject:   strings.TrimSpace(properties["Subject"]),
		Keywords:  strings.TrimSpace(properties["Keywords"]),
		PageCount: pdfInfo.PageCount,
	}, nil

}
