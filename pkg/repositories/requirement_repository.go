package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

func NewRequirementRepository(db *sql.DB) (RequirementRepository, error) {
	return &repository{
		db: db,
	}, nil
}

func (r *repository) GetByIDRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error) {
	return types.Requirement{}, nil
}

func (r *repository) GetByIDWithQuestionsRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error) {
	return types.Requirement{}, nil
}
