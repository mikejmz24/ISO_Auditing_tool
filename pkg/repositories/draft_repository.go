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

func (r *repository) CreateDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
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

func (r *repository) GetByID(ctx context.Context, draft types.Draft) (types.Draft, error) {
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

func (r *repository) UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
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

func (r *repository) Delete(ctx context.Context, draft types.Draft) (types.Draft, error) {
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

// UpdateRequirementAndDeleteDraft atomically updates a requirement and deletes the associated draft
func (r *repository) UpdateRequirementAndDeleteDraft(ctx context.Context, requirement types.Requirement, draft types.Draft) (types.Requirement, error) {
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
