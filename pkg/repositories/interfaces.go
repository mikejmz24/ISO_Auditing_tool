// Only performs database operations
// Craate files per entity for better maintainability
package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
)

// Repository interface defines the methods for interacting with the database
type DraftRepositoryInterface interface {
	GetAllDrafts(ctx context.Context) ([]types.Draft, error)
	CreateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
	GetDraftByID(ctx context.Context, draft types.Draft) (types.Draft, error)
	UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
	DeleteDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
	GetDraftsByTypeAndObject(ctx context.Context, typeID, objectID int) ([]types.Draft, error)
	UpdateRequirementAndDeleteDraft(ctx context.Context, requirement types.Requirement, draft types.Draft) (types.Requirement, error)
	// Add methods for REST, filtering, searching, etc..
}

type MaterializedJSONQueryRepositoryInterface interface {
	GetByIDWithFullHierarchyMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error)
	GetByNameMaterializedJSONQuery(ctx context.Context, name string) (types.MaterializedJSONQuery, error)
	GetByEntityTypeAndIDMaterializedJSONQuery(ctx context.Context, entityType string, entityID int) (types.MaterializedJSONQuery, error)
	CreateMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error)
	UpdateMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error)
	// Add methods for filtering, searching, etc...
}

type MaterializedHTMLQueryRepositoryInterface interface {
	GetByNameMaterializedHTMLQuery(ctx context.Context, name string) (types.MaterializedHTMLQuery, error)
	CreateMaterializedHTMLQuery(ctx context.Context, materializedQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error)
	UpdateMaterializedHTMLQuery(ctx context.Context, materializedQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error)
	// Add methods for filtering, searching, etc...
}

type StandardRepositoryInterface interface {
	GetAllStandards(ctx context.Context) ([]types.Standard, error)
	GetByIDStandard(ctx context.Context, standard types.Standard) (types.Standard, error)
	GetByIDWithFullHierarchyStandard(ctx context.Context, standard types.Standard) (types.Standard, error)

	// Add methods for filtering, searching, etc...
}

type RequirementRepositoryInterface interface {
	GetByIDRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error)
	GetByIDWithQuestionsRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error)
	UpdateRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error)
	UpdateRequirementAndDeleteDraft(ctx context.Context, requirement types.Requirement, draft types.Draft) (types.Requirement, error)
	// Add methods for filtering, searching, etc...
}

type QuestionRepositoryInterface interface {
	GetByIDQuestion(ctx context.Context, question types.Question) (types.Question, error)
	GetByIDWithEvidenceQuestion(ctx context.Context, question types.Question) (types.Question, error)

	// Add methods for filtering, searching, etc...
}

type EvidenceRepositoryInterface interface {
	GetByIDEvidence(ctx context.Context, evidence types.Evidence) (types.Evidence, error)

	// Add methods for filtering, searching, etc...
}

// // repository struct holds the database connection
// type repository struct {
// 	db *sql.DB
// }
