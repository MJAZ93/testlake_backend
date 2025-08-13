package main

import (
	"log"

	"testlake/app"
	"testlake/dao"
	_ "testlake/docs"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dao.Connect()

	app.ServeApplication()
}
