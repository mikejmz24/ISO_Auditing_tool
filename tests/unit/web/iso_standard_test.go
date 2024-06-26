package web_test

import (
	apiController "ISO_Auditing_Tool/cmd/api/controllers"
	webController "ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WebIsoStandardControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockRepo    *testutils.MockIsoStandardRepository
	standard    types.ISOStandard
	formData    string
	updatedData string
}

func (suite *WebIsoStandardControllerTestSuite) SetupTest() {
	fmt.Println("Setting up test")
	suite.setupMockRepo()
	suite.setupRouter()
	suite.loadTestData("../../testdata/iso_standards_test01.json")
	fmt.Printf("Setup complete: router=%v, mockRepo=%v, sampleData=%v\n", suite.router, suite.mockRepo, suite.standard)
}

func (suite *WebIsoStandardControllerTestSuite) setupMockRepo() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	fmt.Printf("Mock Repository initialized: %v\n", suite.mockRepo)
}

func (suite *WebIsoStandardControllerTestSuite) setupRouter() {
	suite.router = setupRouter(suite.mockRepo)
	fmt.Printf("Router initialized: %v\n", suite.router)
}

func (suite *WebIsoStandardControllerTestSuite) loadTestData(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load test data file: %v", err))
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("Failed to read test data: %w", err))
	}

	var testData struct {
		Standard    types.ISOStandard `json:"standard"`
		FormData    string            `json:"formData"`
		UpdatedData string            `json:"updatedData"`
	}

	if err := json.Unmarshal(data, &testData); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal test data: %v", err))
	}

	suite.standard = testData.Standard
	suite.formData = testData.FormData
	suite.updatedData = testData.UpdatedData
}

func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
	apiController := apiController.NewApiIsoStandardController(repo)
	webController := webController.NewWebIsoStandardController(apiController)

	router := gin.Default()
	templatesPath := filepath.Join("..", "..", "..", "templates", "*.templ")
	router.LoadHTMLGlob(templatesPath)

	webGroup := router.Group("/web")
	{
		webGroup.GET("/iso_standards", webController.GetAllISOStandards)
		webGroup.GET("/iso_standards/:id", webController.GetISOStandardByID)
		webGroup.GET("/iso_standards/add", webController.RenderAddISOStandardForm)
		webGroup.POST("/iso_standards", webController.CreateISOStandard)
		webGroup.PUT("/iso_standards/:id", webController.UpdateISOStandard)
		webGroup.DELETE("/iso_standards/:id", webController.DeleteISOStandard)
	}
	fmt.Printf("Router setup with routes: %v\n", webGroup)
	return router
}

func (suite *WebIsoStandardControllerTestSuite) TestGetAllISOStandards_Success() {
	suite.mockRepo.On("GetAllISOStandards").Return([]types.ISOStandard{suite.standard}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/web/iso_standards", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestGetISOStandardByID_Success() {
	suite.mockRepo.On("GetISOStandardByID", suite.standard.ID).Return(suite.standard, nil)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/web/iso_standards/%d", suite.standard.ID), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestGetISOStandardByID_NotFound() {
	suite.mockRepo.On("GetISOStandardByID", suite.standard.ID).Return(types.ISOStandard{}, fmt.Errorf("Not Found"))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/web/iso_standards/%d", suite.standard.ID), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(int64(1), nil)

	req, _ := http.NewRequest(http.MethodPost, "/web/iso_standards", bytes.NewBufferString(suite.formData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusFound, resp.Code) // Assuming a redirect after successful creation
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestCreateISOStandard_InvalidData() {
	req, _ := http.NewRequest(http.MethodPost, "/web/iso_standards", bytes.NewBufferString(`{"invalid":"data"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *WebIsoStandardControllerTestSuite) TestCreateISOStandard_InternalServerError() {
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(int64(0), fmt.Errorf("Internal Server Error"))

	req, _ := http.NewRequest(http.MethodPost, "/web/iso_standards", bytes.NewBufferString(suite.formData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestUpdateISOStandard_Success() {
	suite.mockRepo.On("UpdateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(nil)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/web/iso_standards/%d", suite.standard.ID), bytes.NewBufferString(suite.updatedData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestUpdateISOStandard_NotFound() {
	suite.mockRepo.On("UpdateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(fmt.Errorf("not found"))

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/web/iso_standards/%d", suite.standard.ID), bytes.NewBufferString(suite.updatedData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestDeleteISOStandard_Success() {
	suite.mockRepo.On("DeleteISOStandard", suite.standard.ID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/web/iso_standards/%d", suite.standard.ID), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestDeleteISOStandard_NotFound() {
	suite.mockRepo.On("DeleteISOStandard", suite.standard.ID).Return(fmt.Errorf("Not Found"))

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/web/iso_standards/%d", suite.standard.ID), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestWebIsoStandardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(WebIsoStandardControllerTestSuite))
}
