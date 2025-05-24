package testutils

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

// SetupTestDB creates a test database with sqlmock
func SetupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}
