// pkg/services/materialized_json_service.go
package services

import (
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type MaterializedJSONService struct {
	JSONRepo         repositories.MaterializedJSONQueryRepository
	StandardRepo     repositories.StandardRepository
	RequirementRepo  repositories.RequirementRepository
	QuestionRepo     repositories.QuestionRepository
	EvidenceRepo     repositories.EvidenceRepository
	EventBus         *events.EventBus
	debounceTimers   map[string]*time.Timer
	debounceInterval time.Duration
	mutex            sync.Mutex
}

// Service calls the MaterializedJSONQuery, Standard, Requirement, Question, Evidence repos and the Event Bus
func NewMaterializedJSONService(
	jsonRepo repositories.MaterializedJSONQueryRepository,
	standardRepo repositories.StandardRepository,
	requirementRepo repositories.RequirementRepository,
	questionRepo repositories.QuestionRepository,
	evidenceRepo repositories.EvidenceRepository,
	eventBus *events.EventBus,
) *MaterializedJSONService {
	service := &MaterializedJSONService{
		JSONRepo:         jsonRepo,
		StandardRepo:     standardRepo,
		RequirementRepo:  requirementRepo,
		QuestionRepo:     questionRepo,
		EvidenceRepo:     evidenceRepo,
		EventBus:         eventBus,
		debounceTimers:   make(map[string]*time.Timer),
		debounceInterval: 2 * time.Second,
	}

	// Subscribe to entity events
	eventBus.Subscribe(events.EntityChanged, service.handleEntityChange)
	// For backward compatibility
	eventBus.Subscribe(events.DataCreated, service.handleLegacyEvent)
	eventBus.Subscribe(events.DataUpdated, service.handleLegacyEvent)
	eventBus.Subscribe(events.DataDeleted, service.handleLegacyEvent)

	return service
}

// HandleEntityChange implements the EntityService interface
func (s *MaterializedJSONService) HandleEntityChange(ctx context.Context, payload events.EntityChangePayload) error {
	// Determine what needs updating based on entity type
	entityType := payload.EntityType
	entityID, ok := payload.EntityID.(int)
	if !ok {
		return fmt.Errorf("expected entity ID to be an integer, got %T", payload.EntityID)
	}

	// Debounce the update to avoid rapid successive updates
	updateKey := fmt.Sprintf("%s_%d", entityType, entityID)
	s.debounceUpdate(updateKey, func() {
		// Use a background context for the debounced function
		bgCtx := context.Background()

		// Update the specific entity
		if err := s.updateEntity(bgCtx, entityType, entityID, payload.Data); err != nil {
			// Log the error
			fmt.Printf("Error updating %s %d: %v\n", entityType, entityID, err)
		}

		// Update any parent entities if needed
		if err := s.updateParentEntities(bgCtx, entityType, entityID, payload.ParentType, payload.ParentID); err != nil {
			fmt.Printf("Error updating parent entities for %s %d: %v\n", entityType, entityID, err)
		}

		// Trigger HTML updates if needed
		if err := s.triggerHTMLUpdate(bgCtx, entityType, entityID); err != nil {
			fmt.Printf("Error triggering HTML update for %s %d: %v\n", entityType, entityID, err)
		}
	})

	return nil
}

// Event handler for the EventBus
func (s *MaterializedJSONService) handleEntityChange(ctx context.Context, event events.Event) error {
	payload, err := events.GetEntityChangePayload(event)
	if err != nil {
		return err
	}

	return s.HandleEntityChange(ctx, payload)
}

// Legacy event handler for backward compatibility
func (s *MaterializedJSONService) handleLegacyEvent(ctx context.Context, event events.Event) error {
	legacyPayload, err := events.GetDataChangePayload(event)
	if err != nil {
		return err
	}

	// Convert to new format
	payload := events.EntityChangePayload{
		EntityType:    events.EntityType(legacyPayload.EntityType),
		EntityID:      legacyPayload.EntityID,
		ChangeType:    events.ChangeType(legacyPayload.ChangeType),
		AffectedQuery: legacyPayload.AffectedQuery,
	}

	return s.HandleEntityChange(ctx, payload)
}

func (s *MaterializedJSONService) debounceUpdate(key string, fn func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Cancel existing timer if any
	if timer, exists := s.debounceTimers[key]; exists {
		timer.Stop()
	}

	// Create new timer
	s.debounceTimers[key] = time.AfterFunc(s.debounceInterval, fn)
}

func (s *MaterializedJSONService) updateEntity(ctx context.Context, entityType events.EntityType, entityID int, data interface{}) error {
	switch entityType {
	case events.EntityStandard:
		return s.updateStandard(ctx, entityID, data)
	case events.EntityRequirement:
		return s.updateRequirement(ctx, entityID, data)
	case events.EntityQuestion:
		return s.updateQuestion(ctx, entityID, data)
	case events.EntityEvidence:
		return s.updateEvidence(ctx, entityID, data)
	default:
		return fmt.Errorf("unknown entity type: %s", entityType)
	}
}

func (s *MaterializedJSONService) updateStandard(ctx context.Context, standardID int, data interface{}) error {
	var standard types.Standard

	// If data is provided, use it directly
	if data != nil {
		if std, ok := data.(types.Standard); ok {
			standard = std
		} else {
			// Try to convert from map
			if dataMap, ok := data.(map[string]interface{}); ok {
				if id, ok := dataMap["id"].(float64); ok {
					standard.ID = int(id)
				}
				// Map other fields as needed
			}
		}
	}

	// If we don't have complete data, fetch it
	if standard.ID == 0 {
		standard.ID = standardID
		fetchedStandard, err := s.StandardRepo.GetByIDStandard(ctx, standard)
		if err != nil {
			return err
		}
		standard = fetchedStandard
	}

	// Convert to JSON
	jsonData, err := json.Marshal(standard)
	if err != nil {
		return err
	}

	// Update or create materialized query
	return s.updateMaterializedQuery(ctx, "standard", standardID, jsonData)
}

func (s *MaterializedJSONService) updateRequirement(ctx context.Context, requirementID int, data interface{}) error {
	var requirement types.Requirement

	// If data is provided, use it directly
	if data != nil {
		if req, ok := data.(types.Requirement); ok {
			requirement = req
		} else {
			// Try to convert from map
			if dataMap, ok := data.(map[string]interface{}); ok {
				if id, ok := dataMap["id"].(float64); ok {
					requirement.ID = int(id)
				}
				// Map other fields as needed
			}
		}
	}

	// If we don't have complete data, fetch it with questions
	if requirement.ID == 0 {
		requirement.ID = requirementID
		var err error
		// Assuming you have a GetByIDWithQuestions method that returns requirement with questions
		requirement, err = s.RequirementRepo.GetByIDWithQuestionsRequirement(ctx, requirement)
		if err != nil {
			return err
		}
	}

	// Convert to JSON
	jsonData, err := json.Marshal(requirement)
	if err != nil {
		return err
	}

	// Update or create materialized query
	return s.updateMaterializedQuery(ctx, "requirement", requirementID, jsonData)
}

func (s *MaterializedJSONService) updateQuestion(ctx context.Context, questionID int, data interface{}) error {
	var question types.Question

	// If data is provided, use it directly
	if data != nil {
		if q, ok := data.(types.Question); ok {
			question = q
		} else {
			// Try to convert from map
			if dataMap, ok := data.(map[string]interface{}); ok {
				if id, ok := dataMap["id"].(float64); ok {
					question.ID = int(id)
				}
				// Map other fields as needed
			}
		}
	}

	// If we don't have complete data, fetch it with evidence
	if question.ID == 0 {
		question.ID = questionID
		var err error
		// Assuming you have a GetByIDWithEvidence method that returns question with evidence
		question, err = s.QuestionRepo.GetByIDWithEvidenceQuestion(ctx, question)
		if err != nil {
			return err
		}
	}

	// Convert to JSON
	jsonData, err := json.Marshal(question)
	if err != nil {
		return err
	}

	// Update or create materialized query
	return s.updateMaterializedQuery(ctx, "question", questionID, jsonData)
}

func (s *MaterializedJSONService) updateEvidence(ctx context.Context, evidenceID int, data interface{}) error {
	var evidence types.Evidence

	// If data is provided, use it directly
	if data != nil {
		if ev, ok := data.(types.Evidence); ok {
			evidence = ev
		} else {
			// Try to convert from map
			if dataMap, ok := data.(map[string]interface{}); ok {
				if id, ok := dataMap["id"].(float64); ok {
					evidence.ID = int(id)
				}
				// Map other fields as needed
			}
		}
	}

	// If we don't have complete data, fetch it
	if evidence.ID == 0 {
		evidence.ID = evidenceID
		var err error
		evidence, err = s.EvidenceRepo.GetByIDEvidence(ctx, evidence)
		if err != nil {
			return err
		}
	}

	// Get parent question to update it
	question := types.Question{ID: evidence.QuestionID}
	return s.updateQuestion(ctx, question.ID, nil)
}

func (s *MaterializedJSONService) updateParentEntities(ctx context.Context, entityType events.EntityType, entityID int, parentType events.EntityType, parentID interface{}) error {
	switch entityType {
	case events.EntityEvidence:
		// If we already know the parent question ID, use it
		if parentType == events.EntityQuestion && parentID != nil {
			questionID, ok := parentID.(int)
			if ok {
				return s.updateQuestion(ctx, questionID, nil)
			}
		}

		// Otherwise fetch the evidence to find its question
		evidence := types.Evidence{ID: entityID}
		fetchedEvidence, err := s.EvidenceRepo.GetByIDEvidence(ctx, evidence)
		if err != nil {
			return err
		}
		return s.updateQuestion(ctx, fetchedEvidence.QuestionID, nil)

	case events.EntityQuestion:
		// If we already know the parent requirement ID, use it
		if parentType == events.EntityRequirement && parentID != nil {
			requirementID, ok := parentID.(int)
			if ok {
				return s.updateRequirement(ctx, requirementID, nil)
			}
		}

		// Otherwise fetch the question to find its requirement
		question := types.Question{ID: entityID}
		fetchedQuestion, err := s.QuestionRepo.GetByIDQuestion(ctx, question)
		if err != nil {
			return err
		}
		return s.updateRequirement(ctx, fetchedQuestion.RequirementID, nil)

	case events.EntityRequirement:
		// If we already know the parent standard ID, use it
		if parentType == events.EntityStandard && parentID != nil {
			standardID, ok := parentID.(int)
			if ok {
				return s.updateStandardFull(ctx, standardID)
			}
		}

		// Otherwise fetch the requirement to find its standard
		requirement := types.Requirement{ID: entityID}
		fetchedRequirement, err := s.RequirementRepo.GetByIDRequirement(ctx, requirement)
		if err != nil {
			return err
		}
		return s.updateStandardFull(ctx, fetchedRequirement.StandardID)

	case events.EntityStandard:
		// Update full standard hierarchy
		return s.updateStandardFull(ctx, entityID)
	}

	return nil
}

func (s *MaterializedJSONService) updateStandardFull(ctx context.Context, standardID int) error {
	// This builds the complete hierarchy for a standard
	standard := types.Standard{ID: standardID}
	fetchedStandard, err := s.fetchStandardWithFullHierarchy(ctx, standard)
	if err != nil {
		return err
	}

	// Convert to JSON
	jsonData, err := json.Marshal(fetchedStandard)
	if err != nil {
		return err
	}

	// Update or create materialized query
	return s.updateMaterializedQuery(ctx, "standard_full", standardID, jsonData)
}

func (s *MaterializedJSONService) fetchStandardWithFullHierarchy(ctx context.Context, standard types.Standard) (types.Standard, error) {
	// Fetch standard with its complete hierarchy (requirements, questions, evidence)
	// This would call your repository method that fetches the complete hierarchy
	return s.StandardRepo.GetByIDWithFullHierarchyStandard(ctx, standard)
}

func (s *MaterializedJSONService) updateMaterializedQuery(ctx context.Context, entityType string, entityID int, jsonData json.RawMessage) error {
	queryName := fmt.Sprintf("%s_%d", entityType, entityID)

	// Create the materialized query object
	materializedQuery := types.MaterializedJSONQuery{
		Name:       queryName,
		EntityType: entityType,
		EntityID:   entityID,
		Data:       jsonData,
	}

	// Try to get existing query
	existingQuery, err := s.JSONRepo.GetByNameMaterializedJSONQuery(ctx, materializedQuery.Name)
	if err == nil {
		// Update existing
		materializedQuery.ID = existingQuery.ID
		materializedQuery.Version = existingQuery.Version + 1
		_, err = s.JSONRepo.UpdateMaterializedJSONQuery(ctx, materializedQuery)
		return err
	}

	// Create new
	materializedQuery.Version = 1
	_, err = s.JSONRepo.CreateMaterializedJSONQuery(ctx, materializedQuery)
	return err
}

func (s *MaterializedJSONService) triggerHTMLUpdate(ctx context.Context, entityType events.EntityType, entityID int) error {
	// Only trigger HTML updates for standard-level changes or when we've updated a standard_full
	if entityType == events.EntityStandard {
		event := events.NewMaterializedQueryUpdatedEvent(
			fmt.Sprintf("standard_%d", entityID),
			"",  // query definition not needed here
			nil, // data not needed here
			0,   // version not needed here
			0,   // error count not needed here
			"",  // last error not needed here
		)
		return s.EventBus.Publish(ctx, event)
	}

	// Also publish an event for standard_full if appropriate
	if entityType == events.EntityStandard || entityType == events.EntityRequirement {
		// For requirements, we need to find the standard ID
		standardID := entityID
		if entityType == events.EntityRequirement {
			requirement := types.Requirement{ID: entityID}
			fetchedRequirement, err := s.RequirementRepo.GetByIDRequirement(ctx, requirement)
			if err != nil {
				return err
			}
			standardID = fetchedRequirement.StandardID
		}

		event := events.NewMaterializedQueryUpdatedEvent(
			fmt.Sprintf("standard_full_%d", standardID),
			"",  // query definition not needed here
			nil, // data not needed here
			0,   // version not needed here
			0,   // error count not needed here
			"",  // last error not needed here
		)
		return s.EventBus.Publish(ctx, event)
	}

	return nil
}
