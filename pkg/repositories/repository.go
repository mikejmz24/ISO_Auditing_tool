// Only performs database operations
// Craate files per entity for better maintainability
package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

// Repository interface defines the methods for interacting with the database
type DraftRepository interface {
	CreateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
	GetByID(ctx context.Context, draft types.Draft) (types.Draft, error)
	UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
	Delete(ctx context.Context, draft types.Draft) (types.Draft, error)

	// Add methods for REST, filtering, searching, etc..
}

type MaterializedJSONQueryRepository interface {
	GetByIDWithFullHierarchyMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error)
	GetByNameMaterializedJSONQuery(ctx context.Context, name string) (types.MaterializedJSONQuery, error)
	GetByEntityTypeAndIDMaterializedJSONQuery(ctx context.Context, entityType string, entityID int) (types.MaterializedJSONQuery, error)
	CreateMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error)
	UpdateMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error)
	// Add methods for filtering, searching, etc...
}

type MaterializedHTMLQueryRepository interface {
	GetByNameMaterializedHTMLQuery(ctx context.Context, name string) (types.MaterializedHTMLQuery, error)
	CreateMaterializedHTMLQuery(ctx context.Context, materializedQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error)
	UpdateMaterializedHTMLQuery(ctx context.Context, materializedQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error)
	// Add methods for filtering, searching, etc...
}

type StandardRepository interface {
	GetAllStandards(ctx context.Context) ([]types.Standard, error)
	GetByIDStandard(ctx context.Context, standard types.Standard) (types.Standard, error)
	GetByIDWithFullHierarchyStandard(ctx context.Context, standard types.Standard) (types.Standard, error)

	// Add methods for filtering, searching, etc...
}

type RequirementRepository interface {
	GetByIDRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error)
	GetByIDWithQuestionsRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error)
	UpdateRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error)
	UpdateRequirementAndDeleteDraft(ctx context.Context, requirement types.Requirement, draft types.Draft) (types.Requirement, error)
	// Add methods for filtering, searching, etc...
}

type QuestionRepository interface {
	GetByIDQuestion(ctx context.Context, question types.Question) (types.Question, error)
	GetByIDWithEvidenceQuestion(ctx context.Context, question types.Question) (types.Question, error)

	// Add methods for filtering, searching, etc...
}

type EvidenceRepository interface {
	GetByIDEvidence(ctx context.Context, evidence types.Evidence) (types.Evidence, error)

	// Add methods for filtering, searching, etc...
}

// repository struct holds the database connection
type repository struct {
	db *sql.DB
}
