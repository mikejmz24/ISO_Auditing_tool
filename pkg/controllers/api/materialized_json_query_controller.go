// pkg/controllers/api/materialized_json_query_controller.go
package controllers

import (
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/services"
	"ISO_Auditing_Tool/pkg/types"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiMaterializedJSONQueryController struct {
	JSONService *services.MaterializedJSONService
	HTMLService *services.HTMLCacheService
	EventBus    *events.EventBus
}

// MaterializedJSONController takes MaterializedJSONService, HTMLCacheService and Event Bus
func NewApiMaterializedJSONQueryController(
	jsonService *services.MaterializedJSONService,
	htmlService *services.HTMLCacheService,
	eventBus *events.EventBus,
) *ApiMaterializedJSONQueryController {
	return &ApiMaterializedJSONQueryController{
		JSONService: jsonService,
		HTMLService: htmlService,
		EventBus:    eventBus,
	}
}

// GetByName retrieves a materialized JSON query by name
func (c *ApiMaterializedJSONQueryController) GetByName(ctx *gin.Context) {
	name := ctx.Param("name")

	// Create a query object with the name
	query := types.MaterializedJSONQuery{
		Name: name,
	}

	// Get the query from the repository
	materializedQuery, err := c.JSONService.JSONRepo.GetByNameMaterializedJSONQuery(ctx.Request.Context(), query.Name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, materializedQuery)
}

// CreateOrUpdateJSONQuery creates or updates a JSON materialized query
func (c *ApiMaterializedJSONQueryController) CreateOrUpdateJSONQuery(ctx *gin.Context) {
	var requestData types.MaterializedJSONQuery

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if requestData.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query_name is required"})
		return
	}

	// Check if query exists already
	existingQuery, err := c.JSONService.JSONRepo.GetByNameMaterializedJSONQuery(ctx.Request.Context(), requestData.Name)

	var result types.MaterializedJSONQuery
	var statusCode int

	if err != nil || existingQuery.ID == 0 {
		// New query
		requestData.Version = 1
		result, err = c.JSONService.JSONRepo.CreateMaterializedJSONQuery(ctx.Request.Context(), requestData)
		statusCode = http.StatusCreated

		if err == nil {
			// Publish event for creation
			event := events.CreateMaterializedQueryEvent(result)
			c.EventBus.AsyncPublish(ctx.Request.Context(), event)
		}
	} else {
		// Existing query - update it
		requestData.ID = existingQuery.ID
		requestData.Version = existingQuery.Version + 1
		result, err = c.JSONService.JSONRepo.UpdateMaterializedJSONQuery(ctx.Request.Context(), requestData)
		statusCode = http.StatusOK

		if err == nil {
			// Publish event for update
			event := events.UpdateMaterializedQueryEvent(result)
			c.EventBus.AsyncPublish(ctx.Request.Context(), event)
		}
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(statusCode, result)
}

// // CreateOrUpdateHTMLQuery creates or updates an HTML materialized query
// func (cc *ApiMaterializedQueryController) CreateOrUpdateHTMLQuery(c *gin.Context) {
// 	var requestData types.MaterializedHTMLQuery
//
// 	if err := c.ShouldBindJSON(&requestData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	// Validate required fields
// 	if requestData.Name == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "query_name is required"})
// 		return
// 	}
//
// 	if requestData.ViewPath == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "view_path is required"})
// 		return
// 	}
//
// 	// Check if query exists already
// 	existingQuery, err := cc.HTMLService.HTMLRepo.GetByName(c.Request.Context(), requestData)
//
// 	var result types.MaterializedHTMLQuery
// 	var statusCode int
//
// 	if err != nil || existingQuery.ID == 0 {
// 		// New query
// 		requestData.Version = 1
// 		result, err = cc.HTMLService.HTMLRepo.Create(c.Request.Context(), requestData)
// 		statusCode = http.StatusCreated
// 	} else {
// 		// Existing query - update it
// 		requestData.ID = existingQuery.ID
// 		requestData.Version = existingQuery.Version + 1
// 		result, err = cc.HTMLService.HTMLRepo.Update(c.Request.Context(), requestData)
// 		statusCode = http.StatusOK
// 	}
//
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	c.JSON(statusCode, result)
// }

// // GetByEntityType retrieves all materialized JSON queries for a specific entity type
// func (c *ApiMaterializedJSONQueryController) GetByEntityType(ctx *gin.Context) {
// 	entityType := ctx.Param("entity_type")
//
// 	// Get all queries for this entity type
// 	queries, err := c.Service.JSONRepo.GetByEntityType(ctx.Request.Context(), entityType)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, queries)
// }

// GetByEntityTypeAndID retrieves a materialized JSON query for a specific entity type and ID
func (c *ApiMaterializedJSONQueryController) GetByEntityTypeAndID(ctx *gin.Context) {
	entityType := ctx.Param("entity_type")
	entityIDStr := ctx.Param("entity_id")

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Create a query object with the entity type and ID
	query := types.MaterializedJSONQuery{
		EntityType: entityType,
		EntityID:   entityID,
	}

	// Get the query from the repository
	materializedQuery, err := c.JSONService.JSONRepo.GetByEntityTypeAndIDMaterializedJSONQuery(ctx.Request.Context(), query.EntityType, query.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, materializedQuery)
}

// RefreshEntityData manually refreshes the materialized JSON data for an entity
func (c *ApiMaterializedJSONQueryController) RefreshEntityData(ctx *gin.Context) {
	// Parse request
	var request struct {
		EntityType string `json:"entity_type" binding:"required"`
		EntityID   int    `json:"entity_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to the correct entity type constant
	var entityType events.EntityType
	switch request.EntityType {
	case "standard":
		entityType = events.EntityStandard
	case "requirement":
		entityType = events.EntityRequirement
	case "question":
		entityType = events.EntityQuestion
	case "evidence":
		entityType = events.EntityEvidence
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity type"})
		return
	}

	// Create an entity change event
	event := events.NewEntityChangeEvent(
		entityType,
		request.EntityID,
		events.ChangeUpdated,
		"",  // No specific affected query
		"",  // No parent type
		nil, // No parent ID
		nil, // No data
	)

	// Publish the event
	if err := c.EventBus.Publish(ctx.Request.Context(), event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "refresh triggered"})
}

// RefreshAllStandards refreshes all standards and their hierarchy
func (c *ApiMaterializedJSONQueryController) RefreshAllStandards(ctx *gin.Context) {
	// Get all standards
	standards, err := c.JSONService.StandardRepo.GetAllStandards(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh each standard
	for _, standard := range standards {
		event := events.NewEntityChangeEvent(
			events.EntityStandard,
			standard.ID,
			events.ChangeUpdated,
			"",       // No specific affected query
			"",       // No parent type
			nil,      // No parent ID
			standard, // Include the standard data
		)

		// Publish asynchronously to not block
		c.EventBus.AsyncPublish(ctx.Request.Context(), event)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "refresh triggered",
		"count":  len(standards),
	})
}

// ForceRegenerateHTML forces regeneration of HTML caches for a standard
func (c *ApiMaterializedJSONQueryController) ForceRegenerateHTML(ctx *gin.Context) {
	standardIDStr := ctx.Param("standard_id")

	standardID, err := strconv.Atoi(standardIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid standard ID"})
		return
	}

	// Trigger HTML regeneration
	if err := c.HTMLService.RegenerateHTML(ctx.Request.Context(), standardID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "HTML regenerated"})
}

// GetStandardWithHierarchy gets a materialized standard with its full hierarchy
func (c *ApiMaterializedJSONQueryController) GetStandardWithHierarchy(ctx *gin.Context) {
	standardIDStr := ctx.Param("standard_id")

	standardID, err := strconv.Atoi(standardIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid standard ID"})
		return
	}

	// Try to get from materialized query first
	query := types.MaterializedJSONQuery{
		Name: fmt.Sprintf("standard_full_%d", standardID),
	}

	materializedQuery, err := c.JSONService.JSONRepo.GetByNameMaterializedJSONQuery(ctx.Request.Context(), query.Name)
	if err == nil {
		// Return the cached data
		ctx.Header("Content-Type", "application/json")
		ctx.Writer.Write(materializedQuery.Data)
		return
	}

	// If not found, try to generate it
	standard := types.Standard{ID: standardID}
	standardWithHierarchy, err := c.JSONService.StandardRepo.GetByIDWithFullHierarchyStandard(ctx.Request.Context(), standard)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Standard not found"})
		return
	}

	// Return the freshly generated data
	ctx.JSON(http.StatusOK, standardWithHierarchy)

	// Trigger async update of the materialized query
	event := events.NewEntityChangeEvent(
		events.EntityStandard,
		standardID,
		events.ChangeUpdated,
		"",                    // No specific affected query
		"",                    // No parent type
		nil,                   // No parent ID
		standardWithHierarchy, // Include the data
	)

	c.EventBus.AsyncPublish(ctx.Request.Context(), event)
}
