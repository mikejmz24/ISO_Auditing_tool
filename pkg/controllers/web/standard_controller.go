// Only handles HTML request validation and response formatting for Standards
package controllers

import (
	"ISO_Auditing_Tool/pkg/services"
	"github.com/gin-gonic/gin"
)

type WebStandardController struct {
	Service *services.StandardService
}

func NewWebStandardController(service *services.StandardService) *WebStandardController {
	return &WebStandardController{Service: service}
}

func (cc *WebStandardController) GetByID(c *gin.Context) {

}
