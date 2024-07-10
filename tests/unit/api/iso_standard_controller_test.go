package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ISO_Auditing_Tool/cmd/api/controllers"
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
}

func (suite *IsoStandardControllerTestSuite) SetupTest() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	suite.router = gin.Default()
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

func (suite *IsoStandardControllerTestSuite) performRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
	if body == nil {
		body = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, url, body)
	require.NoError(suite.T(), err, "error creating new request")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *IsoStandardControllerTestSuite) validateResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBodyContains string) {
	assert.Equal(suite.T(), expectedStatus, w.Code)
	assert.Contains(suite.T(), w.Body.String(), expectedBodyContains)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Success() {
	expectedStandard := suite.jsonData[0]
	expectedStandards := []types.ISOStandard{expectedStandard}
	suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	w := suite.performRequest("GET", "/api/iso_standards", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Success() {
	expectedStandard := suite.jsonData[0]
	suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(expectedStandard, nil)

	w := suite.performRequest("GET", "/api/iso_standards/1", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
	expectedStandard := suite.jsonData[0]

	formData, err := json.Marshal(expectedStandard)
	require.NoError(suite.T(), err, "failed to marshal JSON data")

	suite.mockRepo.On("CreateISOStandard", expectedStandard).Return(expectedStandard, nil)
	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBuffer(formData))
	suite.validateResponse(w, http.StatusCreated, `"name":"Leadership"`)
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_InvalidData() {
	invalidFormData := `{"name": ""}`
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()

	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))
	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid data"}`)
}

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

	suite.mockRepo.On("UpdateISOStandard", expectedStandard).Return(errors.New("not found"))

	w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(formData))
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_Success() {
	suite.mockRepo.On("DeleteISOStandard", int64(1)).Return(nil)

	w := suite.performRequest("DELETE", "/api/iso_standards/1", nil)
	suite.validateResponse(w, http.StatusOK, `{"message":"ISO standard deleted"}`)
}

func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_NotFound() {
	suite.mockRepo.On("DeleteISOStandard", int64(2)).Return(errors.New("not found"))

	w := suite.performRequest("DELETE", "/api/iso_standards/2", nil)
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_NotFound() {
	suite.mockRepo.On("GetISOStandardByID", int64(2)).Return(types.ISOStandard{}, errors.New("not found"))

	w := suite.performRequest("GET", "/api/iso_standards/2", nil)
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Error() {
	suite.mockRepo.On("GetAllISOStandards").Return(nil, errors.New("database error"))

	w := suite.performRequest("GET", "/api/iso_standards", nil)
	suite.validateResponse(w, http.StatusInternalServerError, `{"error":"database error"}`)
}

func TestApiIsoStandardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(IsoStandardControllerTestSuite))
}
