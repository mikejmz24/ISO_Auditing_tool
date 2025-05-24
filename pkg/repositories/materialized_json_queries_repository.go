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
type MaterializedJSONQueryRepository struct {
	db *sql.DB
}

// Ensure DraftRepository implements DraftRepositoryInterface
var _ MaterializedJSONQueryRepositoryInterface = (*MaterializedJSONQueryRepository)(nil)

func NewMaterializedJSONQueryRepository(db *sql.DB) (MaterializedJSONQueryRepositoryInterface, error) {
	return &MaterializedJSONQueryRepository{db: db}, nil
}

func (r *MaterializedJSONQueryRepository) GetByIDWithFullHierarchyMaterializedJSONQuery(ctx context.Context, materializedJSONQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error) {
	return types.MaterializedJSONQuery{}, nil
}

func (r *MaterializedJSONQueryRepository) GetByNameMaterializedJSONQuery(ctx context.Context, name string) (types.MaterializedJSONQuery, error) {
	query := `
  SELECT 
    id, query_name, query_definition, entity_type, entity_id, 
		data, version,  error_count, last_error, created_at, updated_at
  FROM materialized_json_queries
  WHERE query_name = ?;
  `

	row := r.db.QueryRowContext(ctx, query, name)

	var (
		createdAt, updatedAt  []uint8
		materializedJSONQuery types.MaterializedJSONQuery
	)

	err := row.Scan(
		&materializedJSONQuery.ID,
		&materializedJSONQuery.Name,
		&materializedJSONQuery.Definition,
		&materializedJSONQuery.EntityType,
		&materializedJSONQuery.EntityID,
		&materializedJSONQuery.Data,
		&materializedJSONQuery.Version,
		&materializedJSONQuery.ErrorCount,
		&materializedJSONQuery.LastError,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return types.MaterializedJSONQuery{}, custom_errors.ErrNotFound
	}

	if err != nil {
		return types.MaterializedJSONQuery{}, fmt.Errorf("Failed to scan materialized JSON query: %w", err)
	}

	if materializedJSONQuery.CreatedAt, err = utils.BytesToTime(createdAt); err != nil {
		return types.MaterializedJSONQuery{}, fmt.Errorf("Failed to parse created_at: %w", err)
	}

	if materializedJSONQuery.UpdatedAt, err = utils.BytesToTimePtr(updatedAt); err != nil {
		return types.MaterializedJSONQuery{}, fmt.Errorf("Failed to parse updated_at: %w", err)
	}

	return materializedJSONQuery, nil
}

func (r *MaterializedJSONQueryRepository) GetByEntityTypeAndIDMaterializedJSONQuery(ctx context.Context, entityType string, entityID int) (types.MaterializedJSONQuery, error) {

	return types.MaterializedJSONQuery{}, nil
}

func (r *MaterializedJSONQueryRepository) CreateMaterializedJSONQuery(ctx context.Context, materializedJSONQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error) {
	query := `
  INSERT INTO materialized_json_queries (
    query_name, query_definition, entity_type, entity_id, data, version, error_count, last_error
  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
  `

	result, err := r.db.ExecContext(
		ctx,
		query,
		materializedJSONQuery.Name,
		materializedJSONQuery.Definition,
		materializedJSONQuery.EntityType,
		materializedJSONQuery.EntityID,
		materializedJSONQuery.Data,
		materializedJSONQuery.Version,
		materializedJSONQuery.ErrorCount,
		materializedJSONQuery.LastError,
	)
	if err != nil {
		errRes := fmt.Errorf("Failed to create materialized JSON query: %w", err)

		materializedJSONQuery.ErrorCount = materializedJSONQuery.ErrorCount + 1
		materializedJSONQuery.LastError = errRes.Error()
		return materializedJSONQuery, errRes
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.MaterializedJSONQuery{}, fmt.Errorf("Failed to get last insert ID: %w", err)
	}

	materializedJSONQuery.ID = int(id)
	return materializedJSONQuery, nil
}

func (r *MaterializedJSONQueryRepository) UpdateMaterializedJSONQuery(ctx context.Context, materializedJSONQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error) {
	query := `
	UPDATE materialized_json_queries
	SET 
		query_definition = ?,
		data = ?,
		version = ?,
		error_count = ?,
		last_error = ?
	WHERE query_name = ?;
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		materializedJSONQuery.Definition,
		materializedJSONQuery.Data,
		materializedJSONQuery.Version,
		materializedJSONQuery.ErrorCount,
		materializedJSONQuery.LastError,
		materializedJSONQuery.Name,
	)

	if err != nil {
		errRes := fmt.Errorf("Failed to update materialized JSON query: %w", err)

		materializedJSONQuery.ErrorCount = materializedJSONQuery.ErrorCount + 1
		materializedJSONQuery.LastError = errRes.Error()
		return materializedJSONQuery, errRes
	}

	return materializedJSONQuery, nil
}
