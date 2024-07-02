package controllers

import (
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
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, isoStandards)
}

// GetISOStandardByID retrieves an ISO standard by its ID
func (cc *ApiIsoStandardController) GetISOStandardByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid ISO standard ID")
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
	createdISOStandard, err := cc.Repo.CreateISOStandard(isoStandard)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
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
