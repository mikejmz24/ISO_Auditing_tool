package controllers

import (
	"ISO_Auditing_Tool/cmd/api/controllers"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebIsoStandardController struct {
	ApiController *controllers.ApiIsoStandardController
}

func NewWebIsoStandardController(apiController *controllers.ApiIsoStandardController) *WebIsoStandardController {
	return &WebIsoStandardController{ApiController: apiController}
}

func (wc *WebIsoStandardController) GetAllISOStandards(c *gin.Context) {
	apiContext := &gin.Context{Request: c.Request, Writer: c.Writer}
	wc.ApiController.GetAllISOStandards(apiContext)
	var isoStandards []types.ISOStandard
	if err := apiContext.BindJSON(&isoStandards); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bind JSON"})
		return
	}
	templ.Handler(templates.ISOStandards(isoStandards)).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) GetISOStandardByID(c *gin.Context) {
	// id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ISO standard ID"})
	// 	return
	// }
	//
	// standard, err := cc.Repo.GetISOStandardByID(id)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	apiContext := &gin.Context{Request: c.Request, Writer: c.Writer}
	wc.ApiController.GetISOStandardByID(apiContext)
	var isoStandard types.ISOStandard
	if err := apiContext.BindJSON(&isoStandard); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bind JSON"})
		return
	}

	// c.JSON(http.StatusOK, standard)
	c.HTML(http.StatusOK, "iso_standard.html", gin.H{"isoStandard": isoStandard})
}

func (wc *WebIsoStandardController) RenderAddISOStandardForm(c *gin.Context) {
	templ.Handler(templates.AddISOStandard()).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) CreateISOStandard(c *gin.Context) {
	// var standard types.ISOStandard
	// // if err := c.ShouldBindJSON(&standard); err != nil {
	// if err := c.ShouldBind(&standard); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	//
	// // id, err := cc.Repo.CreateISOStandard(standard)
	// _, err := cc.Repo.CreateISOStandard(standard)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	apiContext := &gin.Context{Request: c.Request, Writer: c.Writer}
	wc.ApiController.CreateISOStandard(apiContext)

	// c.JSON(http.StatusCreated, gin.H{"id": id})
	c.Redirect(http.StatusFound, "/html/iso_standards")
	// templ.Handler(templates.AddISOStandard()).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) UpdateISOStandard(c *gin.Context) {
	// var standard types.ISOStandard
	// if err := c.ShouldBindJSON(&standard); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	//
	// if err := cc.Repo.UpdateISOStandard(standard); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	apiContext := &gin.Context{Request: c.Request, Writer: c.Writer}
	wc.ApiController.UpdateISOStandard(apiContext)
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (wc *WebIsoStandardController) DeleteISOStandard(c *gin.Context) {
	// id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ISO standard ID"})
	// 	return
	// }
	//
	// if err := cc.Repo.DeleteISOStandard(id); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	apiContext := &gin.Context{Request: c.Request, Writer: c.Writer}
	wc.ApiController.DeleteISOStandard(apiContext)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
