package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
	"fmt"
)

func NewDraftRepository(db *sql.DB) (DraftRepository, error) {
	return &repository{
		db: db,
	}, nil
}

func (r *repository) Create(ctx context.Context, draft types.Draft) (types.Draft, error) {
	query := ""

	result, err := r.db.ExecContext(ctx, query, draft)
	if err != nil {
		return types.Draft{}, fmt.Errorf("failed to create draft: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.Draft{}, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	draft.ID = int(id)
	return draft, nil
}
