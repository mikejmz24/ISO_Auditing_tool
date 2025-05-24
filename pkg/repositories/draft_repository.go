package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
	"fmt"

	"github.com/goccy/go-json"
)

// DraftRepository is the concrete implementation
type DraftRepository struct {
	db *sql.DB
}

// Ensure DraftRepository implements DraftRepositoryInterface
var _ DraftRepositoryInterface = (*DraftRepository)(nil)

func NewDraftRepository(db *sql.DB) (DraftRepositoryInterface, error) {
	return &DraftRepository{db: db}, nil
}

func (r *DraftRepository) GetAllDrafts(ctx context.Context) ([]types.Draft, error) {
	query := `
	SELECT id, type_id, object_id, status_id, version, data, diff,
				user_id, approver_id, approval_comment, publish_error, 
				created_at, updated_at, expires_at
	FROM drafts
	ORDER BY created_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query drafts: %w", err)
	}
	defer rows.Close()

	var drafts []types.Draft
	for rows.Next() {
		var draft types.Draft
		var updatedAt, expiresAt sql.NullTime
		var data, diff sql.NullString

		err := rows.Scan(
			&draft.ID,
			&draft.TypeID,
			&draft.ObjectID,
			&draft.StatusID,
			&draft.Version,
			&data,
			&diff,
			&draft.UserID,
			&draft.ApproverID,
			&draft.ApprovalComment,
			&draft.PublishError,
			&draft.CreatedAt,
			&updatedAt,
			&expiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan draft row: %w", err)
		}

		// Handle nullable data field
		if data.Valid {
			draft.Data = json.RawMessage(data.String)
		} else {
			draft.Data = nil
		}

		// Handle nullable diff field
		if diff.Valid {
			draft.Diff = json.RawMessage(data.String)
		} else {
			draft.Diff = nil
		}

		// Handle nullable updated_at field
		if updatedAt.Valid {
			draft.UpdatedAt = updatedAt.Time
		}

		// Handle nullable expires_at field
		if expiresAt.Valid {
			draft.ExpiresAt = expiresAt.Time
		}

		drafts = append(drafts, draft)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over draft rows: %w", err)
	}

	return drafts, nil
}

func (r *DraftRepository) CreateDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
	query := `
  INSERT INTO drafts (
    type_id, object_id, status_id, version, data, diff,
    user_id, approver_id, approval_comment, publish_error
  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
  `

	result, err := r.db.ExecContext(
		ctx,
		query,
		draft.TypeID,
		draft.ObjectID,
		draft.StatusID,
		draft.Version,
		draft.Data,
		draft.Diff,
		draft.UserID,
		draft.ApproverID,
		draft.ApprovalComment,
		draft.PublishError,
	)
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

func (r *DraftRepository) GetDraftByID(ctx context.Context, draft types.Draft) (types.Draft, error) {
	query := ``
	result, err := r.db.ExecContext(
		ctx,
		query,
		draft.ID,
	)
	if err != nil {
		return draft, fmt.Errorf("Failed to get draft by ID: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return draft, fmt.Errorf("Failed to get rows affected: %d, %w", rows, err)
	}

	return draft, nil
}

func (r *DraftRepository) UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
	query := `
  UPDATE drafts
  SET  data = ?
  WHERE id = ?;
  `

	result, err := r.db.ExecContext(
		ctx,
		query,
		draft.Data,
		draft.ID,
	)
	if err != nil {
		return draft, fmt.Errorf("Failed to update draft: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return draft, fmt.Errorf("Failed to get rows affected: %d, %w", rows, err)
	}

	return draft, nil
}

func (r *DraftRepository) DeleteDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
	query := ``
	result, err := r.db.ExecContext(
		ctx,
		query,
		draft.ID,
	)
	if err != nil {
		return draft, fmt.Errorf("Failed to delete draft by ID: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return draft, fmt.Errorf("Failed to get rows affected: %d, %w", rows, err)
	}

	return draft, nil
}

func (r *DraftRepository) GetDraftsByTypeAndObject(ctx context.Context, typeID, objectID int) ([]types.Draft, error) {
	query := ``
	rows, err := r.db.QueryContext(ctx, query, typeID, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query drafts: %w", err)
	}
	defer rows.Close()

	var drafts []types.Draft
	for rows.Next() {
		var draft types.Draft
		var updatedAt sql.NullTime

		err := rows.Scan(
			&draft.ID,
			&draft.TypeID,
			&draft.ObjectID,
			&draft.StatusID,
			&draft.Version,
			&draft.Data,
			&draft.Diff,
			&draft.UserID,
			&draft.ApproverID,
			&draft.ApprovalComment,
			&draft.PublishError,
			&draft.CreatedAt,
			&updatedAt,
			&draft.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan draft row: %w", err)
		}

		// Handle nullable updated_at field
		if updatedAt.Valid {
			draft.UpdatedAt = updatedAt.Time
		}

		drafts = append(drafts, draft)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over draft rows: %w", err)
	}

	return drafts, nil
}

// UpdateRequirementAndDeleteDraft atomically updates a requirement and deletes the associated draft
func (r *DraftRepository) UpdateRequirementAndDeleteDraft(ctx context.Context, requirement types.Requirement, draft types.Draft) (types.Requirement, error) {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return types.Requirement{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // Safe to call even after commit

	// 1. Update requirement in transaction
	updateQuery := `
		UPDATE requirement 
		SET standard_id = ?, requirement_level_id = ?, parent_id = ?, 
		    reference_code = ?, name = ?, description = ?
		WHERE id = ?
	`

	_, err = tx.ExecContext(
		ctx,
		updateQuery,
		requirement.StandardID,
		requirement.LevelID,
		requirement.ParentID,
		requirement.ReferenceCode,
		requirement.Name,
		requirement.Description,
		requirement.ID,
	)
	if err != nil {
		return types.Requirement{}, fmt.Errorf("failed to update requirement: %w", err)
	}

	// 2. Delete draft in transaction
	deleteQuery := `DELETE FROM drafts WHERE id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, draft)
	if err != nil {
		return types.Requirement{}, fmt.Errorf("failed to delete draft: %w", err)
	}

	// 3. Commit transaction
	if err := tx.Commit(); err != nil {
		return types.Requirement{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return requirement, nil
}
