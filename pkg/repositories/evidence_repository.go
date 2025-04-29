package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

func NewEvidenceRepository(db *sql.DB) (EvidenceRepository, error) {
	return &repository{
		db: db,
	}, nil
}

func (r *repository) GetByIDEvidence(ctx context.Context, evidence types.Evidence) (types.Evidence, error) {
	return types.Evidence{}, nil
}
