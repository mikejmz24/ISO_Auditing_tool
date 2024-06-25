package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// Interface for the API controller to allow for easier testing and mocking
type ApiIsoStandardController interface {
	GetAllISOStandards(c *gin.Context)
	GetISOStandardByID(c *gin.Context)
	CreateISOStandard(c *gin.Context)
	UpdateISOStandard(c *gin.Context)
	DeleteISOStandard(c *gin.Context)
}

type WebIsoStandardController struct {
	ApiController ApiIsoStandardController
}

func NewWebIsoStandardController(apiController ApiIsoStandardController) *WebIsoStandardController {
	return &WebIsoStandardController{ApiController: apiController}
}

func (wc *WebIsoStandardController) GetAllISOStandards(c *gin.Context) {
	isoStandards, err := wc.fetchAllISOStandards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ISO standards"})
		return
	}
	templ.Handler(templates.ISOStandards(isoStandards)).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) GetISOStandardByID(c *gin.Context) {
	id := c.Param("id")
	isoStandard, err := wc.fetchISOStandardByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ISO standard"})
		return
	}
	c.HTML(http.StatusOK, "iso_standard.html", gin.H{"isoStandard": isoStandard})
}

func (wc *WebIsoStandardController) RenderAddISOStandardForm(c *gin.Context) {
	templ.Handler(templates.AddISOStandard()).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) CreateISOStandard(c *gin.Context) {
	apiContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	apiContext.Request = c.Request
	apiContext.Writer = c.Writer

	wc.ApiController.CreateISOStandard(apiContext)
	c.Redirect(http.StatusFound, "/html/iso_standards")
}

func (wc *WebIsoStandardController) UpdateISOStandard(c *gin.Context) {
	apiContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	apiContext.Request = c.Request
	apiContext.Writer = c.Writer

	wc.ApiController.UpdateISOStandard(apiContext)
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (wc *WebIsoStandardController) DeleteISOStandard(c *gin.Context) {
	apiContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	apiContext.Request = c.Request
	apiContext.Writer = c.Writer

	wc.ApiController.DeleteISOStandard(apiContext)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Helper functions for fetching data from the API controller
func (wc *WebIsoStandardController) fetchAllISOStandards() ([]types.ISOStandard, error) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/iso_standards", nil)
	apiContext, _ := gin.CreateTestContext(recorder)
	apiContext.Request = req

	wc.ApiController.GetAllISOStandards(apiContext)

	if recorder.Code != http.StatusOK {
		log.Printf("Error fetching ISO standards: %s", recorder.Body.String())
		return nil, fmt.Errorf("error fetching ISO standards")
	}

	var isoStandards []types.ISOStandard
	if err := json.Unmarshal(recorder.Body.Bytes(), &isoStandards); err != nil {
		log.Printf("Error unmarshalling ISO standards: %v", err)
		return nil, err
	}
	return isoStandards, nil
}

func (wc *WebIsoStandardController) fetchISOStandardByID(id string) (*types.ISOStandard, error) {
	requestURL := "/api/iso_standards/" + id
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, requestURL, nil)
	apiContext, _ := gin.CreateTestContext(recorder)
	apiContext.Request = req

	wc.ApiController.GetISOStandardByID(apiContext)

	if recorder.Code != http.StatusOK {
		log.Printf("Error fetching ISO standard by ID: %s", recorder.Body.String())
		return nil, fmt.Errorf("error fetching ISO standard by ID")
	}

	var isoStandard types.ISOStandard
	if err := json.Unmarshal(recorder.Body.Bytes(), &isoStandard); err != nil {
		log.Printf("Error unmarshalling ISO standard: %v", err)
		return nil, err
	}
	return &isoStandard, nil
}
