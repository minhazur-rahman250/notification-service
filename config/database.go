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

	createCertificateTableQuery := `
	CREATE TABLE IF NOT EXISTS certificates (
		id SERIAL PRIMARY KEY,
		certificate_number VARCHAR(100) UNIQUE NOT NULL,
		student_email VARCHAR(255) NOT NULL,
		student_name VARCHAR(255) NOT NULL,
		course_title VARCHAR(255) NOT NULL,
		file_path VARCHAR(500) NOT NULL,
		issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createCertificateTableQuery); err != nil {
		log.Fatal("certificates টেবিল তৈরি করতে ব্যর্থ:", err)
	}

	log.Println("certificates টেবিল প্রস্তুত")

	createCourseIndexTableQuery := `
	CREATE TABLE IF NOT EXISTS course_search_index (
		id SERIAL PRIMARY KEY,
		course_id INTEGER UNIQUE NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10,2) DEFAULT 0,
		teacher_name VARCHAR(255),
		indexed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createCourseIndexTableQuery); err != nil {
		log.Fatal("course_search_index টেবিল তৈরি করতে ব্যর্থ:", err)
	}

	log.Println("course_search_index টেবিল প্রস্তুত")
}