package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "b1a"
)

var db *sql.DB

func init() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to PostgreSQL!")

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS urls (
          id SERIAL PRIMARY KEY,
          url TEXT NOT NULL UNIQUE,
          shortUrl TEXT NOT NULL UNIQUE, 
          author TEXT NOT NULL, 
          sharedWith TEXT NOT NULL,
          topic TEXT NOT NULL,
          message TEXT NOT NUll, 
          clicked INT NOT NULL,
          rating FLOAT NOT NULL
  );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Table 'urls' created successfully!")

	createTableQuery = `
        CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE, 
    password TEXT NOT NULL, 
    lastVisited TEXT[] NOT NULL, 
    sharedWithMe TEXT[] NOT NULL
);
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Table 'users' created successfully!")
}
