package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// db, err := DatabaseInit()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// var storage Storage
	// storage = PostgresStore{db}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", store)

	// server := NewAPIServer(os.Getenv("SERVER_PORT"), store)
	// server.Run()
}
