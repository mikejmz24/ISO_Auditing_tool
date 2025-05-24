package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

// EvidenceRepository is the concrete implementation
type EvidenceRepository struct {
	db *sql.DB
}

// Ensure EvidenceRepository implements EvidenceRepositoryInterface
var _ EvidenceRepositoryInterface = (*EvidenceRepository)(nil)

func NewEvidenceRepository(db *sql.DB) (EvidenceRepositoryInterface, error) {
	return &EvidenceRepository{db: db}, nil
}

func (r *EvidenceRepository) GetByIDEvidence(ctx context.Context, evidence types.Evidence) (types.Evidence, error) {
	return types.Evidence{}, nil
}
