package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
	"fmt"
)

// DraftRepository is the concrete implementation
type StandardRepository struct {
	db *sql.DB
}

// Ensure DraftRepository implements DraftRepositoryInterface
var _ StandardRepositoryInterface = (*StandardRepository)(nil)

func NewStandardRepository(db *sql.DB) (StandardRepositoryInterface, error) {
	return &StandardRepository{db: db}, nil
}

func (r *StandardRepository) GetAllStandards(ctx context.Context) ([]types.Standard, error) {
	return []types.Standard{}, nil
}

func (r *StandardRepository) GetByIDWithFullHierarchyStandard(ctx context.Context, standard types.Standard) (types.Standard, error) {
	return types.Standard{}, nil
}

func (r *StandardRepository) GetByIDStandard(ctx context.Context, standard types.Standard) (types.Standard, error) {
	query := `
  INSERT INTO drafts (
    type_id, object_id, status_id, version, data, diff,
    user_id, approver_id, approval_comment, publish_error
  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
  `

	result, err := r.db.ExecContext(
		ctx,
		query,
		standard.ID,
	)
	if err != nil {
		return types.Standard{}, fmt.Errorf("failed to create draft: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.Standard{}, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	standard.ID = int(id)
	return standard, nil
}
