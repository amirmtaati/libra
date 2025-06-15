package storage

import (
	"fmt"
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
