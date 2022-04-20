package database

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBInit() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: "" + os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" +
			os.Getenv("MYSQL_DATABASE") + "?charset=utf8&parseTime=True&loc=Local",
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error conecting DB")
	}

	return db
}
