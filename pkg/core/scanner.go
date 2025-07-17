package core

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type File struct {
	title string
	path  string
	ext   string
}

type Scanner struct {
	path  string
	Files []File
}

func NewScanner(path string) *Scanner {
	return &Scanner{
		path:  path,
		Files: make([]File, 0),
	}
}

func (s *Scanner) ScanDir() ([]File, error) {
	var files []File

	err := filepath.WalkDir(s.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Error walking directory: %v", err)
		}

		files = append(files, File{
			title: d.Name(),
			path:  path,
			ext:   filepath.Ext(d.Name()),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error extracting file data: %v", err)
	}

	return files, nil
}
