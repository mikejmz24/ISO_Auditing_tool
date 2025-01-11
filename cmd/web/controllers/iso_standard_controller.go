package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/templates"
	"ISO_Auditing_Tool/templates/iso_standards"

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
	templ.Handler(templates.Base(iso_standards.List(isoStandards))).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) GetISOStandardByID(c *gin.Context) {
	id := c.Param("id")
	isoStandard, err := wc.fetchISOStandardByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ISO standard not found"})
		return
	}
	templ.Handler(templates.Base(iso_standards.Detail(*isoStandard))).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) RenderAddISOStandardForm(c *gin.Context) {
	templ.Handler(templates.Base(iso_standards.Add())).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) CreateISOStandard(c *gin.Context) {
	var formData types.ISOStandardForm
	// var formData types.ISOStandard
	var customErr custom_errors.CustomError

	if c.ContentType() != "application/x-www-form-urlencoded" {
		customErr = *custom_errors.ErrInvalidFormData
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	if err := c.Request.ParseForm(); err != nil {
		customErr = *custom_errors.ErrInvalidFormData
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	if len(c.Request.PostForm) == 0 {
		customErr = *custom_errors.EmptyData("Form")
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	if len(c.Request.PostForm) == 1 {
		for key, values := range c.Request.PostForm {
			if key == "" || len(values) == 0 || len(values) == 1 && values[0] == "" && key != "name" {
				customErr = *custom_errors.ErrInvalidFormData
				c.JSON(customErr.StatusCode, customErr)
				return
			}
		}
	}

	if err := c.ShouldBind(&formData); err != nil {
		customErr = *custom_errors.ErrInvalidFormData
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	if _, exists := c.Request.PostForm["name"]; !exists {
		customErr = *custom_errors.MissingField("name")
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	if formData.Name == "" {
		customErr = *custom_errors.EmptyField("string", "name")
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	if isInvalidString(formData.Name) {
		customErr = *custom_errors.InvalidDataType("name", "string")
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	// Convert form data to ISOStandard type
	isoStandard := formData.ToISOStandard()

	// Marshal to JSON for API call
	jsonData, err := json.Marshal(isoStandard)
	// jsonData, err := json.Marshal(formData)
	if err != nil {
		customErr = *custom_errors.NewCustomError(http.StatusBadRequest, "Failed to process data", nil)
		c.JSON(customErr.StatusCode, customErr)
		// c.JSON(201, jsonData)
		return
	}

	// Create test context for API call
	recorder := httptest.NewRecorder()
	apiContext, _ := gin.CreateTestContext(recorder)
	// apiContext.Request = c.Request.Clone(c.Request.Context())
	apiContext.Request = c.Request
	apiContext.Request.Header.Set("Content-Type", "application/json")
	apiContext.Request.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	wc.ApiController.CreateISOStandard(apiContext)

	response := recorder.Result()
	defer response.Body.Close()
	responseBody, _ := io.ReadAll(response.Body)

	if strings.Contains(string(responseBody), "error") {
		customErr = *custom_errors.NewCustomError(http.StatusInternalServerError, "Failed to create ISO Standard", nil)
		c.JSON(customErr.StatusCode, customErr)
		return
	}

	switch response.StatusCode {
	case http.StatusCreated, http.StatusOK:
		c.Redirect(http.StatusFound, "/web/iso_standards")
		return
	case http.StatusBadRequest:
		customErr = *custom_errors.ErrInvalidFormData
		c.JSON(customErr.StatusCode, customErr)
		return
	default:
		c.JSON(response.StatusCode, string(responseBody))
	}
}

func (wc *WebIsoStandardController) UpdateISOStandard(c *gin.Context) {
	var formData map[string]string
	if err := c.Bind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	jsonData, err := json.Marshal(formData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal form data"})
		return
	}

	apiContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	apiContext.Request = c.Request
	apiContext.Writer = c.Writer
	apiContext.Request.Body = io.NopCloser(strings.NewReader(string(jsonData)))
	apiContext.Request.Header.Set("Content-Type", "application/json")

	wc.ApiController.UpdateISOStandard(apiContext)

	if apiContext.Writer.Status() == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	} else {
		c.JSON(apiContext.Writer.Status(), gin.H{"error": "Failed to update ISO standard"})
	}
}

func (wc *WebIsoStandardController) DeleteISOStandard(c *gin.Context) {
	apiContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	apiContext.Request = c.Request
	apiContext.Writer = c.Writer

	wc.ApiController.DeleteISOStandard(apiContext)

	if apiContext.Writer.Status() == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	} else {
		c.JSON(apiContext.Writer.Status(), gin.H{"error": "Failed to delete ISO standard"})
	}
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

// Helper function to fetch a single ISO standard
func (wc *WebIsoStandardController) fetchISOStandardByID(id string) (*types.ISOStandard, error) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/iso_standards/"+id, nil)
	apiContext, _ := gin.CreateTestContext(recorder)
	apiContext.Request = req
	apiContext.Params = gin.Params{{Key: "id", Value: id}}

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

func isInvalidString(input string) bool {
	if input == "true" || input == "false" {
		return true
	}

	if _, err := strconv.Atoi(input); err == nil {
		return true
	}

	if _, err := strconv.ParseFloat(input, 64); err == nil {
		return true
	}

	return false
}
