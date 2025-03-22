// Only handles API request validation and response formattinf for materialized queries
package controllers

import (
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/services"
	"ISO_Auditing_Tool/pkg/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiMaterializedQueryController struct {
	Service *services.MaterializedQueryService
}

func NewApiMaterializedQueryController(service *services.MaterializedQueryService) *ApiMaterializedQueryController {
	return &ApiMaterializedQueryController{Service: service}
}

func (cc *ApiMaterializedQueryController) GetByName(c *gin.Context) {
	name := c.Param("name")

	materializedQuery, err := cc.Service.GetByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, materializedQuery)
}

func (cc *ApiMaterializedQueryController) Create(c *gin.Context) {

}

func (cc *ApiMaterializedQueryController) CreateOrUpdateMaterializedQuery(c *gin.Context) {
	var requestData types.MaterializedQuery

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := events.CreateMaterializedQueryEvent(requestData)

	if err := cc.Service.PublishEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Publish error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Materialized query creation/update request accepted"})
}
