package sqltools

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func (db *RealDB) SetupDatabase() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:secret@localhost:5432/notes?sslmode=disable"
	}
	dbpointer, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
		return err
	}

	if err = dbpointer.Ping(); err != nil {
		log.Fatal("Failed to connect to the database: ", err)
		return err
	}

	db.DB = dbpointer

	CreateUserTable(dbpointer)
	CreateLoggedInUserTable(dbpointer)
	CreateForgotPasswordTable(dbpointer)
	CreateStatementsTable(dbpointer)

	return nil
}

func CreateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(50) NOT NULL UNIQUE CHECK (user_id <> ''),
		name VARCHAR(150) NOT NULL CHECK (name <> ''),
		email VARCHAR(100) NOT NULL UNIQUE CHECK (email <> ''),
		password VARCHAR(100) NOT NULL CHECK (password <> ''),
		code VARCHAR(100),
		verified BOOL NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Could not create users table: ", err)
		return err
	}
	return nil
}

func CreateLoggedInUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS loggedin_users (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(50) NOT NULL UNIQUE CHECK (user_id <> ''),
		code VARCHAR(200),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Could not create loggedin_users table: ", err)
		return err
	}
	return nil
}

func CreateForgotPasswordTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS forgotpassword (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(50) NOT NULL UNIQUE CHECK (user_id <> ''),
		code VARCHAR(200),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Could not create forgotpassword table: ", err)
		return err
	}
	return nil
}

func CreateStatementsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS statements (
		id SERIAL PRIMARY KEY,
		statement_id VARCHAR(50) NOT NULL UNIQUE CHECK (statement_id <> ''),
		user_id VARCHAR(50) NOT NULL CHECK (user_id <> ''),
		content TEXT NOT NULL CHECK (content <> ''),
		association TEXT NOT NULL CHECK (association <> ''),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Could not create statements table: ", err)
		return err
	}
	return nil
}

func DeleteUserTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS users;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Could not delete table: ", err)
		return err
	}
	return nil
}

func AlterUsersTable(db *sql.DB) error {
	query := `
		ALTER TABLE users 
		ADD COLUMN IF NOT EXISTS rank INT NOT NULL DEFAULT 0 CHECK (rank >= 0);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Could not alter table: ", err)
		return err
	}
	return nil
}
