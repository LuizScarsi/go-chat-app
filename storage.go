package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	DeleteUser(int) error
	UpdateUser(*User) error
	GetUserByID(int) (*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (p PostgresStore) CreateUser(*User) error {
	fmt.Println("creating user...")
	return nil
}

func (p PostgresStore) DeleteUser(int) error {

	return nil
}

func (p PostgresStore) UpdateUser(*User) error {

	return nil
}

func (p PostgresStore) GetUserByID(int) (*User, error) {

	return nil, nil
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_DB_NAME"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_SSL_MODE"),
	)
	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
