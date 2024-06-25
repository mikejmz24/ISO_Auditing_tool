package api_test

import (
	"ISO_Auditing_Tool/cmd/api/controllers"
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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IsoStandardControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockRepo    *testutils.MockIsoStandardRepository
	standard    types.ISOStandard
	formData    []byte
	updatedData []byte
}

func (suite *IsoStandardControllerTestSuite) setupMockRepo() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	if suite.mockRepo == nil {
		panic("mockRepo is nil")
	}
	fmt.Printf("Mock Repository initialized: %v\n", suite.mockRepo)
}

func (suite *IsoStandardControllerTestSuite) setupRouter() {
	suite.router = setupRouter(suite.mockRepo)
	if suite.router == nil {
		panic("router is nil")
	}
	fmt.Printf("Router initialized: %v\n", suite.router)
}

func (suite *IsoStandardControllerTestSuite) setupSampleData() {
	filePath := "../../testdata/iso_standards_test01.json"
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to open JSON file: %w", err))
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("failed to read JSON file: %w", err))
	}

	var jsonData struct {
		ISOStandards []types.ISOStandard `json:"iso_standards"`
	}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal JSON data: %w", err))
	}
	// Assuming the first standard in the JSON file for testing purposes
	if len(jsonData.ISOStandards) > 0 {
		suite.standard = jsonData.ISOStandards[0]
	} else {
		panic("no ISO standards found in JSON data")
	}

	suite.formData, err = json.Marshal(suite.standard)
	if err != nil {
		panic(fmt.Errorf("failed to marshal JSON data: %w", err))
	}

	fmt.Printf("Unmarshalled Data: %+v\n", suite.standard)

	suite.updatedData = []byte(`{
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
    }`)
}

func (suite *IsoStandardControllerTestSuite) SetupTest() {
	fmt.Println("Setting up test")
	suite.setupMockRepo()
	suite.setupRouter()
	suite.setupSampleData()
	fmt.Printf("Setup complete: router=%v, mockRepo=%v, sampleData=%v\n", suite.router, suite.mockRepo, suite.standard)
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
	fmt.Printf("Router setup with routes: %v\n", api)
	return router
}

func (suite *IsoStandardControllerTestSuite) performRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
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

func (suite *IsoStandardControllerTestSuite) validateResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBodyContains string) {
	fmt.Printf("Validating response: status=%d, body=%s\n", w.Code, w.Body.String())
	require.Equal(suite.T(), expectedStatus, w.Code)
	assert.Contains(suite.T(), w.Body.String(), expectedBodyContains)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *IsoStandardControllerTestSuite) TestAPIGetAllISOStandards() {
	fmt.Println("Running TestAPIGetAllISOStandards")
	expectedStandards := []types.ISOStandard{suite.standard}
	suite.mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	w := suite.performRequest("GET", "/api/iso_standards", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestAPIGetISOStandardByID() {
	fmt.Println("Running TestAPIGetISOStandardByID")
	suite.mockRepo.On("GetISOStandardByID", int64(1)).Return(suite.standard, nil)

	w := suite.performRequest("GET", "/api/iso_standards/1", nil)
	suite.validateResponse(w, http.StatusOK, "ISO 9001")
}

func (suite *IsoStandardControllerTestSuite) TestAPICreateISOStandard() {
	fmt.Println("Running TestAPICreateISOStandard")
	expectedID := int64(1)
	suite.mockRepo.On("CreateISOStandard", suite.standard).Return(expectedID, nil)

	fmt.Printf("Form Data: %s\n", suite.formData) // Print form data for debugging

	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBuffer(suite.formData))
	fmt.Printf("Response Body: %s\n", w.Body.String()) // Print response body for debugging
	suite.validateResponse(w, http.StatusCreated, `"id":1`)
}

func (suite *IsoStandardControllerTestSuite) TestAPICreateISOStandardInvalidData() {
	fmt.Println("Running TestAPICreateISOStandardInvalidData")
	invalidFormData := `{"name": ""}`

	suite.mockRepo.On("CreateISOStandard", mock.Anything).Return(int64(0), nil).Maybe()

	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBufferString(invalidFormData))
	suite.validateResponse(w, http.StatusBadRequest, `{"error":"Invalid data"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPIUpdateISOStandard() {
	fmt.Println("Running TestAPIUpdateISOStandard")
	updatedStandard := suite.standard
	updatedStandard.Name = "ISO 9001 Updated"
	updatedStandard.Clauses = &[]types.Clause{
		{
			ID:   1,
			Name: "Clause 1",
			Sections: &[]types.Section{
				{
					ID:   1,
					Name: "Section 1",
					Questions: &[]types.Question{
						{
							ID:   1,
							Text: "Question 1",
						},
					},
				},
			},
		},
	}

	suite.mockRepo.On("UpdateISOStandard", mock.MatchedBy(func(isoStandard types.ISOStandard) bool {
		return isoStandard.ID == updatedStandard.ID &&
			isoStandard.Name == updatedStandard.Name &&
			len(*isoStandard.Clauses) == len(*updatedStandard.Clauses) &&
			(*isoStandard.Clauses)[0].Name == (*updatedStandard.Clauses)[0].Name
	})).Return(nil)

	w := suite.performRequest("PUT", "/api/iso_standards/1", bytes.NewBuffer(suite.updatedData))
	suite.validateResponse(w, http.StatusOK, `{"status":"updated"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPIUpdateISOStandardNotFound() {
	fmt.Println("Running TestAPIUpdateISOStandardNotFound")
	suite.mockRepo.On("UpdateISOStandard", suite.standard).Return(errors.New("not found"))

	w := suite.performRequest("PUT", "/api/iso_standards/2", bytes.NewBuffer(suite.formData))
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPIDeleteISOStandard() {
	fmt.Println("Running TestAPIDeleteISOStandard")
	suite.mockRepo.On("DeleteISOStandard", int64(1)).Return(nil)

	fmt.Println("Mock setup complete for TestAPIDeleteISOStandard")
	w := suite.performRequest("DELETE", "/api/iso_standards/1", nil)
	fmt.Println("Request performed for TestAPIDeleteISOStandard")
	suite.validateResponse(w, http.StatusOK, `{"message":"ISO standard deleted"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPIDeleteISOStandardNotFound() {
	fmt.Println("Running TestAPIDeleteISOStandardNotFound")
	suite.mockRepo.On("DeleteISOStandard", int64(2)).Return(errors.New("not found"))

	fmt.Println("Mock setup complete for TestAPIDeleteISOStandardNotFound")
	w := suite.performRequest("DELETE", "/api/iso_standards/2", nil)
	fmt.Println("Request performed for TestAPIDeleteISOStandardNotFound")
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPIGetISOStandardByIDNotFound() {
	fmt.Println("Running TestAPIGetISOStandardByIDNotFound")
	suite.mockRepo.On("GetISOStandardByID", int64(2)).Return(types.ISOStandard{}, errors.New("not found"))

	w := suite.performRequest("GET", "/api/iso_standards/2", nil)
	suite.validateResponse(w, http.StatusNotFound, `{"error":"ISO standard not found"}`)
}

func (suite *IsoStandardControllerTestSuite) TestAPICreateISOStandardInternalServerError() {
	fmt.Println("Running TestAPICreateISOStandardInternalServerError")
	suite.mockRepo.On("CreateISOStandard", suite.standard).Return(int64(0), errors.New("internal server error"))

	w := suite.performRequest("POST", "/api/iso_standards", bytes.NewBuffer(suite.formData))
	suite.validateResponse(w, http.StatusInternalServerError, `{"error":"internal server error"}`)
}

func TestIsoStandardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(IsoStandardControllerTestSuite))
}
