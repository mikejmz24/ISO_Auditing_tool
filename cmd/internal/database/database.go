package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	Close() error
	DB() *sql.DB
	Migrate() error
}

type service struct {
	db *sql.DB
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	if dbname == "" || password == "" || username == "" || port == "" || host == "" {
		log.Fatal("One or more required environment variables (DB_DATABASE, DB_PASSWORD, DB_USERNAME, DB_PORT, DB_HOST) are not set")
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	return s.db.Close()
}

func (s *service) DB() *sql.DB {
	return s.db
}

func (s *service) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS iso_standard (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS clause (
			id INT AUTO_INCREMENT PRIMARY KEY,
			iso_standard_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			FOREIGN KEY (iso_standard_id) REFERENCES iso_standard(id)
		)`,
		`CREATE TABLE IF NOT EXISTS section (
			id INT AUTO_INCREMENT PRIMARY KEY,
			clause_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			FOREIGN KEY (clause_id) REFERENCES clause(id)
		)`,
		`CREATE TABLE IF NOT EXISTS question (
			id INT AUTO_INCREMENT PRIMARY KEY,
			section_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			FOREIGN KEY (section_id) REFERENCES section(id)
		)`,
		`CREATE TABLE IF NOT EXISTS evidence (
			id INT AUTO_INCREMENT PRIMARY KEY,
			expected VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS evidence_provided (
			id INT AUTO_INCREMENT PRIMARY KEY,
			evidence_id INT NOT NULL,
			provided TEXT NOT NULL,
			FOREIGN KEY (evidence_id) REFERENCES evidence(id)
		)`,
		`CREATE TABLE IF NOT EXISTS comment (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			text TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(id)
		)`,
		`CREATE TABLE IF NOT EXISTS audit (
			id INT AUTO_INCREMENT PRIMARY KEY,
			datetime DATETIME NOT NULL,
			iso_standard_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			team VARCHAR(255) NOT NULL,
			user_id INT NOT NULL,
			FOREIGN KEY (iso_standard_id) REFERENCES iso_standard(id),
			FOREIGN KEY (user_id) REFERENCES user(id)
		)`,
		`CREATE TABLE IF NOT EXISTS audit_questions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			audit_id INT NOT NULL,
			evidence_provided_id INT NOT NULL,
			question_id INT NOT NULL,
			FOREIGN KEY (audit_id) REFERENCES audit(id),
			FOREIGN KEY (evidence_provided_id) REFERENCES evidence_provided(id),
			FOREIGN KEY (question_id) REFERENCES question(id)
		)`,
		`CREATE TABLE IF NOT EXISTS audit_question_comments (
			id INT AUTO_INCREMENT PRIMARY KEY,
			audit_question_id INT NOT NULL,
			comment_id INT NOT NULL,
			FOREIGN KEY (audit_question_id) REFERENCES audit_questions(id),
			FOREIGN KEY (comment_id) REFERENCES comment(id)
		)`,
	}
	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	log.Println("Database tables created successfully")
	return nil
}
