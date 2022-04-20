package main

import (
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	"github.com/yanchevdimitar/RSS-Reader-Service/app/database"
	"github.com/yanchevdimitar/RSS-Reader-Service/app/services/queue"
)

func main() {
	currentWorkDirectory, _ := os.Getwd()
	err := godotenv.Load(currentWorkDirectory + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := database.DBInit()
	database.Migrations(db)
	go queue.NewDefaultConsumer(database.NewMySQLRSSRepository(db)).Process()
	go queue.NewDefaultPublisher(database.NewMySQLRSSRepository(db)).Process()

	runtime.Goexit()
}
