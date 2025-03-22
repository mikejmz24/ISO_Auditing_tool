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
	CreateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
	UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)

	// Add methods for REST, filtering, searching, etc..
}

type StandardRepository interface {
	GetByIDStandard(ctx context.Context, standard types.Standard) (types.Standard, error)

	// Add methods for filtering, searching, etc...
}

type MaterializedQueryRepository interface {
	CreateMaterializedQuery(ctx context.Context, MaterializedQuery types.MaterializedQuery) (types.MaterializedQuery, error)
	GetByNameMaterializedQuery(ctx context.Context, name string) (types.MaterializedQuery, error)
	UpdateMaterializedQuery(ctx context.Context, materializedQuery types.MaterializedQuery) (types.MaterializedQuery, error)
	// Add methods for filtering, searching, etc...
}

// repository struct holds the database connection
type repository struct {
	db *sql.DB
}
