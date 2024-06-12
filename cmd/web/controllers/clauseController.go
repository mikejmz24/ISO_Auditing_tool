package html

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

func (cc *HtmlClauseController) ShowAddClauseForm(c *gin.Context) {
	templ.Handler(templates.AddClause()).ServeHTTP(c.Writer, c.Request)
}

func (cc *HtmlClauseController) AddClause(c *gin.Context) {
	var clause types.Clause

	clause.Name = c.PostForm("clauseName")
	sections := c.PostForm("sections")
	sectionList := strings.Split(sections, ",")
	for _, section := range sectionList {
		clause.Sections = append(clause.Sections, types.Section{Name: strings.TrimSpace(section)})
	}

	_, err := cc.Repo.CreateClause(clause)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add clause"})
		return
	}

	c.Redirect(http.StatusFound, "/clauses")
}
