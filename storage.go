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

func NewPostgresStore() (*PostgresStore, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_DB_NAME"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_SSL_MODE"),
	)

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

func (s *PostgresStore) Init() error {
	return s.CreateUserTable()
}

func (s *PostgresStore) CreateUserTable() error {
	query := `create table if not exists users (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		nick_name varchar(50)
	)`
	result, err := s.db.Exec(query)

	fmt.Println(result)
	return err
}

func (s *PostgresStore) CreateUser(user *User) error {
	fmt.Println("Creating user...")
	fmt.Println(user.FirstName)
	query := fmt.Sprintf("insert into users (id, first_name, last_name, nick_name) values (%v, '%v', '%v', '%v')", user.ID, user.FirstName, user.LastName, user.NickName)
	result, err := s.db.Exec(query)
	fmt.Println(result)
	return err
}

func (s *PostgresStore) UpdateUser(*User) error {
	return nil
}

func (s *PostgresStore) DeleteUser(id int) error {
	return nil
}

func (s *PostgresStore) GetUserByID(id int) (*User, error) {
	return nil, nil
}
