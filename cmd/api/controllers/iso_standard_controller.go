package controllers

import (
	"ISO_Auditing_Tool/pkg/custom_errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"

	"github.com/gin-gonic/gin"
)

type ApiIsoStandardController struct {
	Repo repositories.IsoStandardRepository
}

// NewApiIsoStandardController creates a new ApiIsoStandardController
func NewApiIsoStandardController(repo repositories.IsoStandardRepository) *ApiIsoStandardController {
	return &ApiIsoStandardController{Repo: repo}
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// GetAllISOStandards retrieves all ISO standards
func (cc *ApiIsoStandardController) GetAllISOStandards(c *gin.Context) {
	isoStandards, err := cc.Repo.GetAllISOStandards()
	if err != nil {
		_ = c.Error(custom_errors.ErrFailedToFetchISOStandards)
		return
	}
	c.JSON(http.StatusOK, isoStandards)
}

// GetISOStandardByID retrieves an ISO standard by its ID
func (cc *ApiIsoStandardController) GetISOStandardByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(custom_errors.ErrInvalidID)
		return
	}
	isoStandard, err := cc.Repo.GetISOStandardByID(id)
	if err != nil {
		if err.Error() == "not found" {
			respondWithError(c, http.StatusNotFound, "ISO standard not found")
		} else {
			respondWithError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, isoStandard)
}

// // CreateISOStandard creates a new ISO standard
// func (cc *ApiIsoStandardController) CreateISOStandard(c *gin.Context) {
// 	// Read the request body
// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		respondWithError(c, http.StatusBadRequest, "Could not read request body")
// 		return
// 	}
//
// 	// Parse JSON data into structured map for initial validation
// 	var rawData map[string]interface{}
// 	if err := json.Unmarshal(body, &rawData); err != nil {
// 		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
// 		return
// 	}
//
// 	// Validate required fields and check for unexpected fields
// 	requiredFields := []string{"name"}
// 	if err := validateRequiredFields(rawData, requiredFields); err != nil {
// 		respondWithError(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
//
// 	// Check if name is a string
// 	if _, ok := rawData["name"].(string); !ok {
// 		respondWithError(c, http.StatusBadRequest, "Invalid Data - name must be a string")
// 		return
// 	}
//
// 	// Additional validation for business rules
// 	if name, ok := rawData["name"].(string); !ok || name == "" {
// 		respondWithError(c, http.StatusBadRequest, "ISO Standard name should not be empty")
// 		return
// 	}
//
// 	// Parse JSON data into structured ISOStandard object
// 	var isoStandard types.ISOStandard
// 	if err := json.Unmarshal(body, &isoStandard); err != nil {
// 		respondWithError(c, http.StatusBadRequest, "Could not decode JSON data")
// 		return
// 	}
//
// 	// Attempt to create ISO standard in the repository
// 	createdISOStandard, err := cc.Repo.CreateISOStandard(isoStandard)
// 	if err != nil {
// 		respondWithError(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}
//
// 	// Respond with the created ISO standard
// 	c.JSON(http.StatusCreated, createdISOStandard)
// }

// CreateISOStandard creates a new ISO standard
func (cc *ApiIsoStandardController) CreateISOStandard(c *gin.Context) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Could not read request body")
		return
	}

	// Parse JSON data into structured map for initial validation
	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate required fields and check for unexpected fields
	requiredFields := []string{"name"}
	if err := validateRequiredFields(rawData, requiredFields); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if name is a string
	if _, ok := rawData["name"].(string); !ok {
		respondWithError(c, http.StatusBadRequest, "Invalid Data - name must be a string")
		return
	}

	// Additional validation for business rules
	if name, ok := rawData["name"].(string); !ok || name == "" {
		respondWithError(c, http.StatusBadRequest, "ISO Standard name should not be empty")
		return
	}

	// Parse JSON data into structured ISOStandard object
	var isoStandard types.ISOStandard
	if err := json.Unmarshal(body, &isoStandard); err != nil {
		// if err := c.ShouldBindJSON(&isoStandard); err != nil {
		respondWithError(c, http.StatusBadRequest, "Could not decode JSON data")
		return
	}

	// Attempt to create ISO standard in the repository
	createdISOStandard, err := cc.Repo.CreateISOStandard(isoStandard)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with the created ISO standard
	c.JSON(http.StatusCreated, createdISOStandard)
}

// validateRequiredFields validates if all required fields are present in the given data.
func validateRequiredFields(data map[string]interface{}, requiredFields []string) error {
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			return fmt.Errorf("Missing required field: %s", field)
		}
	}
	return nil
}

// UpdateISOStandard updates an existing ISO standard
func (cc *ApiIsoStandardController) UpdateISOStandard(c *gin.Context) {
	var isoStandard types.ISOStandard
	if err := c.ShouldBindJSON(&isoStandard); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid data")
		return
	}

	// Additional validation for business rules
	if isoStandard.Name == "" {
		respondWithError(c, http.StatusBadRequest, "Invalid data")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid ISO standard ID")
		return
	}

	isoStandard.ID = id
	if err := cc.Repo.UpdateISOStandard(isoStandard); err != nil {
		if err.Error() == "not found" {
			respondWithError(c, http.StatusNotFound, "ISO standard not found")
		} else {
			respondWithError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteISOStandard deletes an ISO standard
func (cc *ApiIsoStandardController) DeleteISOStandard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid ISO standard ID")
		return
	}
	if err := cc.Repo.DeleteISOStandard(id); err != nil {
		if err.Error() == "not found" {
			respondWithError(c, http.StatusNotFound, "ISO standard not found")
		} else {
			respondWithError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ISO standard deleted"})
}
