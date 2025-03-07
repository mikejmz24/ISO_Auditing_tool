// Only performs database operations
// Craate files per entity for better maintainability
package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

// Repository interface defines the methods for interacting with the database
type DraftRepository interface {
	Create(ctx context.Context, draft types.Draft) (types.Draft, error)

	// Add methods for REST, filtering, searching, etc..
}

// repository struct holds the database connection
type repository struct {
	db *sql.DB
}
