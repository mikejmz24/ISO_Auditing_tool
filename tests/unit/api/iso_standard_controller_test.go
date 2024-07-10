// package api_test
//
// import (
//
//	"ISO_Auditing_Tool/cmd/api/controllers"
//	"ISO_Auditing_Tool/pkg/types"
//	"ISO_Auditing_Tool/tests/testutils"
//	"bytes"
//	"encoding/json"
//	"errors"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"os"
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
//	type IsoStandardControllerTestSuite struct {
//		suite.Suite
//		router      *gin.Engine
//		mockRepo    *testutils.MockIsoStandardRepository
//		standard    types.ISOStandard
//		formData    []byte
//		updatedData []byte
//	}
//
//	func (suite *IsoStandardControllerTestSuite) SetupSuite() {
//		suite.setupMockRepo()
//		suite.setupRouter()
//		suite.setupSampleData()
//	}
//
//	func (suite *IsoStandardControllerTestSuite) setupMockRepo() {
//		suite.mockRepo = new(testutils.MockIsoStandardRepository)
//		require.NotNil(suite.T(), suite.mockRepo, "mockRepo is nil")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) setupRouter() {
//		suite.router = gin.Default()
//		controller := controllers.NewApiIsoStandardController(suite.mockRepo)
//		api := suite.router.Group("/api")
//		{
//			api.GET("/iso_standards", controller.GetAllISOStandards)
//			api.GET("/iso_standards/:id", controller.GetISOStandardByID)
//			api.POST("/iso_standards", controller.CreateISOStandard)
//			api.PUT("/iso_standards/:id", controller.UpdateISOStandard)
//			api.DELETE("/iso_standards/:id", controller.DeleteISOStandard)
//		}
//		require.NotNil(suite.T(), suite.router, "router is nil")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) setupSampleData() {
//		filePath := "../../testdata/iso_standards_test01.json"
//		file, err := os.Open(filePath)
//		require.NoError(suite.T(), err, "failed to open JSON file")
//		defer file.Close()
//
//		data, err := io.ReadAll(file)
//		require.NoError(suite.T(), err, "failed to read JSON file")
//
//		var jsonData struct {
//			ISOStandards []types.ISOStandard `json:"iso_standards"`
//		}
//		err = json.Unmarshal(data, &jsonData)
//		require.NoError(suite.T(), err, "failed to unmarshal JSON data")
//		require.NotEmpty(suite.T(), jsonData.ISOStandards, "no ISO standards found in JSON data")
//
//		suite.standard = jsonData.ISOStandards[0]
//
//		suite.formData, err = json.Marshal(suite.standard)
//		require.NoError(suite.T(), err, "failed to marshal JSON data")
//
//		suite.updatedData = []byte(`{
//			"id": 1,
//			"name": "ISO 9001 Updated",
//			"clauses": [{
//				"id": 1,
//				"name": "Clause 1",
//				"sections": [{
//					"id": 1,
//					"name": "Section 1",
//					"questions": [{
//						"id": 1,
//						"text": "Question 1"
//					}]
//				}]
//			}]
//		}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) performRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
//		if body == nil {
//			body = bytes.NewBuffer([]byte{})
//		}
//
//		req, err := http.NewRequest(method, url, body)
//		require.NoError(suite.T(), err, "error creating new request")
//
//		if body != nil {
//			req.Header.Set("Content-Type", "application/json")
//		}
//
//		w := httptest.NewRecorder()
//		suite.router.ServeHTTP(w, req)
//		return w
//	}
//
//	func (suite *IsoStandardControllerTestSuite) validateResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBodyContains string) {
//		assert.Equal(suite.T(), expectedStatus, w.Code)
//		assert.Contains(suite.T(), w.Body.String(), expectedBodyContains)
//		suite.mockRepo.AssertExpectations(suite.T())
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Success() {
//		expectedStandards := []types.ISOStandard{suite.standard}
//		suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)
//
//		w := suite.performRequest("GET", "/api/iso_standards", nil)
//		suite.validateResponse(w, http.StatusOK, "ISO 9001")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Success() {
//		suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(suite.standard, nil)
//
//		w := suite.performRequest("GET", "/api/iso_standards/1", nil)
//		suite.validateResponse(w, http.StatusOK, "ISO 9001")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
//		suite.mockRepo.On("CreateISOStandard", suite.standard).Return(suite.standard, nil)
//		w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBuffer(suite.formData))
//		suite.validateResponse(w, http.StatusCreated, `"name":"Leadership"`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_InvalidData() {
//		invalidFormData := `{"name": ""}`
//		suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()
//
//		w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))
//		suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid data"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestUpdateISOStandard_Success() {
//		updatedStandard := suite.standard
//		updatedStandard.Name = "ISO 9001 Updated"
//		updatedStandard.Clauses = &[]types.Clause{
//			{
//				ID:   1,
//				Name: "Clause 1",
//				Sections: &[]types.Section{
//					{
//						ID:   1,
//						Name: "Section 1",
//						Questions: &[]types.Question{
//							{
//								ID:   1,
//								Text: "Question 1",
//							},
//						},
//					},
//				},
//			},
//		}
//
//		suite.mockRepo.On("UpdateISOStandard", mock.MatchedBy(func(isoStandard types.ISOStandard) bool {
//			return isoStandard.ID == updatedStandard.ID &&
//				isoStandard.Name == updatedStandard.Name &&
//				len(*isoStandard.Clauses) == len(*updatedStandard.Clauses) &&
//				(*isoStandard.Clauses)[0].Name == (*updatedStandard.Clauses)[0].Name
//		})).Return(nil)
//
//		w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(suite.updatedData))
//		suite.validateResponse(w, http.StatusOK, `{"status":"updated"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestUpdateISOStandard_NotFound() {
//		suite.mockRepo.On("UpdateISOStandard", suite.standard).Return(errors.New("not found"))
//
//		w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(suite.formData))
//		suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_Success() {
//		suite.mockRepo.On("DeleteISOStandard", int64(1)).Return(nil)
//
//		w := suite.performRequest("DELETE", "/api/iso_standards/1", nil)
//		suite.validateResponse(w, http.StatusOK, `{"message":"ISO standard deleted"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_NotFound() {
//		suite.mockRepo.On("DeleteISOStandard", int64(2)).Return(errors.New("not found"))
//
//		w := suite.performRequest("DELETE", "/api/iso_standards/2", nil)
//		suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_NotFound() {
//		suite.mockRepo.On("GetISOStandardByID", int64(2)).Return(types.ISOStandard{}, errors.New("not found"))
//
//		w := suite.performRequest("GET", "/api/iso_standards/2", nil)
//		suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_InternalServerError() {
//		suite.mockRepo.On("CreateISOStandard", suite.standard).Return(types.ISOStandard{}, errors.New("internal server error"))
//
//		w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBuffer(suite.formData))
//		suite.validateResponse(w, http.StatusInternalServerError, `{"error":"internal server error"}`)
//	}
//
//	func TestApiIsoStandardControllerTestSuite(t *testing.T) {
//		suite.Run(t, new(IsoStandardControllerTestSuite))
//	}
//
// package api_test
//
// import (
//
//	"ISO_Auditing_Tool/cmd/api/controllers"
//	"ISO_Auditing_Tool/pkg/types"
//	"ISO_Auditing_Tool/tests/testutils"
//	"bytes"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"os"
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
//	type IsoStandardControllerTestSuite struct {
//		suite.Suite
//		router      *gin.Engine
//		mockRepo    *testutils.MockIsoStandardRepository
//		standard    types.ISOStandard
//		formData    []byte
//		updatedData []byte
//	}
//
//	func (suite *IsoStandardControllerTestSuite) SetupSuite() {
//		suite.setupMockRepo()
//		suite.setupRouter()
//		suite.setupSampleData()
//	}
//
//	func (suite *IsoStandardControllerTestSuite) setupMockRepo() {
//		suite.mockRepo = new(testutils.MockIsoStandardRepository)
//		require.NotNil(suite.T(), suite.mockRepo, "mockRepo is nil")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) setupRouter() {
//		suite.router = gin.Default()
//		controller := controllers.NewApiIsoStandardController(suite.mockRepo)
//		api := suite.router.Group("/api")
//		{
//			api.GET("/iso_standards", controller.GetAllISOStandards)
//			api.GET("/iso_standards/:id", controller.GetISOStandardByID)
//			api.POST("/iso_standards", controller.CreateISOStandard)
//			api.PUT("/iso_standards/:id", controller.UpdateISOStandard)
//			api.DELETE("/iso_standards/:id", controller.DeleteISOStandard)
//		}
//		require.NotNil(suite.T(), suite.router, "router is nil")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) setupSampleData() {
//		filePath := "../../testdata/iso_standards_test01.json"
//		file, err := os.Open(filePath)
//		require.NoError(suite.T(), err, "failed to open JSON file")
//		defer file.Close()
//
//		data, err := io.ReadAll(file)
//		require.NoError(suite.T(), err, "failed to read JSON file")
//
//		var jsonData struct {
//			ISOStandards []types.ISOStandard `json:"iso_standards"`
//		}
//		err = json.Unmarshal(data, &jsonData)
//		require.NoError(suite.T(), err, "failed to unmarshal JSON data")
//		require.NotEmpty(suite.T(), jsonData.ISOStandards, "no ISO standards found in JSON data")
//
//		suite.standard = jsonData.ISOStandards[0]
//
//		suite.formData, err = json.Marshal(suite.standard)
//		require.NoError(suite.T(), err, "failed to marshal JSON data")
//
//		suite.updatedData = []byte(`{
//			"id": 1,
//			"name": "ISO 9001 Updated",
//			"clauses": [{
//				"id": 1,
//				"name": "Clause 1",
//				"sections": [{
//					"id": 1,
//					"name": "Section 1",
//					"questions": [{
//						"id": 1,
//						"text": "Question 1"
//					}]
//				}]
//			}]
//		}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) performRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
//		if body == nil {
//			body = bytes.NewBuffer([]byte{})
//		}
//
//		req, err := http.NewRequest(method, url, body)
//		require.NoError(suite.T(), err, "error creating new request")
//
//		if body != nil {
//			req.Header.Set("Content-Type", "application/json")
//		}
//
//		w := httptest.NewRecorder()
//		suite.router.ServeHTTP(w, req)
//		return w
//	}
//
//	func (suite *IsoStandardControllerTestSuite) validateResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBodyContains string) {
//		assert.Equal(suite.T(), expectedStatus, w.Code)
//		assert.Contains(suite.T(), w.Body.String(), expectedBodyContains)
//		suite.mockRepo.AssertExpectations(suite.T())
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Success() {
//		expectedStandards := []types.ISOStandard{suite.standard}
//		suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)
//
//		w := suite.performRequest("GET", "/api/iso_standards", nil)
//		suite.validateResponse(w, http.StatusOK, "ISO 9001")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Success() {
//		suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(suite.standard, nil)
//
//		w := suite.performRequest("GET", "/api/iso_standards/1", nil)
//		suite.validateResponse(w, http.StatusOK, "ISO 9001")
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
//		fmt.Printf("suite.standard: %v", suite.standard)
//		fmt.Println()
//		fmt.Println()
//		fmt.Printf("suite.formData: %v", suite.formData)
//		fmt.Println()
//		fmt.Println()
//
//		suite.mockRepo.On("CreateISOStandard", suite.standard).Return(suite.standard, nil)
//		w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBuffer(suite.formData))
//		fmt.Printf("w: %v", w)
//		fmt.Println()
//		fmt.Println()
//		suite.validateResponse(w, http.StatusCreated, `"name":"Leadership"`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_InvalidData() {
//		invalidFormData := `{"name": ""}`
//		suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(types.ISOStandard{}, nil).Maybe()
//
//		w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))
//		suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid data"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestUpdateISOStandard_Success() {
//		updatedStandard := suite.standard
//		updatedStandard.Name = "ISO 9001 Updated"
//		updatedStandard.Clauses = []*types.Clause{
//			{
//				ID:   1,
//				Name: "Clause 1",
//				Sections: []*types.Section{
//					{
//						ID:   1,
//						Name: "Section 1",
//						Questions: []*types.Question{
//							{
//								ID:   1,
//								Text: "Question 1",
//							},
//						},
//					},
//				},
//			},
//		}
//
//		suite.mockRepo.On("UpdateISOStandard", mock.MatchedBy(func(isoStandard types.ISOStandard) bool {
//			return isoStandard.ID == updatedStandard.ID &&
//				isoStandard.Name == updatedStandard.Name &&
//				len(isoStandard.Clauses) == len(updatedStandard.Clauses) &&
//				(isoStandard.Clauses)[0].Name == (updatedStandard.Clauses)[0].Name
//		})).Return(nil)
//
//		w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(suite.updatedData))
//		suite.validateResponse(w, http.StatusOK, `{"status":"updated"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestUpdateISOStandard_NotFound() {
//		suite.mockRepo.On("UpdateISOStandard", suite.standard).Return(errors.New("not found"))
//
//		w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(suite.formData))
//		suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_Success() {
//		suite.mockRepo.On("DeleteISOStandard", int64(1)).Return(nil)
//
//		w := suite.performRequest("DELETE", "/api/iso_standards/1", nil)
//		suite.validateResponse(w, http.StatusOK, `{"message":"ISO standard deleted"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestDeleteISOStandard_NotFound() {
//		suite.mockRepo.On("DeleteISOStandard", int64(2)).Return(errors.New("not found"))
//
//		w := suite.performRequest("DELETE", "/api/iso_standards/2", nil)
//		suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_NotFound() {
//		suite.mockRepo.On("GetISOStandardByID", int64(2)).Return(types.ISOStandard{}, errors.New("not found"))
//
//		w := suite.performRequest("GET", "/api/iso_standards/2", nil)
//		suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
//	}
//
//	func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Error() {
//		suite.mockRepo.On("GetAllISOStandards").Return(nil, errors.New("database error"))
//
//		w := suite.performRequest("GET", "/api/iso_standards", nil)
//		suite.validateResponse(w, http.StatusInternalServerError, `{"error":"database error"}`)
//	}
//
//	func TestApiIsoStandardControllerTestSuite(t *testing.T) {
//		suite.Run(t, new(IsoStandardControllerTestSuite))
//	}
package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	// "fmt"
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
}

func (suite *IsoStandardControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
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

func (suite *IsoStandardControllerTestSuite) loadSampleData() types.ISOStandard {
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

	return jsonData.ISOStandards[0]
}

func (suite *IsoStandardControllerTestSuite) TestGetAllISOStandards_Success() {
	expectedStandard := suite.loadSampleData()
	expectedStandards := []types.ISOStandard{expectedStandard}
	suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	w := suite.performRequest("GET", "/api/iso_standards", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestGetISOStandardByID_Success() {
	expectedStandard := suite.loadSampleData()
	suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(expectedStandard, nil)

	w := suite.performRequest("GET", "/api/iso_standards/1", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestCreateISOStandard_Success() {
	expectedStandard := suite.loadSampleData()

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
	expectedStandard := suite.loadSampleData()
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
	expectedStandard := suite.loadSampleData()

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
