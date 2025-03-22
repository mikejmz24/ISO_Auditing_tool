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
