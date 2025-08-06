package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type File struct {
	Title string
	Path  string
	Ext   string
}

type Scanner struct {
	path  string
	Files []File
}

func NewFile(title, path, ext string) File {
	return File{
		Title: title,
		Path:  path,
		Ext:   ext,
	}
}

func NewScanner(path string) *Scanner {
	return &Scanner{
		path: path,
	}
}

func (s *Scanner) ScanDir() ([]File, error) {
	var files []File

	err := filepath.WalkDir(s.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Error walking directory: %v", err)
		}

		if !d.IsDir() {
			files = append(files, File{
				Title: d.Name(),
				Path:  path,
				Ext:   filepath.Ext(d.Name()),
			})

		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error extracting file data: %v", err)
	}

	return files, nil
}
