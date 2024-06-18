package web

import (
	"ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(mockRepo *testutils.MockIsoStandardRepository) *gin.Engine {
	controller := controllers.NewHtmlIsoStandardController(mockRepo)
	router := gin.Default()
	html := router.Group("/html")
	{
		html.GET("/iso_standards", controller.GetAllISOStandards)
		html.GET("/iso_standards/add", controller.RenderAddISOStandardForm)
		html.POST("/iso_standards/add", controller.CreateISOStandard)
	}
	return router
}

func TestGetAllISOStandards(t *testing.T) {
	mockRepo := new(testutils.MockIsoStandardRepository)
	mockRepo.On("GetAllISOStandards").Return([]types.ISOStandard{
		// {ID: 1, Name: "ISO 9001", Description: "Quality Management"},
		{ID: 1, Name: "ISO 9001"},
	}, nil)

	router := setupRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/html/iso_standards", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ISO 9001")

	mockRepo.AssertExpectations(t)
}

func TestRenderAddISOStandardForm(t *testing.T) {
	mockRepo := new(testutils.MockIsoStandardRepository)
	router := setupRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/html/iso_standards/add", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Add ISO Standard")
}

func TestCreateISOStandard(t *testing.T) {
	mockRepo := new(testutils.MockIsoStandardRepository)
	mockRepo.On("CreateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(int64(1), nil)

	router := setupRouter(mockRepo)

	formData := "name=ISO 9001"
	req, _ := http.NewRequest("POST", "/html/iso_standards/add", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code) // Check for redirect
	assert.Equal(t, "/html/iso_standards", w.Header().Get("Location"))

	mockRepo.AssertExpectations(t)
}
