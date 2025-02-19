package controllers

import (
	"ISO_Auditing_Tool/pkg/custom_errors"
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"strconv"

	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
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
		errResp := custom_errors.FailedToFetch(c.Request.Context(), "ISO Standards")
		c.JSON(errResp.StatusCode, errResp)
		return
	}
	c.JSON(http.StatusOK, isoStandards)
}

// GetISOStandardByID retrieves an ISO standard by its ID
func (cc *ApiIsoStandardController) GetISOStandardByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errResp := custom_errors.InvalidID(c.Request.Context(), "ISO Standard")
		c.JSON(errResp.StatusCode, errResp)
		return
	}
	isoStandard, err := cc.Repo.GetISOStandardByID(id)
	if err != nil {
		errResp := custom_errors.NotFound(context.TODO(), "ISO Standard")
		c.JSON(errResp.StatusCode, errResp)
		return
	}
	c.JSON(http.StatusOK, isoStandard)
}

// CreateISOStandard creates a new ISO standard
func (cc *ApiIsoStandardController) CreateISOStandard(c *gin.Context) {
	// // Read the request body
	body, err := readRequestBody(c)
	if err != nil {
		// respondWithError(c, http.StatusBadRequest, "Could not read request body")
		errResp := custom_errors.ErrInvalidJSON.ToResponse()
		c.JSON(custom_errors.ErrInvalidData.StatusCode, errResp)
		return
	}

	// Restore the body for later use
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if len(body) == 0 {
		customErr := custom_errors.EmptyData(c.Request.Context(), "JSON")
		c.JSON(customErr.StatusCode, customErr.ToResponse())
		return
	}

	rawData, err := parseJSONToMap(body)
	if err != nil {
		if customErr, ok := err.(*custom_errors.CustomError); ok {
			c.JSON(customErr.StatusCode, customErr.ToResponse())
			return
		}
		// Fallback for other types of errors
		errResp := custom_errors.ErrInvalidJSON.ToResponse()
		c.JSON(custom_errors.ErrInvalidJSON.StatusCode, errResp)
		return
	}

	if err := validateFields(c.Request.Context(), rawData); err != nil {
		if customErr, ok := err.(*custom_errors.CustomError); ok {
			c.JSON(customErr.StatusCode, customErr.ToResponse())
			return
		}
		// Fallback for non-custom errors
		errResp := custom_errors.ErrInvalidJSON.ToResponse()
		c.JSON(custom_errors.ErrInvalidJSON.StatusCode, errResp)
		return
	}

	isoStandard, err := parseISOStandard(body)
	if err != nil {
		// respondWithError(c, http.StatusBadRequest, "Could not decode JSON data")
		errResp := custom_errors.ErrInvalidJSON.ToResponse()
		c.JSON(custom_errors.ErrInvalidJSON.StatusCode, errResp)
		return
	}

	// Attempt to create ISO standard in the repository
	createdISOStandard, err := cc.Repo.CreateISOStandard(isoStandard)
	if err != nil {
		if customErr, ok := err.(*custom_errors.CustomError); ok {
			c.JSON(customErr.StatusCode, customErr.ToResponse())
			return
		}
		// Fallback for non-custom errors
		errResp := custom_errors.ErrInvalidJSON.ToResponse()
		c.JSON(custom_errors.ErrInvalidJSON.StatusCode, errResp)
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
		// _ = c.Error(custom_errors.InvalidID("ISO Standard"))
		_ = c.Error(custom_errors.InvalidID(c.Request.Context(), "ISO Standard"))
		return
	}

	isoStandard.ID = id
	if err := cc.Repo.UpdateISOStandard(isoStandard); err != nil {
		_ = c.Error(custom_errors.NotFound(c.Request.Context(), "ISO Standard"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteISOStandard deletes an ISO standard
func (cc *ApiIsoStandardController) DeleteISOStandard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(custom_errors.InvalidID(c.Request.Context(), "ISO Standard"))
		return
	}
	if err := cc.Repo.DeleteISOStandard(id); err != nil {
		_ = c.Error(custom_errors.NotFound(c.Request.Context(), "ISO Standard"))
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
		return nil, custom_errors.ErrInvalidJSON
	}
	return result, nil
}

func validateFields(ctx context.Context, data map[string]interface{}) error {
	requiredFields := []string{"name"}

	for _, field := range requiredFields {
		if err := validateField(ctx, data, field, "string"); err != nil {
			return err
		}
	}
	return nil
}

func validateField(ctx context.Context, data map[string]interface{}, field, entity string) error {
	value, ok := data[field]
	if !ok {
		// return fmt.Errorf("Missing required field: %s", field)
		return custom_errors.MissingField(ctx, field)
	}

	strVal, ok := value.(string)
	if !ok {
		return custom_errors.InvalidDataType(ctx, field, "string")
	}

	if strVal == "" {
		return custom_errors.EmptyField(ctx, entity, field)
	}

	return nil
}
func parseISOStandard(data []byte) (types.ISOStandard, error) {
	var isoStandard types.ISOStandard
	err := json.Unmarshal(data, &isoStandard)
	return isoStandard, err
}
