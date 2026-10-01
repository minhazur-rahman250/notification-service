package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDatabase - NestJS এর মতোই একই PostgreSQL ডাটাবেজে কানেক্ট করে,
// কিন্তু নিজস্ব টেবিল ব্যবহার করে (notifications)
func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("ডাটাবেজ কানেকশন তৈরি করতে ব্যর্থ:", err)
	}

	// আসল কানেকশন যাচাই করা (sql.Open শুধু object বানায়, connect করে না)
	if err = db.Ping(); err != nil {
		log.Fatal("ডাটাবেজে পিং করতে ব্যর্থ:", err)
	}

	// টেবিল না থাকলে তৈরি করে নেয় (TypeORM এর synchronize এর মতো, কিন্তু manual)
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		recipient_email VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		message TEXT NOT NULL,
		is_read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal("notifications টেবিল তৈরি করতে ব্যর্থ:", err)
	}

	DB = db
	log.Println("PostgreSQL এর সাথে কানেকশন সফল, notifications টেবিল প্রস্তুত")
}