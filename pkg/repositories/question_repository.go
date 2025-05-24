package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

// DraftRepository is the concrete implementation
type QuestionRepository struct {
	db *sql.DB
}

// Ensure DraftRepository implements DraftRepositoryInterface
var _ QuestionRepositoryInterface = (*QuestionRepository)(nil)

func NewQuestionRepository(db *sql.DB) (QuestionRepositoryInterface, error) {
	return &QuestionRepository{db: db}, nil
}

func (r *QuestionRepository) GetByIDQuestion(ctx context.Context, question types.Question) (types.Question, error) {
	return types.Question{}, nil
}

func (r *QuestionRepository) GetByIDWithEvidenceQuestion(ctx context.Context, question types.Question) (types.Question, error) {
	return types.Question{}, nil
}
