package database

import (
	"log"

	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	checkTable := db.Migrator().HasTable("rss")
	if !checkTable {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&RSS{})
		if err != nil {
			log.Fatalf("Error running Migrations")
		}
	}
}
