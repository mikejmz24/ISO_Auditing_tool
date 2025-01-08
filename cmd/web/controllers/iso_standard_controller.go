package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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
	// Check content type to determine how to handle the input
	contentType := c.GetHeader("Content-Type")

	var formData types.ISOStandardForm
	var err error

	if strings.Contains(contentType, "application/json") {
		// Handle JSON input
		if err = c.ShouldBindJSON(&formData); err != nil {
			// Check for specific validation errors
			if strings.Contains(err.Error(), "Name") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": custom_errors.MissingField("name").Message,
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad JSON error at the beginning",
			})
			return
		}
	} else {
		// Handle form input
		// if err = c.ShouldBind(&formData); err != nil {
		// 	// For form submissions, render the form with errors
		// 	templ.Handler(templates.Base(iso_standards.Add())).ServeHTTP(c.Writer, c.Request)
		// 	return
		// }
		// TODO: Ensure formData includes the "Name" field and is correctly bound
		if err = c.ShouldBind(&formData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": custom_errors.ErrInvalidFormData.Error(),
			})
			return
		}
	}

	// Convert form data to ISOStandard type
	isoStandard := formData.ToISOStandard()

	// Marshal to JSON for API call
	jsonData, err := json.Marshal(isoStandard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process form data",
		})
		return
	}

	// Create test context for API call
	recorder := httptest.NewRecorder()
	apiContext, _ := gin.CreateTestContext(recorder)
	apiContext.Request = c.Request.Clone(c.Request.Context())
	apiContext.Request.Header.Set("Content-Type", "application/json")
	apiContext.Request.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	wc.ApiController.CreateISOStandard(apiContext)

	response := recorder.Result()
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusCreated:
		// For JSON requests, return status created
		if strings.Contains(contentType, "application/json") {
			c.JSON(http.StatusCreated, gin.H{"status": "created"})
			return
		}
		// For form submissions, redirect
		c.Redirect(http.StatusFound, "/web/iso_standards")
	case http.StatusBadRequest:
		if strings.Contains(contentType, "application/json") {
			responseBody, _ := io.ReadAll(response.Body)
			c.JSON(http.StatusBadRequest, string(responseBody))
		} else {
			// templ.Handler(templates.Base(iso_standards.Add())).ServeHTTP(c.Writer, c.Request)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": custom_errors.ErrInvalidFormData.Error(),
			})
		}
	default:
		responseBody, _ := io.ReadAll(response.Body)
		c.JSON(response.StatusCode, string(responseBody))
	}
}

// formData := make(map[string]string)
//
//	if err := c.Bind(&formData); err != nil {
//		_ = c.Error(custom_errors.ErrInvalidFormData)
//		return
//	}
//
// jsonData, err := json.Marshal(formData)
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal form data"})
//		return
//	}
//
// // Debug: Print the JSON data
// fmt.Println("Marshalled JSON:", string(jsonData))
//
// // Initialize httptest.ResponseRecorder to capture the response
// recorder := httptest.NewRecorder()
// apiContext, _ := gin.CreateTestContext(recorder)
// apiContext.Request = c.Request
// apiContext.Request.Header.Set("Content-Type", "application/json")
// apiContext.Request.Body = io.NopCloser(strings.NewReader(string(jsonData)))
//
// // Call the API controller to create the ISO standard
// wc.ApiController.CreateISOStandard(apiContext)
//
// // Read the response from recorder
// response := recorder.Result()
// defer response.Body.Close()
// responseBody, _ := io.ReadAll(response.Body)
//
// apiStatus := response.StatusCode
//
// // fmt.Println("API Controller HTTP Status:", apiStatus)   // Debug Line
// // fmt.Println("API Response Body:", string(responseBody)) // Debug Line
//
//	if apiStatus == http.StatusCreated {
//		c.Redirect(http.StatusFound, "/web/iso_standards")
//	} else {
//
//		// c.JSON(apiContext.Writer.Status(), gin.H{"error": "Failed to create ISO standard"})
//		c.JSON(apiContext.Writer.Status(), string(responseBody))
// }
// }

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
