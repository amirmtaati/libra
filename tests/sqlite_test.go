package test 

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/amirmtaati/libra/internal/storage"
	"path/filepath"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := storage.NewDatabase(dbPath)
	defer db.Close()

	require.NoError(t, err, "Failed to create database")
	require.NotNil(t, db, "Database should not be nil")
	require.NotNil(t, db.DB, "db.DB should not be nil")

	err = db.Ping()
	assert.NoError(t, err, "Failed to ping the database")
}
