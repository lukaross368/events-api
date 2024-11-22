package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/api.db"
	}

	DB, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	log.Println("Database connection established")

	createTables()
}

func createTables() {

	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`

	if _, err := DB.Exec(createUserTable); err != nil {
		log.Fatalf("Could not create user table: %v", err)
	} else {
		log.Println("User table created or already exists")
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	if _, err := DB.Exec(createEventsTable); err != nil {
		log.Fatalf("Could not create events table: %v", err)
	} else {
		log.Println("Events table created or already exists")
	}

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY (event_id) REFERENCES events(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`
	if _, err := DB.Exec(createRegistrationsTable); err != nil {
		log.Fatalf("Could not create registrations table: %v", err)
	} else {
		log.Println("Registrations table created or already exists")
	}
}
