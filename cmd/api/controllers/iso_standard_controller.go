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

// CreateISOStandard creates a new ISO standard
func (cc *ApiIsoStandardController) CreateISOStandard(c *gin.Context) {
	// // Read the request body
	body, err := readRequestBody(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Could ont read request body")
		return
	}

	rawData, err := parseJSONToMap(body)
	if err != nil {
		_ = c.Error(custom_errors.ErrInvalidJSON)
		return
	}

	if err := validateFields(rawData); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	isoStandard, err := parseISOStandard(body)
	if err != nil {
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

func readRequestBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseJSONToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func validateFields(data map[string]interface{}) error {
	requiredFields := []string{"name"}

	for _, field := range requiredFields {
		if err := validateField(data, field, "ISO Standard"); err != nil {
			return err
		}
	}
	return nil
}

func validateField(data map[string]interface{}, field, entity string) error {
	value, ok := data[field]
	if !ok {
		return fmt.Errorf("Missing required field: %s", field)
	}

	strVal, ok := value.(string)
	if !ok {
		return custom_errors.InvalidDataType(field, "string")
	}

	if strVal == "" {
		return custom_errors.EmptyField(entity, field)
	}

	return nil
}
func parseISOStandard(data []byte) (types.ISOStandard, error) {
	var isoStandard types.ISOStandard
	err := json.Unmarshal(data, &isoStandard)
	return isoStandard, err
}
