package main

import (
	"fmt"
	"go-chat-app/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}