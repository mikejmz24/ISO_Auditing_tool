package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
	"fmt"
)

func NewStandardRepository(db *sql.DB) (StandardRepository, error) {
	return &repository{
		db: db,
	}, nil
}

func (r *repository) GetByIDStandard(ctx context.Context, standard types.Standard) (types.Standard, error) {
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
