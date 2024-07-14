package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Init() (*sql.DB, error){
	connStr := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=%v", 
				os.Getenv("PG_USER"), os.Getenv("PG_DB_NAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_SSL_MODE"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Database initialized")
	return db, nil
}