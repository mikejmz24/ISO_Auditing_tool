// Only handles API request validation and response formatting for drafts
package controllers

import (
	"ISO_Auditing_Tool/pkg/services"
	"ISO_Auditing_Tool/pkg/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiDraftController struct {
	Service *services.DraftService
}

// NewAPIDraftController creates a new instance of ApiDraftController
func NewAPIDraftController(service *services.DraftService) *ApiDraftController {
	return &ApiDraftController{Service: service}
}

func (cc *ApiDraftController) Create(c *gin.Context) {
	var draft types.Draft

	if err := c.ShouldBindJSON(&draft); err != nil {
		// TODO: Implement custom errors
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft, err := cc.Service.Create(c.Request.Context(), draft)
	if err != nil {
		// TODO: Implement custom errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": draft.ID})
}

func (cc *ApiDraftController) Update(c *gin.Context) {
	var draft types.Draft

	if err := c.ShouldBindJSON(&draft); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft.ID = id
	if draft, err := cc.Service.Update(c.Request.Context(), draft); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "draft": draft})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated row with ID": draft.ID})
}
