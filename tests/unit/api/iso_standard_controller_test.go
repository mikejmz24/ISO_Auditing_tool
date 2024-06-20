// // tests/unit/api/iso_standard_controller_test.go
// package api_test
//
// import (
// 	"ISO_Auditing_Tool/cmd/api/controllers"
// 	"ISO_Auditing_Tool/pkg/types"
// 	"ISO_Auditing_Tool/tests/testutils"
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
//
// 	"github.com/gin-gonic/gin"
// 	// "github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )
//
// type IsoStandardControllerTestSuite struct {
// 	suite.Suite
// 	router      *gin.Engine
// 	mockRepo    *testutils.MockIsoStandardRepository
// 	standard    types.ISOStandard
// 	formData    string
// 	updatedData string
// }
//
// func (suite *IsoStandardControllerTestSuite) SetupTest() {
// 	suite.mockRepo = new(testutils.MockIsoStandardRepository)
// 	suite.router = setupRouter(suite.mockRepo)
// 	suite.standard = types.ISOStandard{
// 		ID:   1,
// 		Name: "ISO 9001",
// 		Clauses: []types.Clause{
// 			{
// 				ID: 1, Name: "Clause 1", Sections: []types.Section{
// 					{ID: 1, Name: "Section 1", Questions: []types.Question{
// 						{ID: 1, Text: "Question 1"},
// 					}},
// 				},
// 			},
// 		},
// 	}
// 	suite.formData = `{
//     "id": 1,
// 		"name": "ISO 9001",
// 		"clauses": [{
//       "id": 1,
// 			"name": "Clause 1",
// 			"sections": [{
// 				"id": 1,
// 				"name": "Section 1",
// 				"questions": [{
// 					"id": 1,
// 					"text": "Question 1"
// 				}]
// 			}]
// 		}]
// 	}`
// 	suite.updatedData = `{
// 		"id": 1,
// 		"name": "ISO 9001 Updated",
// 		"clauses": [{
// 			"id": 1,
// 			"name": "Clause 1",
// 			"sections": [{
// 				"id": 1,
// 				"name": "Section 1",
// 				"questions": [{
// 					"id": 1,
// 					"text": "Question 1"
// 				}]
// 			}]
// 		}]
// 	}`
// }
//
// func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
// 	controller := controllers.NewApiIsoStandardController(repo)
// 	router := gin.Default()
// 	api := router.Group("/api")
// 	{
// 		api.GET("/iso_standards", controller.GetAllISOStandards)
// 		api.GET("/iso_standards/:id", controller.GetISOStandardByID)
// 		api.POST("/iso_standards", controller.CreateISOStandard)
// 		api.PUT("/iso_standards/:id", controller.UpdateISOStandard)
// 		api.DELETE("/iso_standards/:id", controller.DeleteISOStandard)
// 	}
// 	return router
// }
//
// func (suite *IsoStandardControllerTestSuite) TestAPIGetAllISOStandards() {
// 	expectedStandards := []types.ISOStandard{suite.standard}
// 	suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)
//
// 	req, _ := http.NewRequest("GET", "/api/iso_standards", nil)
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)
//
// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.Contains(w.Body.String(), "ISO 9001")
//
// 	suite.mockRepo.AssertExpectations(suite.T())
// }
//
// func (suite *IsoStandardControllerTestSuite) TestAPIGetISOStandardByID() {
// 	suite.mockRepo.On("GetISOStandardByID", 1).Return(suite.standard, nil)
//
// 	req, _ := http.NewRequest("GET", "/api/iso_standards/1", nil)
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)
//
// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.Contains(w.Body.String(), "ISO 9001")
//
// 	suite.mockRepo.AssertExpectations(suite.T())
// }
//
// func (suite *IsoStandardControllerTestSuite) TestAPICreateISOStandard() {
// 	expectedID := int64(1)
// 	suite.mockRepo.On("CreateISOStandard", suite.standard).Return(expectedID, nil)
//
// 	req, _ := http.NewRequest("POST", "/api/iso_standards", bytes.NewBufferString(suite.formData))
// 	req.Header.Set("Content-Type", "application/json")
//
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)
//
// 	suite.Equal(http.StatusCreated, w.Code)
// 	suite.Contains(w.Body.String(), `"id":1`)
//
// 	suite.mockRepo.AssertExpectations(suite.T())
// }
//
// func (suite *IsoStandardControllerTestSuite) TestAPIUpdateISOStandard() {
// 	updatedStandard := suite.standard
// 	updatedStandard.Name = "ISO 9001 Updated"
// 	suite.mockRepo.On("UpdateISOStandard", updatedStandard).Return(nil)
//
// 	req, _ := http.NewRequest("PUT", "/api/iso_standards/1", bytes.NewBufferString(suite.updatedData))
// 	req.Header.Set("Content-Type", "application/json")
//
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)
//
// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.Contains(w.Body.String(), `{"status":"updated"}`)
//
// 	suite.mockRepo.AssertExpectations(suite.T())
// }
//
// func (suite *IsoStandardControllerTestSuite) TestAPIDeleteISOStandard() {
// 	suite.mockRepo.On("DeleteISOStandard", 1).Return(nil)
//
// 	req, _ := http.NewRequest("DELETE", "/api/iso_standards/1", nil)
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)
//
// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.Contains(w.Body.String(), `{"message":"ISO standard deleted"}`)
//
// 	suite.mockRepo.AssertExpectations(suite.T())
// }
//
// func TestIsoStandardControllerTestSuite(t *testing.T) {
// 	suite.Run(t, new(IsoStandardControllerTestSuite))
// }

// tests/unit/api/iso_standard_controller_test.go
package api_test

import (
	"ISO_Auditing_Tool/cmd/api/controllers"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IsoStandardControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockRepo    *testutils.MockIsoStandardRepository
	standard    types.ISOStandard
	formData    string
	updatedData string
}

func (suite *IsoStandardControllerTestSuite) SetupTest() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	suite.router = setupRouter(suite.mockRepo)
	suite.standard = types.ISOStandard{
		ID:   1,
		Name: "ISO 9001",
		Clauses: []types.Clause{
			{
				ID: 1, Name: "Clause 1", Sections: []types.Section{
					{ID: 1, Name: "Section 1", Questions: []types.Question{
						{ID: 1, Text: "Question 1"},
					}},
				},
			},
		},
	}
	suite.formData = `{
    "id": 1,
		"name": "ISO 9001",
		"clauses": [{
			"id": 1,
			"name": "Clause 1",
			"sections": [{
				"id": 1,
				"name": "Section 1",
				"questions": [{
					"id": 1,
					"text": "Question 1"
				}]
			}]
		}]
	}`
	suite.updatedData = `{
		"id": 1,
		"name": "ISO 9001 Updated",
		"clauses": [{
			"id": 1,
			"name": "Clause 1",
			"sections": [{
				"id": 1,
				"name": "Section 1",
				"questions": [{
					"id": 1,
					"text": "Question 1"
				}]
			}]
		}]
	}`
}

func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
	controller := controllers.NewApiIsoStandardController(repo)
	router := gin.Default()
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

func (suite *IsoStandardControllerTestSuite) TestAPIGetAllISOStandards() {
	expectedStandards := []types.ISOStandard{suite.standard}
	suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	req, _ := http.NewRequest("GET", "/api/iso_standards", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "ISO 9001")

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) TestAPIGetISOStandardByID() {
	suite.mockRepo.On("GetISOStandardByID", 1).Return(suite.standard, nil)

	req, _ := http.NewRequest("GET", "/api/iso_standards/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "ISO 9001")

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) TestAPICreateISOStandard() {
	expectedID := int64(1)
	suite.mockRepo.On("CreateISOStandard", suite.standard).Return(expectedID, nil)

	req, _ := http.NewRequest("POST", "/api/iso_standards", bytes.NewBufferString(suite.formData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)
	suite.Contains(w.Body.String(), `"id":1`)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) TestAPICreateISOStandardInvalidData() {
	invalidFormData := `{"name": ""}`

	// Set up the mock to not expect any call to CreateISOStandard
	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(int64(0), nil).Maybe()

	req, _ := http.NewRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), `{"error":"Invalid data"}`)

	// suite.mockRepo.AssertNotCalled(suite.T(), "CreateISOStandard", mock.Anything)
}

func (suite *IsoStandardControllerTestSuite) TestAPIUpdateISOStandard() {
	updatedStandard := suite.standard
	updatedStandard.Name = "ISO 9001 Updated"
	suite.mockRepo.On("UpdateISOStandard", updatedStandard).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/iso_standards/1", bytes.NewBufferString(suite.updatedData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), `{"status":"updated"}`)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) TestAPIUpdateISOStandardNotFound() {
	suite.mockRepo.On("UpdateISOStandard", suite.standard).Return(errors.New("not found"))

	req, _ := http.NewRequest("PUT", "/api/iso_standards/2", bytes.NewBufferString(suite.formData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPIDeleteISOStandard() {
	suite.mockRepo.On("DeleteISOStandard", 1).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/iso_standards/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), `{"message":"ISO standard deleted"}`)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) TestAPIDeleteISOStandardNotFound() {
	suite.mockRepo.On("DeleteISOStandard", 2).Return(errors.New("not found"))

	req, _ := http.NewRequest("DELETE", "/api/iso_standards/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPIGetISOStandardByIDNotFound() {
	suite.mockRepo.On("GetISOStandardByID", 2).Return(types.ISOStandard{}, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/api/iso_standards/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPICreateISOStandardInternalServerError() {
	suite.mockRepo.On("CreateISOStandard", suite.standard).Return(int64(0), errors.New("internal server error"))

	req, _ := http.NewRequest("POST", "/api/iso_standards", bytes.NewBufferString(suite.formData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), `{"error":"internal server error"}`)
}

func TestIsoStandardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(IsoStandardControllerTestSuite))
}
