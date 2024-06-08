package html

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HtmlClauseController struct {
	Repo repositories.Repository
}

func NewHtmlClauseController(repo repositories.Repository) *HtmlClauseController {
	return &HtmlClauseController{
		Repo: repo,
	}
}

func (cc *HtmlClauseController) GetAllClauses(c *gin.Context) {
	clauses, err := cc.Repo.GetAllClauses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	templ.Handler(templates.Clauses(clauses)).ServeHTTP(c.Writer, c.Request)
}
