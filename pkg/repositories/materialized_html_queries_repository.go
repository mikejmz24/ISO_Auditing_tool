package repositories

import (
	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/pkg/utils"
	"context"
	"database/sql"
	"fmt"
)

// DraftRepository is the concrete implementation
type MaterializedHTMLQueryRepository struct {
	db *sql.DB
}

// Ensure DraftRepository implements DraftRepositoryInterface
var _ MaterializedHTMLQueryRepositoryInterface = (*MaterializedHTMLQueryRepository)(nil)

func NewMaterializedQueriesHTMLRepository(db *sql.DB) (MaterializedHTMLQueryRepositoryInterface, error) {
	return &MaterializedHTMLQueryRepository{db: db}, nil
}

func (r *MaterializedHTMLQueryRepository) CreateMaterializedHTMLQuery(ctx context.Context, materializedHTMLQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error) {
	query := `
  INSERT INTO materialized_html_queries (
    query_name, view_path, html_content, version, error_count, last_error
  ) VALUES (?, ?, ?, ?, ?, ?);
  `

	result, err := r.db.ExecContext(
		ctx,
		query,
		materializedHTMLQuery.Name,
		materializedHTMLQuery.ViewPath,
		materializedHTMLQuery.HTMLContent,
		materializedHTMLQuery.Version,
		materializedHTMLQuery.ErrorCount,
		materializedHTMLQuery.LastError,
	)
	if err != nil {
		errRes := fmt.Errorf("Failed to create materialized HTML query: %w", err)

		materializedHTMLQuery.ErrorCount = materializedHTMLQuery.ErrorCount + 1
		materializedHTMLQuery.LastError = errRes.Error()
		return materializedHTMLQuery, errRes
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.MaterializedHTMLQuery{}, fmt.Errorf("Failed to get last insert ID: %w", err)
	}

	materializedHTMLQuery.ID = int(id)
	return materializedHTMLQuery, nil
}

func (r *MaterializedHTMLQueryRepository) GetByNameMaterializedHTMLQuery(ctx context.Context, name string) (types.MaterializedHTMLQuery, error) {
	query := `
  SELECT 
    id, query_name, view_path, html_content, version,
    error_count, last_error, created_at, updated_at
  FROM materialized_html_queries
  WHERE query_name = ?;
  `

	row := r.db.QueryRowContext(ctx, query, name)

	var (
		createdAt, updatedAt  []uint8
		materializedHTMLQuery types.MaterializedHTMLQuery
	)

	err := row.Scan(
		&materializedHTMLQuery.ID,
		&materializedHTMLQuery.Name,
		&materializedHTMLQuery.ViewPath,
		&materializedHTMLQuery.HTMLContent,
		&materializedHTMLQuery.Version,
		&materializedHTMLQuery.ErrorCount,
		&materializedHTMLQuery.LastError,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return types.MaterializedHTMLQuery{}, custom_errors.ErrNotFound
	}

	if err != nil {
		return types.MaterializedHTMLQuery{}, fmt.Errorf("Failed to scan materialized HTML query: %w", err)
	}

	if materializedHTMLQuery.CreatedAt, err = utils.BytesToTime(createdAt); err != nil {
		return types.MaterializedHTMLQuery{}, fmt.Errorf("Failed to parse created_at: %w", err)
	}

	if materializedHTMLQuery.UpdatedAt, err = utils.BytesToTimePtr(updatedAt); err != nil {
		return types.MaterializedHTMLQuery{}, fmt.Errorf("Failed to parse updated_at: %w", err)
	}

	return materializedHTMLQuery, nil
}

func (r *MaterializedHTMLQueryRepository) UpdateMaterializedHTMLQuery(ctx context.Context, materializedHTMLQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error) {
	query := `
	UPDATE materialized_html_queries
	SET 
		view_path = ?,
		html_content = ?,
		version = ?,
		error_count = ?,
		last_error = ?
	WHERE query_name = ?;
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		materializedHTMLQuery.ViewPath,
		materializedHTMLQuery.HTMLContent,
		materializedHTMLQuery.Version,
		materializedHTMLQuery.ErrorCount,
		materializedHTMLQuery.LastError,
		materializedHTMLQuery.Name,
	)

	if err != nil {
		errRes := fmt.Errorf("Failed to update materialized HTML query: %w", err)

		materializedHTMLQuery.ErrorCount = materializedHTMLQuery.ErrorCount + 1
		materializedHTMLQuery.LastError = errRes.Error()
		return materializedHTMLQuery, errRes
	}

	return materializedHTMLQuery, nil
}
