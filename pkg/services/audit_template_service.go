// pkg/services/audit_template_service.go
package services

// import (
// 	"ISO_Auditing_Tool/pkg/events"
// 	// "ISO_Auditing_Tool/pkg/templates"
// 	"ISO_Auditing_Tool/pkg/types"
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"sync"
// 	"time"
// )
//
// const (
// 	// Type IDs for drafts
// 	TEMPLATE_TYPE_ID = 1
//
// 	// Status IDs for drafts
// 	DRAFT_STATUS_ID     = 1
// 	PUBLISHED_STATUS_ID = 2
// )
//
// // Repositories interfaces
// type DraftRepository interface {
// 	CreateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
// 	GetDraftByID(ctx context.Context, id int) (types.Draft, error)
// 	UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error)
// 	GetDraftsByTypeAndObject(ctx context.Context, typeID, objectID int) ([]types.Draft, error)
// }
//
// type HTMLRepository interface {
// 	CreateMaterializedHTMLQuery(ctx context.Context, htmlQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error)
// 	GetByNameMaterializedHTMLQuery(ctx context.Context, name string) (types.MaterializedHTMLQuery, error)
// 	UpdateMaterializedHTMLQuery(ctx context.Context, htmlQuery types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error)
// }
//
// type StandardRepository interface {
// 	GetByIDStandard(ctx context.Context, standard types.Standard) (types.Standard, error)
// 	GetAllStandards(ctx context.Context) ([]types.Standard, error)
// }
//
// // AuditTemplateService provides functionality for audit templates
// type AuditTemplateService struct {
// 	draftRepo    DraftRepository
// 	htmlRepo     HTMLRepository
// 	standardRepo StandardRepository
// 	eventBus     *events.EventBus
//
// 	// For debouncing standard update events
// 	debounceInterval time.Duration
// 	debounceMutex    sync.Mutex
// 	pendingStandards map[int]bool
// 	debounceTimer    *time.Timer
// }
//
// // NewAuditTemplateService creates a new audit template service
// func NewAuditTemplateService(
// 	draftRepo DraftRepository,
// 	htmlRepo HTMLRepository,
// 	standardRepo StandardRepository,
// 	eventBus *events.EventBus,
// ) *AuditTemplateService {
// 	service := &AuditTemplateService{
// 		draftRepo:        draftRepo,
// 		htmlRepo:         htmlRepo,
// 		standardRepo:     standardRepo,
// 		eventBus:         eventBus,
// 		debounceInterval: 5 * time.Second,
// 		pendingStandards: make(map[int]bool),
// 	}
//
// 	// Subscribe to standard change events
// 	eventBus.Subscribe(events.TopicEntityChange, service.HandleEntityChange)
//
// 	return service
// }
//
// // SetDebounceInterval allows changing the debounce interval (primarily for testing)
// func (s *AuditTemplateService) SetDebounceInterval(interval time.Duration) {
// 	s.debounceInterval = interval
// }
//
// // CreateTemplateDraft creates a new draft for an audit template
// func (s *AuditTemplateService) CreateTemplateDraft(
// 	ctx context.Context,
// 	standardID int,
// 	templateData map[string]interface{},
// ) (types.Draft, error) {
// 	// Verify standard exists
// 	standard, err := s.standardRepo.GetByIDStandard(ctx, types.Standard{ID: standardID})
// 	if err != nil {
// 		return types.Draft{}, fmt.Errorf("failed to get standard: %w", err)
// 	}
//
// 	// Convert template data to JSON
// 	jsonData, err := json.Marshal(templateData)
// 	if err != nil {
// 		return types.Draft{}, fmt.Errorf("failed to marshal template data: %w", err)
// 	}
//
// 	// Calculate diff for first version
// 	diffData := map[string]interface{}{}
// 	for key, value := range templateData {
// 		diffData[key] = map[string]interface{}{
// 			"old": "",
// 			"new": value,
// 		}
// 	}
//
// 	jsonDiff, err := json.Marshal(diffData)
// 	if err != nil {
// 		return types.Draft{}, fmt.Errorf("failed to marshal diff data: %w", err)
// 	}
//
// 	// Create draft object
// 	now := time.Now().UTC()
// 	draft := types.Draft{
// 		TypeID:    TEMPLATE_TYPE_ID,
// 		ObjectID:  standardID,
// 		StatusID:  DRAFT_STATUS_ID,
// 		Version:   1, // First version
// 		Data:      json.RawMessage(jsonData),
// 		Diff:      json.RawMessage(jsonDiff),
// 		UserID:    1, // Default user or could be passed as parameter
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 		ExpiresAt: now.Add(7 * 24 * time.Hour), // Expires in 7 days
// 	}
//
// 	// Save draft
// 	createdDraft, err := s.draftRepo.CreateDraft(ctx, draft)
// 	if err != nil {
// 		return types.Draft{}, fmt.Errorf("failed to create draft: %w", err)
// 	}
//
// 	return createdDraft, nil
// }
//
// // PreviewTemplateDraft generates HTML preview for a draft
// func (s *AuditTemplateService) PreviewTemplateDraft(ctx context.Context, draftID int) (string, error) {
// 	// Get draft by ID
// 	draft, err := s.draftRepo.GetDraftByID(ctx, draftID)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get draft: %w", err)
// 	}
//
// 	// Generate HTML from template
// 	html, err := s.generateTemplateHTML(draft)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate HTML: %w", err)
// 	}
//
// 	return html, nil
// }
//
// // PublishTemplate publishes a template draft
// func (s *AuditTemplateService) PublishTemplate(ctx context.Context, draftID int) error {
// 	// Get draft by ID
// 	draft, err := s.draftRepo.GetDraftByID(ctx, draftID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get draft: %w", err)
// 	}
//
// 	// Generate HTML from template
// 	html, err := s.generateTemplateHTML(draft)
// 	if err != nil {
// 		return fmt.Errorf("failed to generate HTML: %w", err)
// 	}
//
// 	// Update draft status to published
// 	draft.StatusID = PUBLISHED_STATUS_ID
// 	updatedDraft, err := s.draftRepo.UpdateDraft(ctx, draft)
// 	if err != nil {
// 		return fmt.Errorf("failed to update draft status: %w", err)
// 	}
//
// 	// Create materialized HTML query
// 	htmlQuery := types.MaterializedHTMLQuery{
// 		Name:        fmt.Sprintf("audit_template_%d_v%d", draft.ObjectID, draft.Version),
// 		ViewPath:    fmt.Sprintf("/web/audits/template/%d", draft.ObjectID),
// 		HTMLContent: html,
// 		Version:     draft.Version,
// 	}
//
// 	_, err = s.htmlRepo.CreateMaterializedHTMLQuery(ctx, htmlQuery)
// 	if err != nil {
// 		// Record error in draft
// 		updatedDraft.PublishError = err.Error()
// 		_, _ = s.draftRepo.UpdateDraft(ctx, updatedDraft)
//
// 		return fmt.Errorf("failed to store HTML: %w", err)
// 	}
//
// 	return nil
// }
//
// // GetTemplateDrafts gets all drafts for a standard
// func (s *AuditTemplateService) GetTemplateDrafts(ctx context.Context, standardID int) ([]types.Draft, error) {
// 	drafts, err := s.draftRepo.GetDraftsByTypeAndObject(ctx, TEMPLATE_TYPE_ID, standardID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get drafts: %w", err)
// 	}
//
// 	return drafts, nil
// }
//
// // HandleEntityChange handles entity change events
// func (s *AuditTemplateService) HandleEntityChange(ctx context.Context, payload interface{}) error {
// 	// Type assert to entity change payload
// 	changePayload, ok := payload.(events.EntityChangePayload)
// 	if !ok {
// 		return fmt.Errorf("invalid payload type")
// 	}
//
// 	// Only process standard updates
// 	if changePayload.EntityType != events.EntityStandard {
// 		return nil
// 	}
//
// 	// Handle standard update with debouncing
// 	return s.handleStandardUpdate(ctx, changePayload.EntityID)
// }
//
// // handleStandardUpdate handles standard update with debouncing
// func (s *AuditTemplateService) handleStandardUpdate(ctx context.Context, standardID int) error {
// 	s.debounceMutex.Lock()
// 	defer s.debounceMutex.Unlock()
//
// 	// Add standard to pending list
// 	s.pendingStandards[standardID] = true
//
// 	// If timer already running, let it handle the update
// 	if s.debounceTimer != nil {
// 		return nil
// 	}
//
// 	// Create new timer
// 	s.debounceTimer = time.AfterFunc(s.debounceInterval, func() {
// 		// Process all pending standards
// 		s.debounceMutex.Lock()
// 		standards := make([]int, 0, len(s.pendingStandards))
// 		for id := range s.pendingStandards {
// 			standards = append(standards, id)
// 		}
// 		s.pendingStandards = make(map[int]bool)
// 		s.debounceTimer = nil
// 		s.debounceMutex.Unlock()
//
// 		// Process each standard
// 		for _, id := range standards {
// 			s.regenerateTemplate(context.Background(), id)
// 		}
// 	})
//
// 	// For immediate execution (mainly for testing)
// 	if s.debounceInterval == 0 {
// 		s.debounceTimer.Stop()
// 		standards := make([]int, 0, len(s.pendingStandards))
// 		for id := range s.pendingStandards {
// 			standards = append(standards, id)
// 		}
// 		s.pendingStandards = make(map[int]bool)
// 		s.debounceTimer = nil
//
// 		// Process each standard
// 		for _, id := range standards {
// 			if err := s.regenerateTemplate(ctx, id); err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	return nil
// }
//
// // regenerateTemplate regenerates a template for a standard
// func (s *AuditTemplateService) regenerateTemplate(ctx context.Context, standardID int) error {
// 	// Get standard
// 	standard, err := s.standardRepo.GetByIDStandard(ctx, types.Standard{ID: standardID})
// 	if err != nil {
// 		return fmt.Errorf("failed to get standard: %w", err)
// 	}
//
// 	// Generate template data
// 	templateData := map[string]interface{}{
// 		"title":         fmt.Sprintf("%s Audit Template", standard.Name),
// 		"standard_id":   standard.ID,
// 		"standard_name": standard.Name,
// 		"updated_at":    time.Now().Format(time.RFC3339),
// 		// In a real implementation, would include standard requirements, sections, etc.
// 	}
//
// 	// Create draft
// 	_, err = s.CreateTemplateDraft(ctx, standardID, templateData)
// 	if err != nil {
// 		return fmt.Errorf("failed to create template draft: %w", err)
// 	}
//
// 	return nil
// }
//
// // generateTemplateHTML generates HTML from a template draft using templ
// func (s *AuditTemplateService) generateTemplateHTML(draft types.Draft) (string, error) {
// 	// Parse template data
// 	var templateData map[string]interface{}
// 	if err := json.Unmarshal(draft.Data, &templateData); err != nil {
// 		return "", fmt.Errorf("failed to parse template data: %w", err)
// 	}
//
// 	// Convert requirements to strong type if they exist
// 	var requirements []types.Requirement
// 	if reqs, ok := templateData["requirements"].([]interface{}); ok {
// 		for _, req := range reqs {
// 			if reqMap, ok := req.(map[string]interface{}); ok {
// 				var id int
// 				var name string
//
// 				if idVal, ok := reqMap["id"].(float64); ok {
// 					id = int(idVal)
// 				}
//
// 				if nameVal, ok := reqMap["name"].(string); ok {
// 					name = nameVal
// 				}
//
// 				requirements = append(requirements, types.Requirement{
// 					ID:   id,
// 					Name: name,
// 				})
// 			}
// 		}
// 	}
//
// 	// Create a template view model
// 	viewModel := templates.AuditTemplateViewModel{
// 		Title:        getStringValue(templateData, "title"),
// 		StandardID:   getIntValue(templateData, "standard_id"),
// 		StandardName: getStringValue(templateData, "standard_name"),
// 		Requirements: requirements,
// 	}
//
// 	// Render using templ
// 	var buf bytes.Buffer
// 	if err := templates.AuditTemplate(viewModel).Render(context.Background(), &buf); err != nil {
// 		return "", fmt.Errorf("failed to render template: %w", err)
// 	}
//
// 	return buf.String(), nil
// }
//
// // Helper functions for type conversion
// func getStringValue(data map[string]interface{}, key string) string {
// 	if val, ok := data[key].(string); ok {
// 		return val
// 	}
// 	return ""
// }
//
// func getIntValue(data map[string]interface{}, key string) int {
// 	if val, ok := data[key].(float64); ok {
// 		return int(val)
// 	}
// 	return 0
// }
