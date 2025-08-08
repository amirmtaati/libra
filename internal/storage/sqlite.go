package storage

import (
	"fmt"
	"github.com/amirmtaati/libra/internal/storage/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to crrate database: %v", err)
	}

	// Automigrate core models and the join table for many-to-many relations
	if err := db.AutoMigrate(&models.Book{}, &models.Shelf{}); err != nil {
		return nil, fmt.Errorf("Failed to migrate database: %v", err)
	}

	database := &Database{DB: db}

	return database, nil
}

func (d *Database) Ping() error {
	sql, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sql.Ping()
}

func (d *Database) Close() error {
	sql, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sql.Close()
}
