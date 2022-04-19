package main

import (
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yanchevdimitar/RSS-Reader-Service/app/database"
	"github.com/yanchevdimitar/RSS-Reader-Service/app/services/queue"
)

func main() {
	currentWorkDirectory, _ := os.Getwd()
	err := godotenv.Load(currentWorkDirectory + "/" + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := DBInit()
	Migrations(db)
	go queue.NewDefaultConsumer(database.NewMySQLRSSRepository(db)).Process()
	go queue.NewDefaultPublisher(database.NewMySQLRSSRepository(db)).Process()

	runtime.Goexit()
}

func DBInit() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: "" + os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(127.0.0.1:" + os.Getenv("MYSQL_PORT") + ")/" +
			os.Getenv("MYSQL_DATABASE") + "?charset=utf8&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error conecting DB")
	}

	return db
}

func Migrations(db *gorm.DB) {
	checkTable := db.Migrator().HasTable("rss")
	if !checkTable {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&database.RSS{})
		if err != nil {
			log.Fatalf("Error running Migrations")
		}
	}
}
