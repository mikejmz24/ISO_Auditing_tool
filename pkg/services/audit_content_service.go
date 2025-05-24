// pkg/services/audit_content_service.go
package services

import (
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// AuditContentService manages modifications to audit content (requirements, questions, evidence)
type AuditContentService struct {
	DraftService    *DraftService
	RequirementRepo repositories.RequirementRepositoryInterface
	QuestionRepo    repositories.QuestionRepositoryInterface
	EvidenceRepo    repositories.EvidenceRepositoryInterface
	EventBus        *events.EventBus
}

// NewAuditContentService creates a new AuditContentService
func NewAuditContentService(
	draftService *DraftService,
	requirementRepo repositories.RequirementRepositoryInterface,
	questionRepo repositories.QuestionRepositoryInterface,
	evidenceRepo repositories.EvidenceRepositoryInterface,
	eventBus *events.EventBus,
) *AuditContentService {
	return &AuditContentService{
		DraftService:    draftService,
		RequirementRepo: requirementRepo,
		QuestionRepo:    questionRepo,
		EvidenceRepo:    evidenceRepo,
		EventBus:        eventBus,
	}
}

// ModifyRequirementDescription modifies the description of a requirement
func (s *AuditContentService) ModifyRequirementDescription(ctx context.Context, requirementID int, newDescription, reason string, userID int) error {
	// 1. Get the original requirement
	originalReq, err := s.RequirementRepo.GetByIDRequirement(ctx, types.Requirement{ID: requirementID})
	if err != nil {
		return fmt.Errorf("failed to get requirement: %w", err)
	}

	// 2. Create the modification
	modification := types.AuditContentModification{
		ContentType: "requirement",
		ContentID:   requirementID,
		Action:      "modify",
		ModifiedBy:  userID,
		ModifiedAt:  time.Now(),
		Reason:      reason,
	}

	// 3. Prepare the modified requirement
	modifiedReq := types.RequirementModification{
		ID:                  originalReq.ID,
		StandardID:          originalReq.StandardID,
		LevelID:             originalReq.LevelID,
		ParentID:            originalReq.ParentID,
		ReferenceCode:       originalReq.ReferenceCode,
		Name:                originalReq.Name,
		Description:         newDescription,
		OriginalDescription: originalReq.Description,
	}

	// 4. Serialize content
	modifiedContent, err := json.Marshal(modifiedReq)
	if err != nil {
		return fmt.Errorf("failed to marshal modified content: %w", err)
	}

	originalContent, err := json.Marshal(originalReq)
	if err != nil {
		return fmt.Errorf("failed to marshal original content: %w", err)
	}

	modification.ModifiedContent = modifiedContent
	modification.OriginalContent = originalContent

	// 5. Create draft for immediate publishing
	draft := types.Draft{
		TypeID:   11, // audit_content type (from reference_values)
		ObjectID: requirementID,
		StatusID: 12, // draft status
		Version:  1,
		UserID:   userID,
	}

	// 6. Serialize the modification as draft data
	draftData, err := json.Marshal(modification)
	if err != nil {
		return fmt.Errorf("failed to marshal draft data: %w", err)
	}
	draft.Data = draftData

	// 7. Create the draft
	createdDraft, err := s.DraftService.Create(ctx, draft)
	if err != nil {
		return fmt.Errorf("failed to create draft: %w", err)
	}

	// 8. Immediately publish the change (admin has immediate control)
	return s.PublishRequirementChange(ctx, createdDraft)
}

// GetRequirement gets a requirement (always returns current published version)
func (s *AuditContentService) GetRequirement(ctx context.Context, requirementID int) (types.Requirement, error) {
	return s.RequirementRepo.GetByIDRequirement(ctx, types.Requirement{ID: requirementID})
}

// PublishRequirementChange publishes a draft requirement change
func (s *AuditContentService) PublishRequirementChange(ctx context.Context, draft types.Draft) error {
	// 1. Get the draft
	draft, err := s.DraftService.GetByID(ctx, draft)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	// 2. Parse the modification from draft data
	var modification types.AuditContentModification
	if err := json.Unmarshal(draft.Data, &modification); err != nil {
		return fmt.Errorf("failed to parse draft data: %w", err)
	}

	// 3. Parse the modified requirement
	var modifiedReq types.RequirementModification
	if err := json.Unmarshal(modification.ModifiedContent, &modifiedReq); err != nil {
		return fmt.Errorf("failed to parse modified requirement: %w", err)
	}

	// 4. Atomically update requirement and delete draft
	updatedReq := types.Requirement{
		ID:            modifiedReq.ID,
		StandardID:    modifiedReq.StandardID,
		LevelID:       modifiedReq.LevelID,
		ParentID:      modifiedReq.ParentID,
		ReferenceCode: modifiedReq.ReferenceCode,
		Name:          modifiedReq.Name,
		Description:   modifiedReq.Description,
	}

	finalReq, err := s.RequirementRepo.UpdateRequirementAndDeleteDraft(ctx, updatedReq, draft)
	if err != nil {
		return fmt.Errorf("failed to update requirement and delete draft: %w", err)
	}

	// 5. Publish event asynchronously (after successful atomic operation)
	event := events.NewEntityChangeEvent(
		events.EntityRequirement,
		finalReq.ID,
		events.ChangeUpdated,
		"", // no specific query
		events.EntityStandard,
		finalReq.StandardID,
		finalReq,
	)

	s.EventBus.AsyncPublish(ctx, event)

	return nil
}
