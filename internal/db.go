package internal

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	return sql.Open("postgres", connStr)
}

func Prepare() error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Print("Retrying to connect")
		time.Sleep(time.Second)
	}
	err = godotenv.Load()
	//if _, err := db.Exec("DROP TABLE IF EXISTS users"); err != nil {
	//	return err
	//}
	//if _, err := db.Exec("DROP TABLE IF EXISTS notes"); err != nil {
	//	return err
	//}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL, username VARCHAR UNIQUE, password VARCHAR)"); err != nil {
		return err
	}
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS notes (
        id SERIAL PRIMARY KEY,
        title VARCHAR,
        content TEXT,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR NOT NULL
    )`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
    INSERT INTO users (username, password)
    VALUES ($1, $2)
    ON CONFLICT (username) DO NOTHING
`, "user1", "password1")
	if err != nil {
		log.Fatalf("Error inserting user1: %v", err)
	}
	_, err = db.Exec(`
    INSERT INTO users (username, password)
    VALUES ($1, $2)
    ON CONFLICT (username) DO NOTHING
`, "user2", "password2")
	if err != nil {
		log.Fatalf("Error inserting user1: %v", err)
	}

	log.Print("Connection Successful")
	return nil
}
