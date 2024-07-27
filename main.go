package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := DatabaseInit()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := NewAPIServer(os.Getenv("SERVER_PORT"))
	server.Run()
}