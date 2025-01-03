package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", store)

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(os.Getenv("SERVER_PORT"), store)
	server.Run()
}
