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
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WebIsoStandardControllerTestSuite struct {
	suite.Suite
	router   *gin.Engine
	mockRepo *testutils.MockIsoStandardRepository
	standard types.ISOStandard
	// formData    string
	// updatedData string
}

type MarshallingTestSuite struct {
	suite.Suite
	standard types.ISOStandard
}

type ValidationTestSuite struct {
	suite.Suite
	router *gin.Engine
}

type PersistenceTestSuite struct {
	suite.Suite
	mockRepo *testutils.MockIsoStandardRepository
	standard types.ISOStandard
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

	var testData types.ISOStandard
	if err := json.Unmarshal(data, &testData); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal test data: %v", err))
	}

	suite.standard = testData
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
		webGroup.POST("/iso_standards/add", webController.CreateISOStandard)
		webGroup.PUT("/iso_standards/:id", webController.UpdateISOStandard)
		webGroup.DELETE("/iso_standards/:id", webController.DeleteISOStandard)
	}
	return router
}

// Marshalling Tests

func (suite *MarshallingTestSuite) SetupSuite() {
	suite.standard = types.ISOStandard{
		ID:   1,
		Name: "ISO 9001",
	}
}

func (suite *MarshallingTestSuite) TestMarshalISOStandard() {
	expectedJSON := `{"id":1,"name":"ISO 9001"}`

	data, err := json.Marshal(suite.standard)
	suite.NoError(err)
	suite.JSONEq(expectedJSON, string(data))
}

// Validation Tests

func (suite *ValidationTestSuite) SetupSuite() {
	suite.router = gin.Default()
	suite.router.POST("/web/iso_standards", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}
		c.Status(http.StatusOK)
	})
}

func (suite *WebIsoStandardControllerTestSuite) TestCreateISOStandard_InvalidData() {
	invalidJSON := `{"invalidField": "invalidData"}`

	req, _ := http.NewRequest(http.MethodPost, "/web/iso_standards/add", bytes.NewBuffer([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Contains(resp.Body.String(), "Invalid data")
	suite.mockRepo.AssertExpectations(suite.T())
}

// Persistence Tests

func (suite *PersistenceTestSuite) SetupSuite() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	suite.standard = types.ISOStandard{
		ID:   1,
		Name: "ISO 9001",
	}
}

func (suite *WebIsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
	// Prepare the test data
	standard := types.ISOStandard{
		ID:   1,
		Name: "ISO 9001",
	}

	// If the handler expects form data, use url.Values
	formData := url.Values{
		"name": {standard.Name},
	}

	// Set up the mock repository
	suite.mockRepo.On("CreateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(standard, nil)
	suite.mockRepo.On("GetAllISOStandards").Return([]types.ISOStandard{standard}, nil) // Mock GetAllISOStandards method

	// Perform the POST request to create the ISO standard
	req, err := http.NewRequest(http.MethodPost, "/web/iso_standards/add", strings.NewReader(formData.Encode()))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	// Validate the response code and location header for redirect
	suite.Equal(http.StatusFound, resp.Code)
	location := resp.Header().Get("Location")
	suite.Equal("/web/iso_standards", location)

	// Perform a GET request to the redirected URL
	req, err = http.NewRequest(http.MethodGet, location, nil)
	suite.NoError(err)

	resp = httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	// Validate the response contains the newly created ISO standard
	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "ISO 9001")

	// Ensure the mock expectations are met
	suite.mockRepo.AssertExpectations(suite.T())
	suite.Equal("/web/iso_standards", location)

	// Perform a GET request to the redirected URL
	req, err = http.NewRequest(http.MethodGet, location, nil)
	suite.NoError(err)

	resp = httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	// Validate the response contains the newly created ISO standard
	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "ISO 9001")

	// Ensure the mock expectations are met
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestWebIsoStandardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(WebIsoStandardControllerTestSuite))
	suite.Run(t, new(MarshallingTestSuite))
	suite.Run(t, new(ValidationTestSuite))
	suite.Run(t, new(PersistenceTestSuite))
}
