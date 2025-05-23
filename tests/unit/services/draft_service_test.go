package services_test

import (
	"ISO_Auditing_Tool/pkg/services"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock dependencies
type MockDraftRepository struct {
	mock.Mock
}

// Define test suites
type DraftServiceSuccessSuite struct {
	suite.Suite
	mockRepo *testutils.MockDraftRepository
	service  *services.DraftService
}

type DraftServiceErrorSuite struct {
	suite.Suite
	mockRepo *testutils.MockDraftRepository
	service  *services.DraftService
}

// SetupTest initializes test dependencies before each test
func (suite *DraftServiceSuccessSuite) SetupTest() {
	suite.mockRepo = new(testutils.MockDraftRepository)
	suite.service = services.NewDraftService(suite.mockRepo)
}

// TearDownTest cleans up after each test
func (suite *DraftServiceSuccessSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// SetupTest initializes test dependencies before each test
func (suite *DraftServiceErrorSuite) SetupTest() {
	suite.mockRepo = new(testutils.MockDraftRepository)
	suite.service = services.NewDraftService(suite.mockRepo)
}

// TearDownTest cleans up after each test
func (suite *DraftServiceErrorSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Helper functions
func createTestDraft() types.Draft {
	now := time.Now().UTC().Truncate(time.Second)
	return types.Draft{
		ID:              1,
		TypeID:          1,  // 1 might represent "standard"
		ObjectID:        42, // ID of the standard being drafted
		StatusID:        1,  // 1 might represent "pending"
		Version:         1,
		Data:            []byte(`{"name": "ISO 27001", "description": "Information Security Standard"}`),
		Diff:            []byte(`{"name": {"old": "ISO 27000", "new": "ISO 27001"}}`),
		UserID:          10,
		ApproverID:      0,
		ApprovalComment: "",
		PublishError:    "",
		CreatedAt:       now,
		UpdatedAt:       now,
		ExpiresAt:       now.Add(7 * 24 * time.Hour), // 1 week expiry
	}
}

// --- Success Test Cases ---

// TestNewDraftService_ReturnsServiceWithRepo tests the constructor function
func (suite *DraftServiceSuccessSuite) TestNewDraftService_ReturnsServiceWithRepo() {
	// Arrange & Act
	service := services.NewDraftService(suite.mockRepo)

	// Assert
	assert.NotNil(suite.T(), service)
	assert.Equal(suite.T(), suite.mockRepo, service.Repo)
}

// TestCreate_WhenValidInput_ReturnsDraftWithID tests successful draft creation
func (suite *DraftServiceSuccessSuite) TestCreate_WhenValidInput_ReturnsDraftWithID() {
	// Arrange
	ctx := context.Background()
	input := createTestDraft()
	input.ID = 0 // New draft has no ID yet

	expected := input
	expected.ID = 1

	suite.mockRepo.On("CreateDraft", ctx, input).Return(expected, nil)

	// Act
	result, err := suite.service.Create(ctx, input)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected.ID, result.ID)
	assert.Equal(suite.T(), expected.TypeID, result.TypeID)
	assert.Equal(suite.T(), expected.ObjectID, result.ObjectID)
	assert.Equal(suite.T(), expected.StatusID, result.StatusID)
	assert.Equal(suite.T(), expected.Version, result.Version)
	assert.Equal(suite.T(), string(expected.Data), string(result.Data))
}

// TestUpdate_WhenValidDraft_ReturnsUpdatedDraft tests draft update
func (suite *DraftServiceSuccessSuite) TestUpdate_WhenValidDraft_ReturnsUpdatedDraft() {
	// Arrange
	ctx := context.Background()
	draft := createTestDraft()

	// Updated draft with modified data
	updatedDraft := draft
	updatedDraft.Data = []byte(`{"name": "ISO 27001:2022", "description": "Updated Information Security Standard"}`)
	updatedDraft.Diff = []byte(`{"description": {"old": "Information Security Standard", "new": "Updated Information Security Standard"}}`)

	suite.mockRepo.On("UpdateDraft", ctx, draft).Return(updatedDraft, nil)

	// Act
	result, err := suite.service.Update(ctx, draft)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedDraft.ID, result.ID)
	assert.Equal(suite.T(), string(updatedDraft.Data), string(result.Data))
	assert.Equal(suite.T(), string(updatedDraft.Diff), string(result.Diff))
}

// --- Error Test Cases ---

// TestCreate_WhenRepositoryFails_ReturnsError tests error handling in draft creation
func (suite *DraftServiceErrorSuite) TestCreate_WhenRepositoryFails_ReturnsError() {
	// Arrange
	ctx := context.Background()
	input := createTestDraft()
	input.ID = 0 // New draft has no ID yet

	expectedErr := errors.New("database error")
	suite.mockRepo.On("CreateDraft", ctx, input).Return(types.Draft{}, expectedErr)

	// Act
	result, err := suite.service.Create(ctx, input)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
	assert.Equal(suite.T(), types.Draft{}, result)
}

// TestUpdate_WhenRepositoryFails_ReturnsError tests error handling in draft update
func (suite *DraftServiceErrorSuite) TestUpdate_WhenRepositoryFails_ReturnsError() {
	// Arrange
	ctx := context.Background()
	draft := createTestDraft()

	expectedErr := errors.New("database error")
	suite.mockRepo.On("UpdateDraft", ctx, draft).Return(types.Draft{}, expectedErr)

	// Act
	result, err := suite.service.Update(ctx, draft)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
	assert.Equal(suite.T(), types.Draft{}, result)
}

// Run all the test suites
func TestDraftServiceSuites(t *testing.T) {
	suite.Run(t, new(DraftServiceSuccessSuite))
	suite.Run(t, new(DraftServiceErrorSuite))
}
