package api_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func (suite *IsoStandardControllerTestSuite) SetupSuite() {
	suite.jsonData = loadTestData(suite.T(), "../../testdata/iso_standards_test01.json")
	require.NotEmpty(suite.T(), suite.jsonData, "Test data is empty")
}

func (suite *IsoStandardControllerTestSuite) SetupTest() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	suite.router = setupRouter(suite.mockRepo)
}

func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	controller := controllers.NewApiIsoStandardController(repo)
	api := router.Group("/api")
	{
		api.GET("/iso_standards", controller.GetAllISOStandards)
		api.GET("/iso_standards/:id", controller.GetISOStandardByID)
		api.POST("/iso_standards", controller.CreateISOStandard)
		api.PUT("/iso_standards/:id", controller.UpdateISOStandard)
		api.DELETE("/iso_standards/:id", controller.DeleteISOStandard)
	}
	return router
}

func loadTestData(t *testing.T, filePath string) []types.ISOStandard {
	file, err := os.Open(filePath)
	require.NoError(t, err, "failed to open JSON file")
	defer file.Close()

	data, err := io.ReadAll(file)
	require.NoError(t, err, "failed to read JSON file")

	var jsonData struct {
		ISOStandards []types.ISOStandard `json:"iso_standards"`
	}
	err = json.Unmarshal(data, &jsonData)
	require.NoError(t, err, "failed to unmarshal JSON data")
	require.NotEmpty(t, jsonData.ISOStandards, "no ISO standards found in JSON data")

	return jsonData.ISOStandards
}

func (suite *IsoStandardControllerTestSuite) performRequest(method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	require.NoError(suite.T(), err, "failed to create request")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *IsoStandardControllerTestSuite) validateErrorResponse(w *httptest.ResponseRecorder, expectedError custom_errors.CustomError) {
	// Compare HTTP status codes
	assert.Equal(suite.T(), expectedError.StatusCode, w.Code, "HTTP status code does not match expected")

	responseBody := w.Body.String()
	var errorResponse custom_errors.ErrorResponse
	err := json.Unmarshal([]byte(responseBody), &errorResponse)
	assert.NoError(suite.T(), err, "Failed to unmarshal error response")

	// Compare error code and message
	assert.Equal(suite.T(), expectedError.Code, errorResponse.Code, "Error code does not match expected")
	assert.Equal(suite.T(), expectedError.Message, errorResponse.Message, "Error message does not match expected")
}

func (suite *IsoStandardControllerTestSuite) validateResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBodyContains string) {
	assert.Equal(suite.T(), expectedStatus, w.Code, "HTTP status code does not match expected")

	responseBody := w.Body.String()
	if expectedStatus >= 400 { // For error responses
		var errorResponse struct {
			Error string `json:"error"`
		}
		err := json.Unmarshal([]byte(responseBody), &errorResponse)
		assert.NoError(suite.T(), err, "Failed to unmarshal error response")
		assert.NotEmpty(suite.T(), errorResponse.Error, "Error message should not be empty")
		assert.Equal(suite.T(), expectedBodyContains, errorResponse.Error, "Error message does not match expected")
	} else if expectedBodyContains != "" {
		// Flexible validation for success responses
		assert.Contains(suite.T(), responseBody, expectedBodyContains, "Response body does not contain expected content")
	}
}

// Test cases

func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Success() {
	testCases := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "ReturnOneSuccess",
			setupMock: func() {
				suite.mockRepo.On("GetAllISOStandards").Return([]types.ISOStandard{suite.jsonData[0]}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"ISO 9001"`,
		},
		{
			name: "ReturnManySuccess",
			setupMock: func() {
				suite.mockRepo.On("GetAllISOStandards").Return(suite.jsonData, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"ISO 27001"`,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.mockRepo = new(testutils.MockIsoStandardRepository)
			suite.router = setupRouter(suite.mockRepo)
			tc.setupMock()
			w := suite.performRequest("GET", "/api/iso_standards", nil)
			suite.validateResponse(w, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Validation() {
	testCases := []struct {
		name          string
		setupMock     func()
		expectedError custom_errors.CustomError
	}{
		{
			name: "CannotFetchError",
			setupMock: func() {
				suite.mockRepo.On("GetAllISOStandards").Return(nil, custom_errors.FailedToFetch(context.TODO(), "ISO Standards"))
			},
			expectedError: *custom_errors.FailedToFetch(context.TODO(), "ISO Standards"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.mockRepo = new(testutils.MockIsoStandardRepository)
			suite.router = setupRouter(suite.mockRepo)
			tc.setupMock()
			w := suite.performRequest("GET", "/api/iso_standards", nil)
			suite.validateErrorResponse(w, tc.expectedError)
		})
	}
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Success() {
	testCases := []struct {
		name           string
		id             string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			id:   "1",
			setupMock: func() {
				suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(suite.jsonData[0], nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"ISO 9001"`,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.mockRepo = new(testutils.MockIsoStandardRepository)
			suite.router = setupRouter(suite.mockRepo)
			tc.setupMock()
			w := suite.performRequest("GET", "/api/iso_standards/"+tc.id, nil)
			suite.validateResponse(w, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Validation() {
	testCases := []struct {
		name          string
		id            string
		setupMock     func()
		expectedError custom_errors.CustomError
	}{
		{
			name: "CannotParseIDAsIntError",
			id:   "a",
			setupMock: func() {
				// suite.mockRepo.On("GetISOStandardByID", mock.Anything).Return(types.ISOStandard{}, custom_errors.InvalidID("ISO Standard"))
			},
			expectedError: *custom_errors.InvalidID(context.TODO(), "ISO Standard"),
		},
		{
			name: "NotFound",
			id:   "2",
			setupMock: func() {
				suite.mockRepo.On("GetISOStandardByID", int64(2)).Return(types.ISOStandard{}, custom_errors.NotFound(context.TODO(), "ISO Standard"))
			},
			expectedError: *custom_errors.NotFound(context.TODO(), "ISO Standard"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.mockRepo = new(testutils.MockIsoStandardRepository)
			suite.router = setupRouter(suite.mockRepo)
			tc.setupMock()
			w := suite.performRequest("GET", "/api/iso_standards/"+tc.id, nil)
			suite.validateErrorResponse(w, tc.expectedError)
		})
	}
}

// TODO: Update expectedStatus and ExpectedBody to be a *custom_error
func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_Validation() {
	testCases := []struct {
		name          string
		body          string
		setupMock     func()
		expectedError custom_errors.CustomError
	}{
		{
			name:          "EmptyJSON",
			body:          ``,
			setupMock:     func() {},
			expectedError: *custom_errors.EmptyData(context.TODO(), "JSON"),
		},
		{
			name:          "FieldNameMisspelled",
			body:          `{"nam": "ISO fake name"}`,
			setupMock:     func() {},
			expectedError: *custom_errors.MissingField(context.TODO(), "name"),
		},
		{
			name:          "EmptyName",
			body:          `{"name": ""}`,
			setupMock:     func() {},
			expectedError: *custom_errors.EmptyField(context.TODO(), "string", "name"),
		},
		{
			name:          "BoolInsteadOfString",
			body:          `{"name": true}`,
			setupMock:     func() {},
			expectedError: *custom_errors.InvalidDataType(context.TODO(), "name", "string"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setupMock()
			w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(tc.body))
			suite.validateErrorResponse(w, tc.expectedError)
		})
	}
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
	testCases := []struct {
		name           string
		body           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			body: `{"name":"ISO 9001"}`,
			setupMock: func() {
				suite.mockRepo.On("CreateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(suite.jsonData[0], nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"name":"Leadership"`,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setupMock()
			w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(tc.body))
			suite.validateResponse(w, tc.expectedStatus, tc.expectedBody)
		})
	}
}
func (suite *IsoStandardControllerTestSuite) TestUpdateISOStandard() {
	testCases := []struct {
		name           string
		id             string
		body           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			id:   "1",
			body: `{"name":"ISO 9001 Updated"}`,
			setupMock: func() {
				suite.mockRepo.On("UpdateISOStandard", mock.MatchedBy(func(isoStandard types.ISOStandard) bool {
					return isoStandard.ID == 1 && isoStandard.Name == "ISO 9001 Updated"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"updated"}`,
		},
		{
			name:           "CannotParseIDAsIntError",
			id:             "a",
			body:           `{"name":"ISO 9001 Updated"}`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid ISO Standard ID",
		},
		{
			name: "NotFound",
			id:   "2",
			body: `{"name":"ISO 9001 Updated"}`,
			setupMock: func() {
				suite.mockRepo.On("UpdateISOStandard", mock.Anything).Return(custom_errors.NotFound(context.TODO(), "ISO Standard"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   custom_errors.NotFound(context.TODO(), "ISO Standard").Error(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setupMock()
			w := suite.performRequest("PUT", "/api/iso_standards/"+tc.id, bytes.NewBufferString(tc.body))
			suite.validateResponse(w, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard() {
	testCases := []struct {
		name           string
		id             string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			id:   "1",
			setupMock: func() {
				suite.mockRepo.On("DeleteISOStandard", int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"ISO standard deleted"}`,
		},
		{
			name: "NotFound",
			id:   "2",
			setupMock: func() {
				suite.mockRepo.On("DeleteISOStandard", int64(2)).Return(custom_errors.NotFound(context.TODO(), "ISO Standard"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "ISO Standard not found",
		},
		{
			name:           "CannotParseIDAsIntError",
			id:             "a",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid ISO Standard ID",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setupMock()
			w := suite.performRequest("DELETE", "/api/iso_standards/"+tc.id, nil)
			suite.validateResponse(w, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestApiISOStandardController(t *testing.T) {
	suite.Run(t, new(IsoStandardControllerTestSuite))
}
