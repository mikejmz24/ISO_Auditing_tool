package api

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiClauseController struct {
	Repo repositories.Repository
}

func NewApiClauseController(repo repositories.Repository) *ApiClauseController {
	return &ApiClauseController{
		Repo: repo,
	}
}

func (cc *ApiClauseController) GetAllClauses(c *gin.Context) {
	clauses, err := cc.Repo.GetAllClauses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": clauses})
}
