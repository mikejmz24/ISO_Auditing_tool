package api

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"net/http"
	"strconv"

	"ISO_Auditing_Tool/pkg/types"
	"github.com/gin-gonic/gin"
)

type ApiClauseController struct {
	Repo repositories.Repository
}

// NewClauseController creates a new ApiClauseController
func NewApiClauseController(repo repositories.Repository) *ApiClauseController {
	return &ApiClauseController{Repo: repo}
}

// Get all clauses
func (cc *ApiClauseController) GetAllClauses(c *gin.Context) {
	clauses, err := cc.Repo.GetAllClauses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": clauses})
}

// Get clause by ID
func (cc *ApiClauseController) GetClauseByID(c *gin.Context) {
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

// Create a new clause
func (cc *ApiClauseController) CreateClause(c *gin.Context) {
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

// Update a clause
func (cc *ApiClauseController) UpdateClause(c *gin.Context) {
	var clause types.Clause
	if err := c.ShouldBindJSON(&clause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cc.Repo.UpdateClause(clause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, clause)
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete a clause
func (cc *ApiClauseController) DeleteClause(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clause ID"})
	}
	if err := cc.Repo.DeleteClause(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Clause deleted"})
}

// Get all ISO standards
func (cc *ApiClauseController) GetAllISOStandards(c *gin.Context) {
	var isoStandards []types.ISOStandard
	isoStandards, err := cc.Repo.GetAllISOStandards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, isoStandards)
}

// Get ISO standard by ID
func (cc *ApiClauseController) GetISOStandardByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ISO standard ID"})
		return
	}
	isoStandard, err := cc.Repo.GetISOStandardByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, isoStandard)
}

// Create a new ISO standard
func (cc *ApiClauseController) CreateISOStandard(c *gin.Context) {
	var isoStandard types.ISOStandard
	if err := c.ShouldBindJSON(&isoStandard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := cc.Repo.CreateISOStandard(isoStandard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// Update an ISO standard
func (cc *ApiClauseController) UpdateISOStandard(c *gin.Context) {
	var isoStandard types.ISOStandard
	if err := c.ShouldBindJSON(&isoStandard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cc.Repo.UpdateISOStandard(isoStandard); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete an ISO standard
func (cc *ApiClauseController) DeleteISOStandard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ISO standard ID"})
		return
	}
	if err := cc.Repo.DeleteISOStandard(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ISO standard deleted"})
}

// Get all sections
func (cc *ApiClauseController) GetAllSections(c *gin.Context) {
	sections, err := cc.Repo.GetAllSections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sections)
}

// Get section by ID
func (cc *ApiClauseController) GetSectionByID(c *gin.Context) {
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

// Create a new section
func (cc *ApiClauseController) CreateSection(c *gin.Context) {
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

// Update a section
func (cc *ApiClauseController) UpdateSection(c *gin.Context) {
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

// Delete a section
func (cc *ApiClauseController) DeleteSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invelid section ID"})
		return
	}
	if err := cc.Repo.DeleteSection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Section deleted"})
}

// Get all questions
func (cc *ApiClauseController) GetAllQuestions(c *gin.Context) {
	var questions []types.Question
	questions, err := cc.Repo.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// Get question by ID
func (cc *ApiClauseController) GetQuestionByID(c *gin.Context) {
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

// Create a new question
func (cc *ApiClauseController) CreateQuestion(c *gin.Context) {
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

// Update a question
func (cc *ApiClauseController) UpdateQuestion(c *gin.Context) {
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

// Delete a question
func (cc *ApiClauseController) DeleteQuestion(c *gin.Context) {
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
