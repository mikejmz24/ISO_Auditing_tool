package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

func NewQuestionRepository(db *sql.DB) (QuestionRepository, error) {
	return &repository{
		db: db,
	}, nil
}

func (r *repository) GetByIDQuestion(ctx context.Context, question types.Question) (types.Question, error) {
	return types.Question{}, nil
}

func (r *repository) GetByIDWithEvidenceQuestion(ctx context.Context, question types.Question) (types.Question, error) {
	return types.Question{}, nil
}
