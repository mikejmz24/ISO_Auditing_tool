// package controller_test
//
// import (
// 	"ISO_Auditing_Tool/pkg/repositories"
// 	"ISO_Auditing_Tool/pkg/types"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"testing"
// 	"time"
//
// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestNewDraftController(t *testing.T) {
// 	// Arrange
// 	db, _, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
//
// 	// Act
// 	repo, err := repositories.NewDraftRepository(db)
//
// 	// Assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, repo)
// }
//
// func TestDraftController_Success(t *testing.T) {
// 	// Common setup
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
//
// 	repo, err := repositories.NewDraftRepository(db)
// 	assert.NoError(t, err)
//
// 	ctx := context.Background()
// 	now := time.Now().UTC().Truncate(time.Second)
//
// 	data := json.RawMessage(`{"field1": "value1", "field2": "value2"}`)
// 	diff := json.RawMessage(`{"changed": "content"}`)
//
// 	tests := []struct {
// 		name       string
// 		draft      types.Draft
// 		mockSetup  func()
// 		expectedID int
// 	}{
// 		{
// 			name: "Create",
// 			draft: types.Draft{
// 				TypeID:          1,
// 				ObjectID:        2,
// 				StatusID:        3,
// 				Version:         1,
// 				Data:            data,
// 				Diff:            diff,
// 				UserID:          42,
// 				ApproverID:      0,
// 				ApprovalComment: "",
// 				PublishError:    "",
// 				CreatedAt:       now,
// 				UpdatedAt:       now,
// 				ExpiresAt:       now.Add(24 * time.Hour),
// 			},
// 			mockSetup: func() {
// 				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			expectedID: 1,
// 		},
// 	}
//
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Set up test case
// 			tc.mockSetup()
//
// 			// Act
// 			result, err := repo.Create(ctx, tc.draft)
//
// 			// Assert
// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.expectedID, result.ID)
//
// 			// Verify expectations for this test case
// 			err = mock.ExpectationsWereMet()
// 			assert.NoError(t, err)
// 		})
// 	}
// }
//
// func TestDraftController_Errors(t *testing.T) {
// 	// Common setup
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
//
// 	repo, err := repositories.NewDraftRepository(db)
// 	assert.NoError(t, err)
//
// 	ctx := context.Background()
//
// 	tests := []struct {
// 		name          string
// 		draft         types.Draft
// 		mockSetup     func()
// 		expectedError string
// 	}{
// 		{
// 			name: "An execution error on Create returns failed to create draft",
// 			draft: types.Draft{
// 				TypeID:   1,
// 				ObjectID: 2,
// 				UserID:   42,
// 			},
// 			mockSetup: func() {
// 				execError := errors.New("execution error")
// 				mock.ExpectExec("").WillReturnError(execError)
// 			},
// 			expectedError: "failed to create draft",
// 		},
// 		{
// 			name: "SQL last insert error on Create returns failed to get last insert ID",
// 			draft: types.Draft{
// 				TypeID:   1,
// 				ObjectID: 2,
// 				UserID:   42,
// 			},
// 			mockSetup: func() {
// 				lastIDError := errors.New("last insert ID error")
// 				mockResult := sqlmock.NewErrorResult(lastIDError)
// 				mock.ExpectExec("").WillReturnResult(mockResult)
// 			},
// 			expectedError: "failed to get last insert ID",
// 		},
// 	}
//
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Set up test case
// 			tc.mockSetup()
//
// 			// Act
// 			result, err := repo.Create(ctx, tc.draft)
//
// 			// Assert
// 			assert.Error(t, err)
// 			assert.Equal(t, types.Draft{}, result)
// 			assert.Contains(t, err.Error(), tc.expectedError)
//
// 			// Verify expectations for this test case
// 			err = mock.ExpectationsWereMet()
// 			assert.NoError(t, err)
// 		})
// 	}
// }

package controllers_test

import (
	"ISO_Auditing_Tool/pkg/controllers/api"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/unit/repositories/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Happy Path Test Suite
type ApiDraftControllerHappyPathSuite struct {
	suite.Suite
	mockService *mocks.MockDraftService
	controller  *controllers.ApiDraftController
}

// Error Handling Test Suite
type ApiDraftControllerErrorSuite struct {
	suite.Suite
	mockService *mocks.MockDraftService
	controller  *controllers.ApiDraftController
}

// Common setup
func (suite *ApiDraftControllerHappyPathSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockService = new(mocks.MockDraftService)
	suite.controller = controllers.NewAPIDraftController(suite.mockService)
}

func (suite *ApiDraftControllerErrorSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockService = new(mocks.MockDraftService)
	suite.controller = controllers.NewAPIDraftController(suite.mockService)
}

// Pre-allocated test data for performance
var (
	testDraft = types.Draft{
		ID:       1,
		TypeID:   1,
		ObjectID: 42,
		StatusID: 1,
		Version:  1,
		UserID:   10,
	}

	testDraftJSON, _  = json.Marshal(testDraft)
	validJSONBuffer   = bytes.NewBuffer(testDraftJSON)
	invalidJSONBuffer = bytes.NewBufferString(`{"invalid": json}`)
)

// Helper to create Gin context with minimal overhead
func createTestContext(method, path string, body *bytes.Buffer) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

// --- Happy Path Tests ---

func (suite *ApiDraftControllerHappyPathSuite) TestCreate_ValidInput_ReturnsCreated() {
	// Setup
	suite.mockService.On("Create", mock.Anything, mock.AnythingOfType("types.Draft")).
		Return(testDraft, nil)

	c, w := createTestContext("POST", "/drafts", bytes.NewBuffer(testDraftJSON))

	// Act
	suite.controller.Create(c)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ApiDraftControllerHappyPathSuite) TestUpdate_ValidInput_ReturnsOK() {
	// Setup
	suite.mockService.On("Update", mock.Anything, mock.AnythingOfType("types.Draft")).
		Return(testDraft, nil)

	c, w := createTestContext("PUT", "/drafts/1", bytes.NewBuffer(testDraftJSON))
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	suite.controller.Update(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// --- Error Handling Tests ---

func (suite *ApiDraftControllerErrorSuite) TestCreate_InvalidJSON_ReturnsBadRequest() {
	// Setup
	c, w := createTestContext("POST", "/drafts", invalidJSONBuffer)

	// Act
	suite.controller.Create(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ApiDraftControllerErrorSuite) TestCreate_ServiceError_ReturnsInternalServerError() {
	// Setup
	suite.mockService.On("Create", mock.Anything, mock.AnythingOfType("types.Draft")).
		Return(types.Draft{}, errors.New("service error"))

	c, w := createTestContext("POST", "/drafts", bytes.NewBuffer(testDraftJSON))

	// Act
	suite.controller.Create(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ApiDraftControllerErrorSuite) TestUpdate_InvalidJSON_ReturnsBadRequest() {
	// Setup
	c, w := createTestContext("PUT", "/drafts/1", invalidJSONBuffer)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	suite.controller.Update(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ApiDraftControllerErrorSuite) TestUpdate_InvalidID_ReturnsBadRequest() {
	// Setup
	c, w := createTestContext("PUT", "/drafts/invalid", bytes.NewBuffer(testDraftJSON))
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	// Act
	suite.controller.Update(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ApiDraftControllerErrorSuite) TestUpdate_ServiceError_ReturnsInternalServerError() {
	// Setup
	suite.mockService.On("Update", mock.Anything, mock.AnythingOfType("types.Draft")).
		Return(types.Draft{}, errors.New("service error"))

	c, w := createTestContext("PUT", "/drafts/1", bytes.NewBuffer(testDraftJSON))
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	suite.controller.Update(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// Test runners
func TestApiDraftController_HappyPath(t *testing.T) {
	suite.Run(t, new(ApiDraftControllerHappyPathSuite))
}

func TestApiDraftController_ErrorHandling(t *testing.T) {
	suite.Run(t, new(ApiDraftControllerErrorSuite))
}
