package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(url string) (*gorm.DB, error) {
	if url == "" {
		panic("DATABASE_URL cannot be empty")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
