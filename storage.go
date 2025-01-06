package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
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
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists accounts (
		account_id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		nick_name varchar(50),
		created_at timestamp default current_timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := fmt.Sprintf(`insert into accounts
		(account_id, first_name, last_name, nick_name)
		values (%v, '%v', '%v', '%v')`,
		acc.AccountID, acc.FirstName, acc.LastName, acc.NickName)

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
