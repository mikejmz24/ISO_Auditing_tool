// package web_test
//
// import (
//
//	apiController "ISO_Auditing_Tool/cmd/api/controllers"
//	webController "ISO_Auditing_Tool/cmd/web/controllers"
//	"ISO_Auditing_Tool/pkg/custom_errors"
//	"ISO_Auditing_Tool/pkg/middleware"
//	"ISO_Auditing_Tool/pkg/types"
//	"ISO_Auditing_Tool/tests/testutils"
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"net/url"
//	"os"
//	"path/filepath"
//	"strings"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/require"
//	"github.com/stretchr/testify/suite"
//
// )
//
//	type WebIsoStandardControllerTestSuite struct {
//		suite.Suite
//		router   *gin.Engine
//		mockRepo *testutils.MockIsoStandardRepository
//		standard types.ISOStandard
//	}
//
//	type MarshallingTestSuite struct {
//		suite.Suite
//		standard types.ISOStandard
//	}
//
//	type ValidationTestSuite struct {
//		suite.Suite
//		router *gin.Engine
//	}
//
//	type PersistenceTestSuite struct {
//		suite.Suite
//		mockRepo *testutils.MockIsoStandardRepository
//		standard types.ISOStandard
//	}
//
// var testStandard types.ISOStandard
// var testJSONData []byte
//
//	func TestMain(m *testing.M) {
//		// Load the test data once
//		testStandard = loadTestData("../../testdata/iso_standards_test01.json")
//		testJSONData, _ = json.Marshal(testStandard)
//
//		// Run the tests
//		code := m.Run()
//		os.Exit(code)
//	}
//
//	func loadTestData(filePath string) types.ISOStandard {
//		data := getJSONData(filePath)
//
//		var jsonData struct {
//			ISOStandards []types.ISOStandard `json:"iso_standards"`
//		}
//		if err := json.Unmarshal(data, &jsonData); err != nil {
//			panic(fmt.Sprintf("Failed to unmarshal test data: %v", err))
//		}
//
//		if len(jsonData.ISOStandards) == 0 {
//			panic("No ISO standards found in test data")
//		}
//
//		return jsonData.ISOStandards[0]
//	}
//
//	func getJSONData(filePath string) []byte {
//		file, err := os.Open(filePath)
//		if err != nil {
//			panic(fmt.Sprintf("Failed to load data file: %v", err))
//		}
//		defer file.Close()
//		data, err := io.ReadAll(file)
//		if err != nil {
//			panic(fmt.Errorf("Failed to read data: %w", err))
//		}
//		return data
//	}
//
//	func (suite *WebIsoStandardControllerTestSuite) SetupTest() {
//		suite.setupMockRepo()
//		suite.setupRouter()
//		suite.standard = testStandard
//	}
//
//	func (suite *WebIsoStandardControllerTestSuite) setupMockRepo() {
//		suite.mockRepo = new(testutils.MockIsoStandardRepository)
//		if suite.mockRepo == nil {
//			panic("mockRepo is nil")
//		}
//	}
//
//	func (suite *WebIsoStandardControllerTestSuite) setupRouter() {
//		suite.router = setupRouter(suite.mockRepo)
//		if suite.router == nil {
//			panic("router is nil")
//		}
//	}
//
//	func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
//		apiController := apiController.NewApiIsoStandardController(repo)
//		webController := webController.NewWebIsoStandardController(apiController)
//
//		router := gin.Default()
//		router.Use(middleware.ErrorHandler())
//		templatesPath := filepath.Join("..", "..", "..", "templates", "*.templ")
//		router.LoadHTMLGlob(templatesPath)
//
//		webGroup := router.Group("/web")
//		{
//			webGroup.GET("/iso_standards", webController.GetAllISOStandards)
//			webGroup.GET("/iso_standards/:id", webController.GetISOStandardByID)
//			webGroup.GET("/iso_standards/add", webController.RenderAddISOStandardForm)
//			webGroup.POST("/iso_standards/add", webController.CreateISOStandard)
//			webGroup.PUT("/iso_standards/:id", webController.UpdateISOStandard)
//			webGroup.DELETE("/iso_standards/:id", webController.DeleteISOStandard)
//		}
//		return router
//	}
//
//	func (suite *WebIsoStandardControllerTestSuite) validateErrorResponse(w *httptest.ResponseRecorder, expectedError *custom_errors.CustomError) {
//		assert.Equal(suite.T(), expectedError.StatusCode, w.Code)
//
//		// Print response bodoy for debugging
//		// suite.T().Logf("Resonse body: %s", w.Body.String())
//		fmt.Printf("Response body: %s\n", w.Body.String())
//
//		var response map[string]interface{}
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		require.NoError(suite.T(), err, "failed to unmarshal response body")
//
//		assert.Equal(suite.T(), expectedError.Message, response["error"])
//		if expectedError.Context != nil {
//			assert.Equal(suite.T(), expectedError.Context, response["context"])
//		}
//	}
//
// // Marshalling Tests
//
//	func (suite *MarshallingTestSuite) SetupSuite() {
//		suite.standard = testStandard
//	}
//
//	func (suite *MarshallingTestSuite) TestMarshalISOStandard() {
//		expectedJSON := testJSONData
//		actualJSON, err := json.Marshal(suite.standard)
//		suite.NoError(err)
//		var expectedData interface{}
//		var actualData interface{}
//
//		err = json.Unmarshal(expectedJSON, &expectedData)
//		suite.NoError(err)
//		err = json.Unmarshal(actualJSON, &actualData)
//		suite.NoError(err)
//
//		suite.Equal(expectedData, actualData)
//	}
//
// // Validation Tests
//
//	func (suite *ValidationTestSuite) SetupSuite() {
//		suite.router = gin.Default()
//		suite.router.POST("/web/iso_standards", func(c *gin.Context) {
//			var data map[string]interface{}
//			if err := c.ShouldBindJSON(&data); err != nil {
//				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
//				return
//			}
//			c.Status(http.StatusOK)
//		})
//	}
//
//	func (suite *WebIsoStandardControllerTestSuite) performRequest(method, url string, body io.Reader) *httptest.ResponseRecorder {
//		req, err := http.NewRequest(method, url, body)
//		if err != nil {
//			suite.T().Fatalf("failed to create request: %v", err)
//		}
//		req.Header.Set("Content-Type", "application/json")
//
//		w := httptest.NewRecorder()
//		suite.router.ServeHTTP(w, req)
//		return w
//	}
//
//	func (suite *WebIsoStandardControllerTestSuite) TestCreateISOStandard_InvalidData() {
//		invalidJSON := `"invalidField": "invalidData"`
//
//		// req, _ := http.NewRequest(http.MethodPost, "/web/iso_standards/add", bytes.NewBuffer([]byte(invalidJSON)))
//		w := suite.performRequest("POST", "/web/iso_standards/add", bytes.NewBuffer([]byte(invalidJSON)))
//		// req.Header.Set("Content-Type", "application/json")
//		// resp := httptest.NewRecorder()
//		// suite.router.ServeHTTP(resp, req)
//
//		// suite.Equal(http.StatusBadRequest, resp.Code)
//		// suite.Contains(resp.Body.String(), "Invalid data")
//		// suite.Equal(resp.Body, custom_errors.ErrInvalidFormData)
//		suite.validateErrorResponse(w, custom_errors.ErrInvalidFormData)
//		// suite.mockRepo.AssertExpectareqtions(suite.T())
//	}
//
// // Persistence Tests
//
//	func (suite *PersistenceTestSuite) SetupSuite() {
//		suite.mockRepo = new(testutils.MockIsoStandardRepository)
//		suite.standard = testStandard
//	}
//
//	func (suite *WebIsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
//		standard := suite.standard
//
//		formData := url.Values{
//			"name": {standard.Name},
//		}
//
//		suite.mockRepo.On("CreateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(standard, nil)
//		suite.mockRepo.On("GetAllISOStandards").Return([]types.ISOStandard{standard}, nil)
//
//		req, err := http.NewRequest(http.MethodPost, "/web/iso_standards/add", strings.NewReader(formData.Encode()))
//		suite.NoError(err)
//		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//		resp := httptest.NewRecorder()
//		suite.router.ServeHTTP(resp, req)
//
//		suite.Equal(http.StatusFound, resp.Code)
//		location := resp.Header().Get("Location")
//		suite.Equal("/web/iso_standards", location)
//
//		req, err = http.NewRequest(http.MethodGet, location, nil)
//		suite.NoError(err)
//
//		resp = httptest.NewRecorder()
//		suite.router.ServeHTTP(resp, req)
//
//		suite.Equal(http.StatusOK, resp.Code)
//		suite.Contains(resp.Body.String(), "ISO 9001")
//
//		suite.mockRepo.AssertExpectations(suite.T())
//	}
//
//	func TestWebIsoStandardControllerTestSuite(t *testing.T) {
//		suite.Run(t, new(WebIsoStandardControllerTestSuite))
//		suite.Run(t, new(MarshallingTestSuite))
//		suite.Run(t, new(ValidationTestSuite))
//		suite.Run(t, new(PersistenceTestSuite))
//	}
package web_test

import (
	apiController "ISO_Auditing_Tool/cmd/api/controllers"
	webController "ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/middleware"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestSuite combines all test suites
type TestSuite struct {
	suite.Suite
	router   *gin.Engine
	mockRepo *testutils.MockIsoStandardRepository
	standard types.ISOStandard
}

var (
	testStandard types.ISOStandard
	testJSONData []byte
)

func TestMain(m *testing.M) {
	testStandard = loadTestData("../../testdata/iso_standards_test01.json")
	testJSONData, _ = json.Marshal(testStandard)
	os.Exit(m.Run())
}

func loadTestData(filePath string) types.ISOStandard {
	data := getJSONData(filePath)

	var jsonData struct {
		ISOStandards []types.ISOStandard `json:"iso_standards"`
	}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal test data: %v", err))
	}

	if len(jsonData.ISOStandards) == 0 {
		panic("No ISO standards found in test data")
	}

	return jsonData.ISOStandards[0]
}

func getJSONData(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load data file: %v", err))
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("Failed to read data: %w", err))
	}
	return data
}

func (suite *TestSuite) SetupTest() {
	suite.setupMockRepo()
	suite.setupRouter()
	suite.standard = testStandard
}

func (suite *TestSuite) setupMockRepo() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	if suite.mockRepo == nil {
		panic("mockRepo is nil")
	}
}

func (suite *TestSuite) setupRouter() {
	suite.router = setupRouter(suite.mockRepo)
	if suite.router == nil {
		panic("router is nil")
	}
}

func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
	apiController := apiController.NewApiIsoStandardController(repo)
	webController := webController.NewWebIsoStandardController(apiController)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
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

func (suite *TestSuite) validateErrorResponse(w *httptest.ResponseRecorder, expectedError *custom_errors.CustomError) {
	assert.Equal(suite.T(), expectedError.StatusCode, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err, "failed to unmarshal response body")

	assert.Equal(suite.T(), expectedError.Message, response["error"])
	if expectedError.Context != nil {
		assert.Equal(suite.T(), expectedError.Context, response["context"])
	}
}

func (suite *TestSuite) performRequest(method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		suite.T().Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

// Test cases

func (suite *TestSuite) TestMarshalISOStandard() {
	expectedJSON := testJSONData
	actualJSON, err := json.Marshal(suite.standard)
	suite.NoError(err)
	var expectedData, actualData interface{}

	suite.NoError(json.Unmarshal(expectedJSON, &expectedData))
	suite.NoError(json.Unmarshal(actualJSON, &actualData))

	suite.Equal(expectedData, actualData)
}

func (suite *TestSuite) TestCreateISOStandard_InvalidData() {
	invalidJSON := `"invalidField": "invalidData"`
	w := suite.performRequest("POST", "/web/iso_standards/add", bytes.NewBuffer([]byte(invalidJSON)))
	suite.validateErrorResponse(w, custom_errors.ErrInvalidFormData)
}

func (suite *TestSuite) TestCreateISOStandard_Success() {
	standard := suite.standard

	formData := url.Values{
		"name": {standard.Name},
	}

	suite.mockRepo.On("CreateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(standard, nil)
	suite.mockRepo.On("GetAllISOStandards").Return([]types.ISOStandard{standard}, nil)

	req, err := http.NewRequest(http.MethodPost, "/web/iso_standards/add", strings.NewReader(formData.Encode()))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusFound, resp.Code)
	location := resp.Header().Get("Location")
	suite.Equal("/web/iso_standards", location)

	req, err = http.NewRequest(http.MethodGet, location, nil)
	suite.NoError(err)

	resp = httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "ISO 9001")

	suite.mockRepo.AssertExpectations(suite.T())
}

func TestISOStandardSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
