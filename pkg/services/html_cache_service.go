package services

import (
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/templates"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	// "github.com/a-h/templ"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// HTMLCacheService manages the generation and caching of HTML content
type HTMLCacheService struct {
	HTMLRepo         repositories.MaterializedHTMLQueryRepositoryInterface
	JSONRepo         repositories.MaterializedJSONQueryRepositoryInterface
	StandardRepo     repositories.StandardRepositoryInterface
	RequirementRepo  repositories.RequirementRepositoryInterface
	EventBus         *events.EventBus
	debounceTimers   map[string]*time.Timer
	debounceInterval time.Duration
	mutex            sync.Mutex
}

// NewHTMLCacheService creates a new HTMLCacheService
func NewHTMLCacheService(
	htmlRepo repositories.MaterializedHTMLQueryRepositoryInterface,
	jsonRepo repositories.MaterializedJSONQueryRepositoryInterface,
	standardRepo repositories.StandardRepositoryInterface,
	requirementRepo repositories.RequirementRepositoryInterface,
	eventBus *events.EventBus,
) *HTMLCacheService {
	service := &HTMLCacheService{
		HTMLRepo:         htmlRepo,
		JSONRepo:         jsonRepo,
		StandardRepo:     standardRepo,
		RequirementRepo:  requirementRepo,
		EventBus:         eventBus,
		debounceTimers:   make(map[string]*time.Timer),
		debounceInterval: 3 * time.Second, // Slightly longer than MaterializedJSONService
	}

	// Subscribe to materialized query events
	eventBus.Subscribe(events.MaterializedQueryCreated, service.handleMaterializedQueryEvent)
	eventBus.Subscribe(events.MaterializedQueryUpdated, service.handleMaterializedQueryEvent)

	return service
}

// HandleQueryEvent implements the events.HTMLCacheService interface
func (s *HTMLCacheService) HandleQueryEvent(ctx context.Context, eventType events.EventType, payload events.MaterializedQueryPayload) error {
	log.Printf("HTML Cache service handling query event: %s for query %s", eventType, payload.QueryName)
	return s.RefreshHTMLForQuery(ctx, payload.QueryName)
}

// RefreshHTMLForQuery implements the events.HTMLCacheService interface
func (s *HTMLCacheService) RefreshHTMLForQuery(ctx context.Context, queryName string) error {
	// Debounce updates to avoid multiple rapid refreshes
	s.debounceHTMLUpdate(queryName, func() {
		bgCtx := context.Background()

		// Determine what kind of query this is
		if strings.HasPrefix(queryName, "standard_full_") {
			standardID := extractIDFromQueryName(queryName, "standard_full_")
			if standardID > 0 {
				if err := s.regenerateHTMLForStandard(bgCtx, standardID); err != nil {
					log.Printf("Error regenerating HTML for standard %d: %v", standardID, err)
				}
			}
		} else if strings.HasPrefix(queryName, "standard_") {
			standardID := extractIDFromQueryName(queryName, "standard_")
			if standardID > 0 {
				if err := s.regenerateHTMLForStandard(bgCtx, standardID); err != nil {
					log.Printf("Error regenerating HTML for standard %d: %v", standardID, err)
				}
			}
		}
		// Add other query types as needed
	})

	return nil
}

// Internal event handler for materialized query events
func (s *HTMLCacheService) handleMaterializedQueryEvent(ctx context.Context, event events.Event) error {
	payload, err := events.GetMaterializedQueryPayload(event)
	if err != nil {
		return err
	}

	return s.HandleQueryEvent(ctx, event.Type, payload)
}

// RegenerateHTML forces regeneration of HTML for a standard
func (s *HTMLCacheService) RegenerateHTML(ctx context.Context, standardID int) error {
	return s.regenerateHTMLForStandard(ctx, standardID)
}

// GetCachedHTML retrieves pre-generated HTML for a specific view
func (s *HTMLCacheService) GetCachedHTML(ctx context.Context, viewName string) (string, bool, error) {
	// Create query object with the name
	query := types.MaterializedHTMLQuery{
		Name: viewName,
	}

	// Try to get the HTML from the repository
	htmlQuery, err := s.HTMLRepo.GetByNameMaterializedHTMLQuery(ctx, query.Name)
	if err != nil {
		return "", false, err
	}

	return htmlQuery.HTMLContent, true, nil
}

// Internal methods

func (s *HTMLCacheService) debounceHTMLUpdate(key string, fn func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Cancel existing timer if any
	if timer, exists := s.debounceTimers[key]; exists {
		timer.Stop()
	}

	// Create new timer
	s.debounceTimers[key] = time.AfterFunc(s.debounceInterval, fn)
}

func (s *HTMLCacheService) regenerateHTMLForStandard(ctx context.Context, standardID int) error {
	// Step 1: Get the materialized JSON data for this standard
	jsonQuery := types.MaterializedJSONQuery{
		Name: fmt.Sprintf("standard_full_%d", standardID),
	}

	materializedJSON, err := s.JSONRepo.GetByNameMaterializedJSONQuery(ctx, jsonQuery.Name)
	if err != nil {
		// Try to get the standard data directly if no materialized query exists
		standard := types.Standard{ID: standardID}
		standardData, err := s.StandardRepo.GetByIDWithFullHierarchyStandard(ctx, standard)
		if err != nil {
			return fmt.Errorf("failed to get standard data: %w", err)
		}

		// Create JSON data
		jsonData, err := json.Marshal(standardData)
		if err != nil {
			return fmt.Errorf("failed to marshal standard data: %w", err)
		}

		materializedJSON = types.MaterializedJSONQuery{
			Data: jsonData,
		}
	}

	// Step 2: Parse the JSON data
	var standardData map[string]any
	if err := json.Unmarshal(materializedJSON.Data, &standardData); err != nil {
		return fmt.Errorf("failed to unmarshal standard data: %w", err)
	}

	// Step 3: Create the different HTML views we need
	if err := s.generateAuditView(ctx, standardID, standardData); err != nil {
		return err
	}

	if err := s.generateRequirementsView(ctx, standardID, standardData); err != nil {
		return err
	}

	// Add more views as needed

	return nil
}

func (s *HTMLCacheService) generateAuditView(ctx context.Context, standardID int, standardData map[string]any) error {
	// Get the standard information for the view
	standard := types.Standard{ID: standardID}
	std, err := s.StandardRepo.GetByIDStandard(ctx, standard)
	if err != nil {
		return fmt.Errorf("failed to get standard: %w", err)
	}

	// Generate HTML using templ
	var buf bytes.Buffer
	component := templates.AuditView(std, standardData)
	if err := component.Render(ctx, &buf); err != nil {
		return fmt.Errorf("failed to render audit view HTML: %w", err)
	}

	// Create or update the HTML materialized query
	htmlQuery := types.MaterializedHTMLQuery{
		Name:        fmt.Sprintf("audit_view_%d", standardID),
		ViewPath:    fmt.Sprintf("/web/audits/standard/%d", standardID),
		HTMLContent: buf.String(),
	}

	// Check if it already exists
	existingQuery, err := s.HTMLRepo.GetByNameMaterializedHTMLQuery(ctx, htmlQuery.Name)
	if err == nil {
		// Update existing
		htmlQuery.ID = existingQuery.ID
		htmlQuery.Version = existingQuery.Version + 1
		_, err = s.HTMLRepo.UpdateMaterializedHTMLQuery(ctx, htmlQuery)
		return err
	}

	// Create new
	htmlQuery.Version = 1
	_, err = s.HTMLRepo.CreateMaterializedHTMLQuery(ctx, htmlQuery)
	return err
}

func (s *HTMLCacheService) generateRequirementsView(ctx context.Context, standardID int, standardData map[string]any) error {
	// Get the standard information for the view
	standard := types.Standard{ID: standardID}
	std, err := s.StandardRepo.GetByIDStandard(ctx, standard)
	if err != nil {
		return fmt.Errorf("failed to get standard: %w", err)
	}

	// Generate HTML using templ
	var buf bytes.Buffer
	component := templates.RequirementsView(std, standardData)
	if err := component.Render(ctx, &buf); err != nil {
		return fmt.Errorf("failed to render requirements view HTML: %w", err)
	}

	// Create or update the HTML materialized query
	htmlQuery := types.MaterializedHTMLQuery{
		Name:        fmt.Sprintf("requirements_view_%d", standardID),
		ViewPath:    fmt.Sprintf("/web/requirements/standard/%d", standardID),
		HTMLContent: buf.String(),
	}

	// Check if it already exists
	existingQuery, err := s.HTMLRepo.GetByNameMaterializedHTMLQuery(ctx, htmlQuery.Name)
	if err == nil {
		// Update existing
		htmlQuery.ID = existingQuery.ID
		htmlQuery.Version = existingQuery.Version + 1
		_, err = s.HTMLRepo.UpdateMaterializedHTMLQuery(ctx, htmlQuery)
		return err
	}

	// Create new
	htmlQuery.Version = 1
	_, err = s.HTMLRepo.CreateMaterializedHTMLQuery(ctx, htmlQuery)
	return err
}

// Helper functions

// Extract ID from a query name like "standard_full_123"
func extractIDFromQueryName(queryName, prefix string) int {
	if !strings.HasPrefix(queryName, prefix) {
		return 0
	}

	idStr := queryName[len(prefix):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0
	}

	return id
}
