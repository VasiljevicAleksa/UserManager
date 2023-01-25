package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// open database method, returns gorm db wrapper or error if ocurred
func OpenDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	return db, nil
}
