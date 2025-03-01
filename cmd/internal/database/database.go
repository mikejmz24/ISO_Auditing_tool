package database

import (
	"ISO_Auditing_Tool/cmd/internal/migrations"
	"ISO_Auditing_Tool/cmd/internal/seeds"
	"ISO_Auditing_Tool/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"path/filepath"
)

type Service interface {
	Health() map[string]string
	Close() error
	DB() *sql.DB
	Migrate(file string, direction string) error
	Seed() error
	Truncate() error
	RefreshDatabase() error
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	log.Println("Connecting to database...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Connected to database: %s", dbname)
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

func (s *service) Migrate(file string, direction string) error {
	files, err := utils.FindFilesInDir(file, direction)
	if err != nil {
		return nil
	}

	if s.db != nil {
		log.Printf("Running %s migrations...", direction)
		for _, sqlFile := range files {
			log.Printf("Executing migration: %s", filepath.Base(sqlFile))
			if err := migrations.Migrate(s.db, sqlFile); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", filepath.Base(sqlFile), err)
			}
		}
	}
	return nil
}

func (s *service) Seed() error {
	env := os.Getenv("ENV")
	if env == "development" || env == "test" {
		log.Println("Seeding database...")
		if err := seeds.Seed(s.db); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("Database seeding completed successfully")
	} else {
		log.Println("Skipping database seeding in non-development environment")
	}
	return nil
}

func (s *service) Truncate() error {
	if err := seeds.Truncate(s.db, seeds.DefaultTruncateOptions()); err != nil {
		log.Fatalf("Failed to truncate database: %v", err)
	}
	return nil
}

func (s *service) RefreshDatabase() error {
	if err := seeds.RefreshDatabase(s.db); err != nil {
		log.Fatalf("Failed to refresh database: %v", err)
	}
	return nil
}
