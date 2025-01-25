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
	GetAccountByEmail(string) (*Account, error)
	GetAccounts() ([]*Account, error)
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
		email varchar(50),
		first_name varchar(50),
		last_name varchar(50),
		nick_name varchar(50),
		password_hash varchar(100),
		created_at timestamp default current_timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into accounts
		(email, first_name, last_name, nick_name, password_hash)
		values ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(query,
		acc.Email,
		acc.FirstName,
		acc.LastName,
		acc.NickName,
		acc.PasswordHash,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := "delete from accounts where account_id = $1"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	query := "select * from accounts where email = $1"
	rows, err := s.db.Query(query, email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account with email %s not found", email)
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := "select * from accounts where account_id = $1"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccounts(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccounts(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.AccountID,
		&account.Email,
		&account.FirstName,
		&account.LastName,
		&account.NickName,
		&account.PasswordHash,
		&account.CreatedAt,
	)

	return account, err
}
