package html

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func (cc *HtmlClauseController) GetClauseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clause ID"})
		return
	}

	clause, err := cc.Repo.GetClauseByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clause)
}

func (cc *HtmlClauseController) CreateClause(c *gin.Context) {
	var clause types.Clause
	if err := c.ShouldBindJSON(&clause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := cc.Repo.CreateClause(clause)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (cc *HtmlClauseController) UpdateClause(c *gin.Context) {
	var clause types.Clause
	if err := c.ShouldBindJSON(&clause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Repo.UpdateClause(clause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (cc *HtmlClauseController) DeleteClause(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clause ID"})
		return
	}

	if err := cc.Repo.DeleteClause(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// ISO Standard handlers

func (cc *HtmlClauseController) GetAllISOStandards(c *gin.Context) {
	standards, err := cc.Repo.GetAllISOStandards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, standards)
	templ.Handler(templates.ISOStandards(standards)).ServeHTTP(c.Writer, c.Request)
}

func (cc *HtmlClauseController) GetISOStandardByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ISO standard ID"})
		return
	}

	standard, err := cc.Repo.GetISOStandardByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, standard)
}

func (cc *HtmlClauseController) RenderAddISOStandardForm(c *gin.Context) {
	templ.Handler(templates.AddISOStandard()).ServeHTTP(c.Writer, c.Request)
}

func (cc *HtmlClauseController) CreateISOStandard(c *gin.Context) {
	var standard types.ISOStandard
	// if err := c.ShouldBindJSON(&standard); err != nil {
	if err := c.ShouldBind(&standard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// id, err := cc.Repo.CreateISOStandard(standard)
	_, err := cc.Repo.CreateISOStandard(standard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusCreated, gin.H{"id": id})
	// c.Redirect(http.StatusFound, "/iso_standards")
	templ.Handler(templates.AddISOStandard()).ServeHTTP(c.Writer, c.Request)
}

func (cc *HtmlClauseController) UpdateISOStandard(c *gin.Context) {
	var standard types.ISOStandard
	if err := c.ShouldBindJSON(&standard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Repo.UpdateISOStandard(standard); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (cc *HtmlClauseController) DeleteISOStandard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ISO standard ID"})
		return
	}

	if err := cc.Repo.DeleteISOStandard(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Section handlers

func (cc *HtmlClauseController) GetAllSections(c *gin.Context) {
	sections, err := cc.Repo.GetAllSections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sections)
}

func (cc *HtmlClauseController) GetSectionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid section ID"})
		return
	}

	section, err := cc.Repo.GetSectionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, section)
}

func (cc *HtmlClauseController) CreateSection(c *gin.Context) {
	var section types.Section
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := cc.Repo.CreateSection(section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (cc *HtmlClauseController) UpdateSection(c *gin.Context) {
	var section types.Section
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Repo.UpdateSection(section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (cc *HtmlClauseController) DeleteSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid section ID"})
		return
	}

	if err := cc.Repo.DeleteSection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Question handlers

func (cc *HtmlClauseController) GetAllQuestions(c *gin.Context) {
	questions, err := cc.Repo.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (cc *HtmlClauseController) GetQuestionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question ID"})
		return
	}

	question, err := cc.Repo.GetQuestionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

func (cc *HtmlClauseController) CreateQuestion(c *gin.Context) {
	var question types.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := cc.Repo.CreateQuestion(question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (cc *HtmlClauseController) UpdateQuestion(c *gin.Context) {
	var question types.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Repo.UpdateQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (cc *HtmlClauseController) DeleteQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question ID"})
		return
	}

	if err := cc.Repo.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
