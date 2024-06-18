package controllers

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HtmlIsoStandardController struct {
	Repo repositories.IsoStandardRepository
}

func NewHtmlIsoStandardController(repo repositories.IsoStandardRepository) *HtmlIsoStandardController {
	return &HtmlIsoStandardController{
		Repo: repo,
	}
}

func (cc *HtmlIsoStandardController) GetAllISOStandards(c *gin.Context) {
	standards, err := cc.Repo.GetAllISOStandards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, standards)
	templ.Handler(templates.ISOStandards(standards)).ServeHTTP(c.Writer, c.Request)
}

func (cc *HtmlIsoStandardController) GetISOStandardByID(c *gin.Context) {
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

func (cc *HtmlIsoStandardController) RenderAddISOStandardForm(c *gin.Context) {
	templ.Handler(templates.AddISOStandard()).ServeHTTP(c.Writer, c.Request)
}

func (cc *HtmlIsoStandardController) CreateISOStandard(c *gin.Context) {
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
	c.Redirect(http.StatusFound, "/html/iso_standards")
	// templ.Handler(templates.AddISOStandard()).ServeHTTP(c.Writer, c.Request)
}

func (cc *HtmlIsoStandardController) UpdateISOStandard(c *gin.Context) {
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

func (cc *HtmlIsoStandardController) DeleteISOStandard(c *gin.Context) {
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
