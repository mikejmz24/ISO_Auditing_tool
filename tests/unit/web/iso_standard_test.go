package web_test

import (
	apiController "ISO_Auditing_Tool/cmd/api/controllers"
	webController "ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
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
	formDataBytes, _ := json.Marshal(testData.FormData)
	updatedDataBytes, _ := json.Marshal(testData.UpdatedData)
	suite.formData = string(formDataBytes)
	suite.updatedData = string(updatedDataBytes)
}

func (suite *WebIsoStandardControllerTestSuite) setupMockRepo() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	if suite.mockRepo == nil {
		panic("mockRepo is nil")
	}
	fmt.Printf("Mock Repository initialized: %v\n", suite.mockRepo)
}

func (suite *WebIsoStandardControllerTestSuite) setupRouter() {
	suite.router = setupRouter(suite.mockRepo)
	if suite.router == nil {
		panic("router is nil")
	}
	fmt.Printf("Router initialized: %v\n", suite.router)
}

func (suite *WebIsoStandardControllerTestSuite) SetupTest() {
	fmt.Println("Setting up test")
	suite.setupMockRepo()
	suite.setupRouter()
	suite.loadTestData("../../testdata/iso_standards_test01.json")
	fmt.Printf("Setup complete: router=%v, mockRepo=%v, sampleData=%v\n", suite.router, suite.mockRepo, suite.standard)
}

func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
	apiController := apiController.NewApiIsoStandardController(repo)
	webController := webController.NewWebIsoStandardController(apiController)

	router := gin.Default()
	// Use relative path to the templates directory
	templatesPath := filepath.Join("..", "..", "..", "templates", "*.templ")
	router.LoadHTMLGlob(templatesPath)

	webGroup := router.Group("/web")
	{
		webGroup.GET("/iso_standards", webController.GetAllISOStandards)
		webGroup.GET("/iso_standards/:id", webController.GetISOStandardByID)
		webGroup.POST("/iso_standards", webController.CreateISOStandard)
		webGroup.PUT("/iso_standards/:id", webController.UpdateISOStandard)
		webGroup.DELETE("/iso_standards/:id", webController.DeleteISOStandard)
	}
	fmt.Printf("Router setup with routes: %v\n", webGroup)
	return router
}

func (suite *WebIsoStandardControllerTestSuite) performRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
	fmt.Printf("Performing request: method=%s, url=%s, body is nil: %v\n", method, url, body == nil)

	// Check if body is nil and handle appropriately
	if body == nil {
		fmt.Println("Body is nil, initializing empty buffer")
		body = bytes.NewBuffer([]byte{})
	} else {
		fmt.Printf("Body length: %d\n", body.Len())
	}

	fmt.Printf("Router: %v\n", suite.router)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Printf("Error creating new request: %v\n", err)
		panic(err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	fmt.Println("Serving HTTP request")
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *WebIsoStandardControllerTestSuite) validateResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBodyContains string) {
	fmt.Printf("Validating response: status=%d, body=%s\n", w.Code, w.Body.String())
	suite.Equal(expectedStatus, w.Code)
	suite.Contains(w.Body.String(), expectedBodyContains)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *WebIsoStandardControllerTestSuite) TestWebGetAllISOStandards() {
	fmt.Println("Running TestWebGetAllISOStandards")
	expectedStandards := []types.ISOStandard{suite.standard}
	suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	w := suite.performRequest("GET", "/web/iso_standards", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *WebIsoStandardControllerTestSuite) TestWebGetISOStandardByID() {
	fmt.Println("Running TestWebGetISOStandardByID")
	suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(suite.standard, nil)

	w := suite.performRequest("GET", "/web/iso_standards/1", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *WebIsoStandardControllerTestSuite) TestWebCreateISOStandard() {
	fmt.Println("Running TestWebCreateISOStandard")
	expectedID := int64(1)
	suite.mockRepo.On("CreateISOStandard", suite.standard).Return(expectedID, nil)

	w := suite.performRequest("POST", "/web/iso_standards", bytes.NewBufferString(suite.formData))
	suite.validateResponse(w, http.StatusCreated, `"id":1`)
}

func (suite *WebIsoStandardControllerTestSuite) TestWebCreateISOStandardInvalidData() {
	fmt.Println("Running TestWebCreateISOStandardInvalidData")
	invalidFormData := `{"name": ""}`

	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(int64(0), nil).Maybe()

	w := suite.performRequest("POST", "/web/iso_standards", bytes.NewBufferString(invalidFormData))
	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid data"}`)
}

func (suite *WebIsoStandardControllerTestSuite) TestWebUpdateISOStandard() {
	fmt.Println("Running TestWebUpdateISOStandard")
	updatedStandard := suite.standard
	updatedStandard.Name = "ISO 9001 Updated"
	suite.mockRepo.On("UpdateISOStandard", updatedStandard).Return(nil)

	w := suite.performRequest("PUT", "/web/iso_standards/1", bytes.NewBufferString(suite.updatedData))
	suite.validateResponse(w, http.StatusOK, `{"status":"updated"}`)
}

func (suite *WebIsoStandardControllerTestSuite) TestWebUpdateISOStandardNotFound() {
	fmt.Println("Running TestWebUpdateISOStandardNotFound")
	suite.mockRepo.On("UpdateISOStandard", suite.standard).Return(errors.New("not found"))

	w := suite.performRequest("PUT", "/web/iso_standards/2", bytes.NewBufferString(suite.formData))
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *WebIsoStandardControllerTestSuite) TestWebDeleteISOStandard() {
	fmt.Println("Running TestWebDeleteISOStandard")
	suite.mockRepo.On("DeleteISOStandard", int64(1)).Return(nil)

	fmt.Println("Mock setup complete for TestWebDeleteISOStandard")
	w := suite.performRequest("DELETE", "/web/iso_standards/1", nil)
	fmt.Println("Request performed for TestWebDeleteISOStandard")
	suite.validateResponse(w, http.StatusOK, `{"message":"ISO standard deleted"}`)
}

func (suite *WebIsoStandardControllerTestSuite) TestWebDeleteISOStandardNotFound() {
	fmt.Println("Running TestWebDeleteISOStandardNotFound")
	suite.mockRepo.On("DeleteISOStandard", int64(2)).Return(errors.New("not found"))

	fmt.Println("Mock setup complete for TestWebDeleteISOStandardNotFound")
	w := suite.performRequest("DELETE", "/web/iso_standards/2", nil)
	fmt.Println("Request performed for TestWebDeleteISOStandardNotFound")
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *WebIsoStandardControllerTestSuite) TestWebGetISOStandardByIDNotFound() {
	fmt.Println("Running TestWebGetISOStandardByIDNotFound")
	suite.mockRepo.On("GetISOStandardByID", int64(2)).Return(types.ISOStandard{}, errors.New("not found"))

	w := suite.performRequest("GET", "/web/iso_standards/2", nil)
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *WebIsoStandardControllerTestSuite) TestWebCreateISOStandardInternalServerError() {
	fmt.Println("Running TestWebCreateISOStandardInternalServerError")
	suite.mockRepo.On("CreateISOStandard", suite.standard).Return(int64(0), errors.New("internal server error"))

	w := suite.performRequest("POST", "/web/iso_standards", bytes.NewBufferString(suite.formData))
	suite.validateResponse(w, http.StatusInternalServerError, `{"error":"internal server error"}`)
}

func TestWebIsoStandardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(WebIsoStandardControllerTestSuite))
}
