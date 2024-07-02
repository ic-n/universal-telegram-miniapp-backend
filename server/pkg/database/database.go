package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func New(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
