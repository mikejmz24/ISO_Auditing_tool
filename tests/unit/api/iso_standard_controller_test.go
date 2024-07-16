package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ISO_Auditing_Tool/cmd/api/controllers"
	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/middleware"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IsoStandardControllerTestSuite struct {
	suite.Suite
	router   *gin.Engine
	mockRepo *testutils.MockIsoStandardRepository
	jsonData []types.ISOStandard
}

func (suite *IsoStandardControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	filePath := "../../testdata/iso_standards_test01.json"
	file, err := os.Open(filePath)
	require.NoError(suite.T(), err, "failed to open JSON file")
	defer file.Close()

	data, err := io.ReadAll(file)
	require.NoError(suite.T(), err, "failed to read JSON file")

	var jsonData struct {
		ISOStandards []types.ISOStandard `json:"iso_standards"`
	}
	err = json.Unmarshal(data, &jsonData)
	require.NoError(suite.T(), err, "failed to unmarshal JSON data")
	require.NotEmpty(suite.T(), jsonData.ISOStandards, "no ISO standards found in JSON data")

	suite.jsonData = jsonData.ISOStandards
	fmt.Printf("Loaded ISO standards: %v\n", suite.jsonData)
}

func (suite *IsoStandardControllerTestSuite) SetupTest() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	suite.router = gin.Default()
	suite.router.Use(middleware.ErrorHandler())
	controller := controllers.NewApiIsoStandardController(suite.mockRepo)
	api := suite.router.Group("/api")
	{
		api.GET("/iso_standards", controller.GetAllISOStandards)
		api.GET("/iso_standards/:id", controller.GetISOStandardByID)
		api.POST("/iso_standards", controller.CreateISOStandard)
		api.PUT("/iso_standards/:id", controller.UpdateISOStandard)
		api.DELETE("/iso_standards/:id", controller.DeleteISOStandard)
	}
}

func (suite *IsoStandardControllerTestSuite) performRequest(method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		suite.T().Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *IsoStandardControllerTestSuite) validateResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBodyContains string) {
	assert.Equal(suite.T(), expectedStatus, w.Code)
	assert.Contains(suite.T(), w.Body.String(), expectedBodyContains)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) validateErrorResponse(w *httptest.ResponseRecorder, expectedError *custom_errors.CustomError) {
	assert.Equal(suite.T(), expectedError.StatusCode, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err, "failed to unmarshal response body")
	assert.Equal(suite.T(), expectedError.Message, response["error"])
	if expectedError.Context != nil {
		assert.Equal(suite.T(), expectedError.Context, response["context"])
	}
}

func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Success() {
	expectedStandard := suite.jsonData[0]
	expectedStandards := []types.ISOStandard{expectedStandard}
	suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	w := suite.performRequest("GET", "/api/iso_standards", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Error() {
	suite.mockRepo.On("GetAllISOStandards").Return(nil, custom_errors.NewCustomError(http.StatusInternalServerError, "Failed to fetch ISO Standards", nil))

	w := suite.performRequest("GET", "/api/iso_standards", nil)
	suite.validateErrorResponse(w, custom_errors.NewCustomError(http.StatusInternalServerError, "Failed to fetch ISO Standards", nil))
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Success() {
	expectedStandard := suite.jsonData[0]
	suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(expectedStandard, nil)

	w := suite.performRequest("GET", "/api/iso_standards/1", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Error() {
	suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(types.ISOStandard{}, custom_errors.NewCustomError(http.StatusBadRequest, "Invalid ISO ID", nil))

	w := suite.performRequest("GET", "/api/iso_standards/1", nil)
	suite.validateErrorResponse(w, custom_errors.NewCustomError(http.StatusBadRequest, "Invalid ISO ID", nil))
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
	// Prepare expected data
	expectedStandard := suite.jsonData[0]
	formData, err := json.Marshal(expectedStandard)
	require.NoError(suite.T(), err, "failed to marshal JSON data")

	// Mock repository call and perform request
	suite.mockRepo.On("CreateISOStandard", expectedStandard).Return(expectedStandard, nil)
	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBuffer(formData))

	// Validate response
	suite.validateResponse(w, http.StatusCreated, `"name":"Leadership"`)
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_EmptyName() {
	// Prepare invalid form data
	invalidFormData := `{"name": ""}`
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()

	// Perform request
	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))

	// Validate response for empty name error
	suite.validateResponse(w, http.StatusBadRequest, `{"error":"ISO Standard name should not be empty"}`)
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_InvalidData_EmptyJSON() {
	// Prepare invalid form data (empty JSON)
	invalidFormData := ``
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()

	// Perform request
	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))

	// Validate response for invalid JSON format
	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid JSON format"}`)
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_InvalidJSON_WrongFieldName() {
	// Prepare invalid form data (incorrect field)
	invalidFormData := `{"nam": "ISO fake name"}`
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()

	// Perform request
	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))

	// Validate response for invalid field name
	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Missing required field: name"}`)
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_InvalidDataType_BoolInsteadOfString() {
	invalidFormData := `{"name": true}` // name should be a string, not a boolean
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()

	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))
	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid Data - name must be a string"}`)
}

// func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_AdditionalErrors() {
// 	// Additional tests for other potential errors
// 	// For example:
// 	// - Test for missing required fields other than 'name'
// 	// - Test for unexpected additional fields
//
// 	// Mock repository call with expected response
// 	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()
//
// 	// Perform request with additional error scenarios
// 	// Adjust as per specific additional error cases you want to cover
//
// 	// Example: Missing required field
// 	missingFieldData := `{}` // Missing 'name' field
// 	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(missingFieldData))
// 	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid field name: expected 'name'"}`)
//
// 	// Example: Unexpected additional fields
// 	unexpectedFieldData := `{"name": "ISO 9001", "version": "1.0"}`
// 	w = suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(unexpectedFieldData))
// 	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid field name: expected 'name'"}`)
//
// 	// Add more test cases as necessary for other potential errors
// }

func (suite *IsoStandardControllerTestSuite) TestUpdateISOStandard_Success() {
	expectedStandard := suite.jsonData[0]
	updatedStandard := expectedStandard
	updatedStandard.Name = "ISO 9001 Updated"

	updatedData, err := json.Marshal(updatedStandard)
	require.NoError(suite.T(), err, "failed to marshal JSON data")

	suite.mockRepo.On("UpdateISOStandard", mock.MatchedBy(func(isoStandard types.ISOStandard) bool {
		return isoStandard.ID == updatedStandard.ID &&
			isoStandard.Name == updatedStandard.Name
	})).Return(nil)

	w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(updatedData))
	suite.validateResponse(w, http.StatusOK, `{"status":"updated"}`)
}

func (suite *IsoStandardControllerTestSuite) TestUpdateISOStandard_NotFound() {
	expectedStandard := suite.jsonData[0]

	formData, err := json.Marshal(expectedStandard)
	require.NoError(suite.T(), err, "failed to marshal JSON data")

	suite.mockRepo.On("UpdateISOStandard", expectedStandard).Return(custom_errors.NewCustomError(http.StatusNotFound, "ISO standard not found", nil))

	w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(formData))
	suite.validateErrorResponse(w, custom_errors.NewCustomError(http.StatusNotFound, "ISO standard not found", nil))
}

func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_Success() {
	suite.mockRepo.On("DeleteISOStandard", int64(1)).Return(nil)

	w := suite.performRequest("DELETE", "/api/iso_standards/1", nil)
	suite.validateResponse(w, http.StatusOK, `{"message":"ISO standard deleted"}`)
}

func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_NotFound() {
	suite.mockRepo.On("DeleteISOStandard", int64(2)).Return(custom_errors.NewCustomError(http.StatusNotFound, "ISO standard not found", nil))

	w := suite.performRequest("DELETE", "/api/iso_standards/2", nil)
	suite.validateErrorResponse(w, custom_errors.NewCustomError(http.StatusNotFound, "ISO standard not found", nil))
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_NotFound() {
	suite.mockRepo.On("GetISOStandardByID", int64(2)).Return(types.ISOStandard{}, custom_errors.NewCustomError(http.StatusNotFound, "ISO standard not found", nil))

	w := suite.performRequest("GET", "/api/iso_standards/2", nil)
	suite.validateErrorResponse(w, custom_errors.NewCustomError(http.StatusNotFound, "ISO standard not found", nil))
}

func TestApiIsoStandardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(IsoStandardControllerTestSuite))
}
