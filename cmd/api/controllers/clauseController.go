package controllers

import (
	"ISO_Auditing_Tool/cmd/api/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClauseController struct {
	Repo repositories.ClauseRepository
}

func NewClauseController(repo repositories.ClauseRepository) *ClauseController {
	return &ClauseController{
		Repo: repo,
	}
}

func (cc *ClauseController) GetAllClauses(c *gin.Context) {
	clauses, err := cc.Repo.GetAllClauses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": clauses})
}
