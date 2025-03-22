package repositories

import (
	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func NewMaterializedQueriesRepository(db *sql.DB) (MaterializedQueryRepository, error) {
	return &repository{
		db: db,
	}, nil
}

func (r *repository) CreateMaterializedQuery(ctx context.Context, materializedQuery types.MaterializedQuery) (types.MaterializedQuery, error) {
	query := `
  INSERT INTO materialized_queries (
    query_name, query_definition, data, version, error_count, last_error
  ) VALUES (?, ?, ?, ?, ?, ?);
  `

	result, err := r.db.ExecContext(
		ctx,
		query,
		materializedQuery.Name,
		materializedQuery.Definition,
		materializedQuery.Data,
		materializedQuery.Version,
		materializedQuery.ErrorCount,
		materializedQuery.LastError,
	)
	if err != nil {
		errRes := fmt.Errorf("Failed to create materialized query: %w", err)

		materializedQuery.ErrorCount = materializedQuery.ErrorCount + 1
		materializedQuery.LastError = errRes.Error()
		return materializedQuery, errRes
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.MaterializedQuery{}, fmt.Errorf("Failed to get last insert ID: %w", err)
	}

	materializedQuery.ID = int(id)
	return materializedQuery, nil
}

func (r *repository) GetByNameMaterializedQuery(ctx context.Context, name string) (types.MaterializedQuery, error) {
	query := `
  SELECT 
    id, query_name, query_definition, data, version,
    error_count, last_error, created_at, updated_at
  FROM materialized_queries
  WHERE query_name = ?;
  `

	row := r.db.QueryRowContext(ctx, query, name)

	var (
		createdAt, updatedAt []uint8
		materializedQuery    types.MaterializedQuery
	)

	err := row.Scan(
		&materializedQuery.ID,
		&materializedQuery.Name,
		&materializedQuery.Definition,
		&materializedQuery.Data,
		&materializedQuery.Version,
		&materializedQuery.ErrorCount,
		&materializedQuery.LastError,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return types.MaterializedQuery{}, custom_errors.ErrNotFound
	}

	if err != nil {
		return types.MaterializedQuery{}, fmt.Errorf("Failed to scan materialized query: %w", err)
	}

	if materializedQuery.CreatedAt, err = bytesToTime(createdAt); err != nil {
		return types.MaterializedQuery{}, fmt.Errorf("Failed to parse created_at: %w", err)
	}

	if materializedQuery.UpdatedAt, err = bytesToTimePtr(updatedAt); err != nil {
		return types.MaterializedQuery{}, fmt.Errorf("Failed to parse updated_at: %w", err)
	}

	return materializedQuery, nil
}

func (r *repository) UpdateMaterializedQuery(ctx context.Context, materializedQuery types.MaterializedQuery) (types.MaterializedQuery, error) {
	query := `
	UPDATE materialized_queries
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
		materializedQuery.Definition,
		materializedQuery.Data,
		materializedQuery.Version,
		materializedQuery.ErrorCount,
		materializedQuery.LastError,
		materializedQuery.Name,
	)

	if err != nil {
		errRes := fmt.Errorf("Failed to update materialized query: %w", err)

		materializedQuery.ErrorCount = materializedQuery.ErrorCount + 1
		materializedQuery.LastError = errRes.Error()
		return materializedQuery, errRes
	}

	return materializedQuery, nil
}

func bytesToTime(b []uint8) (time.Time, error) {
	if b == nil {
		return time.Time{}, fmt.Errorf("Nil vallue cannot be converted to time.Time")
	}
	str := string(b)
	layout := "2006-01-02 15:04:05"

	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func bytesToTimePtr(b []uint8) (*time.Time, error) {
	if b == nil {
		return nil, nil
	}
	str := string(b)
	layout := "2006-01-02 15:05:05"
	t, err := time.Parse(layout, str)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
