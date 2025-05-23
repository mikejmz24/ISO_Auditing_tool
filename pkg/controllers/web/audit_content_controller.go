// pkg/controllers/web/audit_content_controller.go
package controllers

import (
	"ISO_Auditing_Tool/pkg/services"
	"ISO_Auditing_Tool/pkg/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebAuditContentController struct {
	Service *services.AuditContentService
}

// NewWebAuditContentController creates a new WebAuditContentController
func NewWebAuditContentController(service *services.AuditContentService) *WebAuditContentController {
	return &WebAuditContentController{Service: service}
}

// RenderEditRequirement renders the requirement editing form
func (c *WebAuditContentController) RenderEditRequirement(ctx *gin.Context) {
	requirementIDStr := ctx.Param("requirement_id")
	requirementID, err := strconv.Atoi(requirementIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requirement ID"})
		return
	}

	// Get the requirement
	requirement, err := c.Service.GetRequirement(ctx.Request.Context(), requirementID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Requirement not found"})
		return
	}

	// Render the edit form template
	// Note: You'll need to create the actual template file
	ctx.HTML(http.StatusOK, "admin/requirement_edit.html", gin.H{
		"requirement": requirement,
		"title":       "Edit Requirement",
	})
}

// UpdateRequirement handles the requirement update form submission
func (c *WebAuditContentController) UpdateRequirement(ctx *gin.Context) {
	requirementIDStr := ctx.Param("requirement_id")
	requirementID, err := strconv.Atoi(requirementIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requirement ID"})
		return
	}

	// Parse form data
	var form types.RequirementEditForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Validate form
	if err := form.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the requirement ID matches
	if form.RequirementID != requirementID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Requirement ID mismatch"})
		return
	}

	// Get current user (you'll need to implement user session management)
	userID := c.getCurrentUserID(ctx) // Placeholder - implement based on your auth system

	// Modify the requirement
	err = c.Service.ModifyRequirementDescription(
		ctx.Request.Context(),
		requirementID,
		form.Description,
		form.Reason,
		userID,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update requirement"})
		return
	}

	// Redirect back to the requirement view or success page
	ctx.Redirect(http.StatusFound, "/admin/requirements/"+requirementIDStr)
}

// ListRequirements shows all requirements for a standard (for admin to choose which to edit)
func (c *WebAuditContentController) ListRequirements(ctx *gin.Context) {
	standardIDStr := ctx.Param("standard_id")
	standardID, err := strconv.Atoi(standardIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid standard ID"})
		return
	}

	// You'll need to add a method to get requirements by standard
	// For now, this is a placeholder
	ctx.HTML(http.StatusOK, "admin/requirements_list.html", gin.H{
		"standard_id": standardID,
		"title":       "Manage Requirements",
	})
}

// getCurrentUserID is a placeholder for getting the current user ID
// Implement this based on your authentication system
func (c *WebAuditContentController) getCurrentUserID(ctx *gin.Context) int {
	// Placeholder implementation
	// In a real system, you'd get this from session, JWT, etc.
	return 1 // Default admin user
}
