package database

import (
	"ISO_Auditing_Tool/internal/migrations"
	"ISO_Auditing_Tool/internal/seeds"
	"ISO_Auditing_Tool/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

// Config holds MySQL database configuration
type Config struct {
	Username        string
	Password        string
	Host            string
	Port            string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() *Config {
	maxOpenConns := 50 // Default value
	maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	if maxOpenConnsStr != "" {
		if val, err := strconv.Atoi(maxOpenConnsStr); err == nil && val > 0 {
			maxOpenConns = val
		}
	}

	maxIdleConns := 50 // Default value
	maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")
	if maxIdleConnsStr != "" {
		if val, err := strconv.Atoi(maxIdleConnsStr); err == nil && val > 0 {
			maxIdleConns = val
		}
	}

	connMaxLifetime := 0 * time.Second // Default: no maximum lifetime
	connMaxLifetimeStr := os.Getenv("DB_CONN_MAX_LIFETIME_SECONDS")
	if connMaxLifetimeStr != "" {
		if val, err := strconv.Atoi(connMaxLifetimeStr); err == nil && val > 0 {
			connMaxLifetime = time.Duration(val) * time.Second
		}
	}

	return &Config{
		Username:        os.Getenv("DB_USERNAME"),
		Password:        os.Getenv("DB_PASSWORD"),
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		Database:        os.Getenv("DB_DATABASE"),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

// Service defines the interface for database operations
type Service interface {
	Health() map[string]string
	Close() error
	DB() *sql.DB
	Migrate(file string, direction string) error
	Seed() error
	Truncate() error
	RefreshDatabase() error
	Ping() error
}

type service struct {
	db     *sql.DB
	config *Config
}

var dbInstance *service

// New creates a new database service with default configuration
func New() Service {
	return NewWithConfig(LoadConfigFromEnv())
}

// NewWithConfig creates a new database service with the specified configuration
func NewWithConfig(config *Config) Service {
	if dbInstance != nil {
		return dbInstance
	}

	// Validate required configuration
	if config.Database == "" || config.Password == "" || config.Username == "" || config.Port == "" || config.Host == "" {
		log.Fatal("One or more required database configuration parameters are not set")
	}

	// Create DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=30s",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	// Connect to database
	log.Println("Connecting to database...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Configure connection pool
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)

	// Create context with timeout for initial connection test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Connected to database: %s on %s:%s", config.Database, config.Host, config.Port)

	dbInstance = &service{
		db:     db,
		config: config,
	}
	return dbInstance
}

// Health returns the health status of the database
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Check connectivity
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("Database health check failed: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Add connection pool statistics
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Add performance warnings based on metrics
	if dbStats.OpenConnections > s.config.MaxOpenConns*80/100 {
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

// Close closes the database connection
func (s *service) Close() error {
	log.Printf("Disconnecting from database: %s", s.config.Database)
	return s.db.Close()
}

// DB returns the underlying sql.DB instance
func (s *service) DB() *sql.DB {
	return s.db
}

// Ping checks if the database is reachable
func (s *service) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}

// Migrate runs database migrations
func (s *service) Migrate(file string, direction string) error {
	files, err := utils.FindFilesInDir("", file, direction)
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
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

// Seed populates the database with initial data
func (s *service) Seed() error {
	env := os.Getenv("ENV")
	if env == "development" || env == "test" {
		log.Println("Seeding database...")
		if err := seeds.Seed(s.db); err != nil {
			return fmt.Errorf("failed to seed database: %w", err)
		}
		log.Println("Database seeding completed successfully")
	} else {
		log.Println("Skipping database seeding in non-development environment")
	}
	return nil
}

// Truncate removes all data from the database tables
func (s *service) Truncate() error {
	log.Println("Truncating database tables...")
	if err := seeds.Truncate(s.db, seeds.DefaultTruncateOptions()); err != nil {
		return fmt.Errorf("failed to truncate database: %w", err)
	}
	log.Println("Database truncation completed successfully")
	return nil
}

// RefreshDatabase resets the database to a clean state
func (s *service) RefreshDatabase() error {
	log.Println("Refreshing database...")
	if err := seeds.RefreshDatabase(s.db); err != nil {
		return fmt.Errorf("failed to refresh database: %w", err)
	}
	log.Println("Database refresh completed successfully")
	return nil
}
