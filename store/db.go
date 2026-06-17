
package store

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"learning-insight-coach/models"
)

var DB *gorm.DB

func Init(dbPath string) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.Document{},
		&models.DocumentChunk{},
		&models.StudentInteraction{},
	); err != nil {
		return err
	}

	DB = db

	log.Println("INFO: database ready at", dbPath)

	return nil
}