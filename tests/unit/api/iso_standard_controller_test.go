// tests/unit/api/iso_standard_controller_test.go
package api_test

import (
	"ISO_Auditing_Tool/cmd/api/controllers"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
	controller := controllers.NewApiIsoStandardController(repo)

	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/iso_standards", controller.GetAllISOStandards)
		api.POST("/iso_standards", controller.CreateISOStandard)
	}
	return router
}

func TestAPIGetAllISOStandards(t *testing.T) {
	mockRepo := new(testutils.MockIsoStandardRepository)
	expectedStandards := []types.ISOStandard{
		{ID: 1, Name: "ISO 9001"},
	}
	mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	router := setupRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/api/iso_standards", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ISO 9001")

	mockRepo.AssertExpectations(t)
}

func TestAPICreateISOStandard(t *testing.T) {
	mockRepo := new(testutils.MockIsoStandardRepository)
	newStandard := types.ISOStandard{Name: "ISO 9001"}
	expectedID := int64(1)
	mockRepo.On("CreateISOStandard", newStandard).Return(expectedID, nil)

	router := setupRouter(mockRepo)

	formData := `{"name": "ISO 9001"}`
	req, _ := http.NewRequest("POST", "/api/iso_standards", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `{"id":1}`)

	mockRepo.AssertExpectations(t)
}
