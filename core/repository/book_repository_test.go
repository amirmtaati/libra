package repository

import (
	"path/filepath"
	"testing"

	"github.com/amirmtaati/libra/core/models"
	"github.com/amirmtaati/libra/core/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) (*storage.Database, error) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := storage.NewDatabase(dbPath)
	return db, err
}

func TestBookRepository_Create(t *testing.T) {
	db, err := setupDB(t)
	require.NoError(t, err)
	defer db.Close()

	err = db.DB.AutoMigrate(&models.Book{})
	require.NoError(t, err)

	repo := NewBookRepository(db)

	book := &models.Book{
		Title:    "Test Book",
		Author:   "Test Author",
		Tags:     "fiction,test",
		Format:   "PDF",
		FilePath: "/path/to/book.pdf",
	}

	repo.Create(book)

	assert.NoError(t, err, "Create should not return an error")
	assert.NotZero(t, book.ID, "Book ID should be set after creation")
	assert.NotZero(t, book.CreatedAt, "CreatedAt should be set")
	assert.NotZero(t, book.UpdatedAt, "UpdatedAt should be set")

	var savedBook models.Book
	err = db.DB.First(&savedBook, book.ID).Error
	require.NoError(t, err, "Should be able to find the created book")

	assert.Equal(t, book.Title, savedBook.Title)
	assert.Equal(t, book.Author, savedBook.Author)
	assert.Equal(t, book.Tags, savedBook.Tags)
	assert.Equal(t, book.Format, savedBook.Format)
	assert.Equal(t, book.FilePath, savedBook.FilePath)

}

func TestBookRepository_GetByID(t *testing.T) {
	db, err := setupDB(t)
	require.NoError(t, err)
	defer db.Close()

	err = db.DB.AutoMigrate(&models.Book{})
	require.NoError(t, err)

	repo := NewBookRepository(db)
	book := &models.Book{
		Title:    "Test Book",
		Author:   "Test Author",
		Tags:     "fiction,test",
		Format:   "PDF",
		FilePath: "/path/to/book.pdf",
	}

	repo.Create(book)

	book, err = repo.GetByID(1)
	assert.NoError(t, err, "GetByID should not return an error")
	assert.NotZero(t, book.ID, "Book ID should be set after creation")
	assert.NotZero(t, book.CreatedAt, "CreatedAt should be set")
	assert.NotZero(t, book.UpdatedAt, "UpdatedAt should be set")
}
