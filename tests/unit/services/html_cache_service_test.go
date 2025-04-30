package services_test

import (
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/services"
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock dependencies
type MockMaterializedHTMLQueryRepository struct {
	mock.Mock
}

func (m *MockMaterializedHTMLQueryRepository) CreateMaterializedHTMLQuery(ctx context.Context, query types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(types.MaterializedHTMLQuery), args.Error(1)
}

func (m *MockMaterializedHTMLQueryRepository) GetByNameMaterializedHTMLQuery(ctx context.Context, name string) (types.MaterializedHTMLQuery, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(types.MaterializedHTMLQuery), args.Error(1)
}

func (m *MockMaterializedHTMLQueryRepository) UpdateMaterializedHTMLQuery(ctx context.Context, query types.MaterializedHTMLQuery) (types.MaterializedHTMLQuery, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(types.MaterializedHTMLQuery), args.Error(1)
}

func (m *MockMaterializedHTMLQueryRepository) GetByIDMaterializedHTMLQuery(ctx context.Context, id int) (types.MaterializedHTMLQuery, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(types.MaterializedHTMLQuery), args.Error(1)
}

func (m *MockMaterializedHTMLQueryRepository) DeleteMaterializedHTMLQuery(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockMaterializedJSONQueryRepository struct {
	mock.Mock
}

func (m *MockMaterializedJSONQueryRepository) GetByIDWithFullHierarchyMaterializedJSONQuery(ctx context.Context, query types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(types.MaterializedJSONQuery), args.Error(1)
}

func (m *MockMaterializedJSONQueryRepository) GetByNameMaterializedJSONQuery(ctx context.Context, name string) (types.MaterializedJSONQuery, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(types.MaterializedJSONQuery), args.Error(1)
}

func (m *MockMaterializedJSONQueryRepository) GetByEntityTypeAndIDMaterializedJSONQuery(ctx context.Context, entityType string, entityID int) (types.MaterializedJSONQuery, error) {
	args := m.Called(ctx, entityType, entityID)
	return args.Get(0).(types.MaterializedJSONQuery), args.Error(1)
}

func (m *MockMaterializedJSONQueryRepository) CreateMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error) {
	args := m.Called(ctx, materializedQuery)
	return args.Get(0).(types.MaterializedJSONQuery), args.Error(1)
}

func (m *MockMaterializedJSONQueryRepository) UpdateMaterializedJSONQuery(ctx context.Context, materializedQuery types.MaterializedJSONQuery) (types.MaterializedJSONQuery, error) {
	args := m.Called(ctx, materializedQuery)
	return args.Get(0).(types.MaterializedJSONQuery), args.Error(1)
}

type MockStandardRepository struct {
	mock.Mock
}

func (m *MockStandardRepository) GetByIDStandard(ctx context.Context, standard types.Standard) (types.Standard, error) {
	args := m.Called(ctx, standard)
	return args.Get(0).(types.Standard), args.Error(1)
}

func (m *MockStandardRepository) GetByIDWithFullHierarchyStandard(ctx context.Context, standard types.Standard) (types.Standard, error) {
	args := m.Called(ctx, standard)
	return args.Get(0).(types.Standard), args.Error(1)
}

func (m *MockStandardRepository) GetAllStandards(ctx context.Context) ([]types.Standard, error) {
	args := m.Called(ctx)
	return args.Get(0).([]types.Standard), args.Error(1)
}

func (m *MockStandardRepository) CreateStandard(ctx context.Context, standard types.Standard) (types.Standard, error) {
	args := m.Called(ctx, standard)
	return args.Get(0).(types.Standard), args.Error(1)
}

func (m *MockStandardRepository) UpdateStandard(ctx context.Context, standard types.Standard) (types.Standard, error) {
	args := m.Called(ctx, standard)
	return args.Get(0).(types.Standard), args.Error(1)
}

func (m *MockStandardRepository) DeleteStandard(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockRequirementRepository struct {
	mock.Mock
}

func (m *MockRequirementRepository) GetByIDRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error) {
	args := m.Called(ctx, requirement)
	return args.Get(0).(types.Requirement), args.Error(1)
}

func (m *MockRequirementRepository) GetByIDWithQuestionsRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error) {
	args := m.Called(ctx, requirement)
	return args.Get(0).(types.Requirement), args.Error(1)
}

// Instead of mocking EventBus, we'll use the real one
// and just verify the results of its operations

// Test Suites
type HTMLCacheServiceSuccessSuite struct {
	suite.Suite
	mockHTMLRepo        *MockMaterializedHTMLQueryRepository
	mockJSONRepo        *MockMaterializedJSONQueryRepository
	mockStandardRepo    *MockStandardRepository
	mockRequirementRepo *MockRequirementRepository
	eventBus            *events.EventBus
	service             *services.HTMLCacheService
}

type HTMLCacheServiceErrorSuite struct {
	suite.Suite
	mockHTMLRepo        *MockMaterializedHTMLQueryRepository
	mockJSONRepo        *MockMaterializedJSONQueryRepository
	mockStandardRepo    *MockStandardRepository
	mockRequirementRepo *MockRequirementRepository
	eventBus            *events.EventBus
	service             *services.HTMLCacheService
}

// SetupTest initializes test dependencies before each test
func (suite *HTMLCacheServiceSuccessSuite) SetupTest() {
	suite.mockHTMLRepo = new(MockMaterializedHTMLQueryRepository)
	suite.mockJSONRepo = new(MockMaterializedJSONQueryRepository)
	suite.mockStandardRepo = new(MockStandardRepository)
	suite.mockRequirementRepo = new(MockRequirementRepository)
	suite.eventBus = events.NewEventBus()

	// We don't need to mock the subscription as we're using a real EventBus

	suite.service = services.NewHTMLCacheService(
		suite.mockHTMLRepo,
		suite.mockJSONRepo,
		suite.mockStandardRepo,
		suite.mockRequirementRepo,
		suite.eventBus,
	)
}

// TearDownTest cleans up after each test
func (suite *HTMLCacheServiceSuccessSuite) TearDownTest() {
	suite.mockHTMLRepo.AssertExpectations(suite.T())
	suite.mockJSONRepo.AssertExpectations(suite.T())
	suite.mockStandardRepo.AssertExpectations(suite.T())
	suite.mockRequirementRepo.AssertExpectations(suite.T())
	// No need to verify event bus expectations as we're using a real one
}

// SetupTest initializes test dependencies before each test
func (suite *HTMLCacheServiceErrorSuite) SetupTest() {
	suite.mockHTMLRepo = new(MockMaterializedHTMLQueryRepository)
	suite.mockJSONRepo = new(MockMaterializedJSONQueryRepository)
	suite.mockStandardRepo = new(MockStandardRepository)
	suite.mockRequirementRepo = new(MockRequirementRepository)
	suite.eventBus = events.NewEventBus()

	// Set up event bus subscriptions expectation
	// suite.mockEventBus.On("Subscribe", events.MaterializedQueryCreated, mock.AnythingOfType("events.Handler")).Return()
	// suite.mockEventBus.On("Subscribe", events.MaterializedQueryUpdated, mock.AnythingOfType("events.Handler")).Return()
	suite.eventBus.Subscribe(events.MaterializedQueryCreated, events.LoggingHandler())
	suite.eventBus.Subscribe(events.MaterializedQueryUpdated, events.LoggingHandler())

	suite.service = services.NewHTMLCacheService(
		suite.mockHTMLRepo,
		suite.mockJSONRepo,
		suite.mockStandardRepo,
		suite.mockRequirementRepo,
		suite.eventBus,
	)
}

// TearDownTest cleans up after each test
func (suite *HTMLCacheServiceErrorSuite) TearDownTest() {
	suite.mockHTMLRepo.AssertExpectations(suite.T())
	suite.mockJSONRepo.AssertExpectations(suite.T())
	suite.mockStandardRepo.AssertExpectations(suite.T())
	suite.mockRequirementRepo.AssertExpectations(suite.T())
	// No need to verify event bus expectations as we're using a real one
}

// Helper functions
func createTestStandard() types.Standard {
	return types.Standard{
		ID:          1,
		Name:        "ISO 27001",
		Description: "Information Security Management",
		Version:     "2022",
		Requirements: []types.Requirement{
			{
				ID:            1,
				StandardID:    1,
				LevelID:       1,
				ParentID:      0,
				ReferenceCode: "A.5.1",
				Name:          "Information Security Policies",
				Description:   "Management direction for information security",
			},
		},
	}
}

func createTestMaterializedJSONQuery() types.MaterializedJSONQuery {
	standardData := createTestStandard()
	jsonData, _ := json.Marshal(standardData)

	return types.MaterializedJSONQuery{
		ID:         1,
		Name:       "standard_full_1",
		EntityType: "standard",
		EntityID:   1,
		Definition: "SELECT * FROM standards WHERE id = 1",
		Data:       jsonData,
		Version:    1,
		ErrorCount: 0,
		LastError:  "",
		CreatedAt:  time.Now(),
		UpdatedAt:  nil,
	}
}

func createTestMaterializedHTMLQuery() types.MaterializedHTMLQuery {
	return types.MaterializedHTMLQuery{
		ID:          1,
		Name:        "audit_view_1",
		ViewPath:    "/web/audits/standard/1",
		HTMLContent: "<div>ISO 27001 Audit View</div>",
		Version:     1,
		ErrorCount:  0,
		LastError:   "",
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
}

// --- Success Test Cases ---

// TestNewHTMLCacheService_ReturnsServiceWithDependencies tests the constructor function
func (suite *HTMLCacheServiceSuccessSuite) TestNewHTMLCacheService_ReturnsServiceWithDependencies() {
	// Already tested in SetupTest, just confirm the service was created
	assert.NotNil(suite.T(), suite.service)
}

// TestGetCachedHTML_WhenHTMLExists_ReturnsHTMLContent tests retrieving cached HTML
func (suite *HTMLCacheServiceSuccessSuite) TestGetCachedHTML_WhenHTMLExists_ReturnsHTMLContent() {
	// Arrange
	ctx := context.Background()
	htmlQuery := createTestMaterializedHTMLQuery()
	viewName := "audit_view_1"

	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, viewName).Return(htmlQuery, nil)

	// Act
	html, found, err := suite.service.GetCachedHTML(ctx, viewName)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), htmlQuery.HTMLContent, html)
}

// TestRefreshHTMLForQuery_WhenStandardQuery_SchedulesUpdate tests refresh HTML functionality
func (suite *HTMLCacheServiceSuccessSuite) TestRefreshHTMLForQuery_WhenStandardQuery_SchedulesUpdate() {
	// Arrange
	ctx := context.Background()
	queryName := "standard_full_1"

	// We can't directly test the debounced function, but we can verify no errors are returned

	// Act
	err := suite.service.RefreshHTMLForQuery(ctx, queryName)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestHandleQueryEvent_WithMaterializedQueryPayload_CallsRefreshHTML tests event handling
func (suite *HTMLCacheServiceSuccessSuite) TestHandleQueryEvent_WithMaterializedQueryPayload_CallsRefreshHTML() {
	// Arrange
	ctx := context.Background()
	payload := events.MaterializedQueryPayload{
		QueryName:    "standard_full_1",
		QuerySQL:     "SELECT * FROM standards WHERE id = 1",
		QueryData:    json.RawMessage(`{"id":1,"name":"ISO 27001"}`),
		QueryVersion: 1,
	}

	// Act
	err := suite.service.HandleQueryEvent(ctx, events.MaterializedQueryUpdated, payload)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestRegenerateHTML_WhenValidStandardID_RegeneratesHTML tests regenerating HTML for a standard
func (suite *HTMLCacheServiceSuccessSuite) TestRegenerateHTML_WhenValidStandardID_RegeneratesHTML() {
	// Arrange
	ctx := context.Background()
	standardID := 1

	// Mocking for regenerateHTMLForStandard which is called by RegenerateHTML
	jsonQuery := createTestMaterializedJSONQuery()
	suite.mockJSONRepo.On("GetByNameMaterializedJSONQuery", ctx, fmt.Sprintf("standard_full_%d", standardID)).Return(jsonQuery, nil)

	standard := createTestStandard()
	suite.mockStandardRepo.On("GetByIDStandard", ctx, types.Standard{ID: standardID}).Return(standard, nil)

	// For audit view generation
	htmlQuery := createTestMaterializedHTMLQuery()
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, fmt.Sprintf("audit_view_%d", standardID)).Return(htmlQuery, nil)

	updatedHTMLQuery := htmlQuery
	updatedHTMLQuery.Version = 2
	suite.mockHTMLRepo.On("UpdateMaterializedHTMLQuery", ctx, mock.MatchedBy(func(q types.MaterializedHTMLQuery) bool {
		return q.Name == fmt.Sprintf("audit_view_%d", standardID) && q.Version == 2
	})).Return(updatedHTMLQuery, nil)

	// For requirements view generation
	reqHTMLQuery := createTestMaterializedHTMLQuery()
	reqHTMLQuery.Name = "requirements_view_1"
	reqHTMLQuery.ViewPath = "/web/requirements/standard/1"
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, fmt.Sprintf("requirements_view_%d", standardID)).Return(reqHTMLQuery, nil)

	updatedReqHTMLQuery := reqHTMLQuery
	updatedReqHTMLQuery.Version = 2
	suite.mockHTMLRepo.On("UpdateMaterializedHTMLQuery", ctx, mock.MatchedBy(func(q types.MaterializedHTMLQuery) bool {
		return q.Name == fmt.Sprintf("requirements_view_%d", standardID) && q.Version == 2
	})).Return(updatedReqHTMLQuery, nil)

	// Act
	err := suite.service.RegenerateHTML(ctx, standardID)

	// Assert
	assert.NoError(suite.T(), err)
}

// --- Error Test Cases ---

// TestGetCachedHTML_WhenHTMLDoesNotExist_ReturnsError tests error when HTML is not found
func (suite *HTMLCacheServiceErrorSuite) TestGetCachedHTML_WhenHTMLDoesNotExist_ReturnsError() {
	// Arrange
	ctx := context.Background()
	viewName := "nonexistent_view"

	expectedErr := errors.New("HTML query not found")
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, viewName).Return(types.MaterializedHTMLQuery{}, expectedErr)

	// Act
	html, found, err := suite.service.GetCachedHTML(ctx, viewName)

	// Assert
	assert.Error(suite.T(), err)
	assert.False(suite.T(), found)
	assert.Equal(suite.T(), "", html)
	assert.Equal(suite.T(), expectedErr, err)
}

// TestRegenerateHTML_WhenJSONQueryNotFound_FetchesStandardDirectly tests fallback behavior
func (suite *HTMLCacheServiceErrorSuite) TestRegenerateHTML_WhenJSONQueryNotFound_FetchesStandardDirectly() {
	// Arrange
	ctx := context.Background()
	standardID := 1

	// Mock JSON query not found
	suite.mockJSONRepo.On("GetByNameMaterializedJSONQuery", ctx, fmt.Sprintf("standard_full_%d", standardID)).
		Return(types.MaterializedJSONQuery{}, errors.New("JSON query not found"))

	// Mock standard being fetched directly
	standard := createTestStandard()
	suite.mockStandardRepo.On("GetByIDWithFullHierarchyStandard", ctx, types.Standard{ID: standardID}).Return(standard, nil)

	// Also mock the GetByIDStandard call that happens in generateAuditView and generateRequirementsView
	suite.mockStandardRepo.On("GetByIDStandard", ctx, types.Standard{ID: standardID}).Return(standard, nil)

	// For audit view generation
	htmlQuery := createTestMaterializedHTMLQuery()
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, fmt.Sprintf("audit_view_%d", standardID)).Return(htmlQuery, nil)

	updatedHTMLQuery := htmlQuery
	updatedHTMLQuery.Version = 2
	suite.mockHTMLRepo.On("UpdateMaterializedHTMLQuery", ctx, mock.MatchedBy(func(q types.MaterializedHTMLQuery) bool {
		return q.Name == fmt.Sprintf("audit_view_%d", standardID) && q.Version == 2
	})).Return(updatedHTMLQuery, nil)

	// For requirements view generation
	reqHTMLQuery := createTestMaterializedHTMLQuery()
	reqHTMLQuery.Name = "requirements_view_1"
	reqHTMLQuery.ViewPath = "/web/requirements/standard/1"
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, fmt.Sprintf("requirements_view_%d", standardID)).Return(reqHTMLQuery, nil)

	updatedReqHTMLQuery := reqHTMLQuery
	updatedReqHTMLQuery.Version = 2
	suite.mockHTMLRepo.On("UpdateMaterializedHTMLQuery", ctx, mock.MatchedBy(func(q types.MaterializedHTMLQuery) bool {
		return q.Name == fmt.Sprintf("requirements_view_%d", standardID) && q.Version == 2
	})).Return(updatedReqHTMLQuery, nil)

	// Act
	err := suite.service.RegenerateHTML(ctx, standardID)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestRegenerateHTML_WhenStandardNotFound_ReturnsError tests error when standard is not found
func (suite *HTMLCacheServiceErrorSuite) TestRegenerateHTML_WhenStandardNotFound_ReturnsError() {
	// Arrange
	ctx := context.Background()
	standardID := 999 // Non-existent ID

	// Mock JSON query not found
	suite.mockJSONRepo.On("GetByNameMaterializedJSONQuery", ctx, fmt.Sprintf("standard_full_%d", standardID)).
		Return(types.MaterializedJSONQuery{}, errors.New("JSON query not found"))

	// Mock standard not found
	expectedErr := errors.New("standard not found")
	suite.mockStandardRepo.On("GetByIDWithFullHierarchyStandard", ctx, types.Standard{ID: standardID}).
		Return(types.Standard{}, expectedErr)

	// Act
	err := suite.service.RegenerateHTML(ctx, standardID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get standard data")
}

// TestRegenerateHTML_WhenAuditViewUpdateFails_ReturnsError tests error when HTML update fails
func (suite *HTMLCacheServiceErrorSuite) TestRegenerateHTML_WhenAuditViewUpdateFails_ReturnsError() {
	// Arrange
	ctx := context.Background()
	standardID := 1

	// Mock JSON query found
	jsonQuery := createTestMaterializedJSONQuery()
	suite.mockJSONRepo.On("GetByNameMaterializedJSONQuery", ctx, fmt.Sprintf("standard_full_%d", standardID)).Return(jsonQuery, nil)

	// Mock standard found
	standard := createTestStandard()
	suite.mockStandardRepo.On("GetByIDStandard", ctx, types.Standard{ID: standardID}).Return(standard, nil)

	// Mock HTML query found but update fails
	htmlQuery := createTestMaterializedHTMLQuery()
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, fmt.Sprintf("audit_view_%d", standardID)).Return(htmlQuery, nil)

	expectedErr := errors.New("update failed")
	suite.mockHTMLRepo.On("UpdateMaterializedHTMLQuery", ctx, mock.MatchedBy(func(q types.MaterializedHTMLQuery) bool {
		return q.Name == fmt.Sprintf("audit_view_%d", standardID) && q.Version == 2
	})).Return(types.MaterializedHTMLQuery{}, expectedErr)

	// Act
	err := suite.service.RegenerateHTML(ctx, standardID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
}

// TestRegenerateHTML_WhenAuditViewCreateRequired_CreatesNewHTMLQuery tests creating new HTML query
func (suite *HTMLCacheServiceErrorSuite) TestRegenerateHTML_WhenAuditViewCreateRequired_CreatesNewHTMLQuery() {
	// Arrange
	ctx := context.Background()
	standardID := 1

	// Mock JSON query found
	jsonQuery := createTestMaterializedJSONQuery()
	suite.mockJSONRepo.On("GetByNameMaterializedJSONQuery", ctx, fmt.Sprintf("standard_full_%d", standardID)).Return(jsonQuery, nil)

	// Mock standard found
	standard := createTestStandard()
	suite.mockStandardRepo.On("GetByIDStandard", ctx, types.Standard{ID: standardID}).Return(standard, nil)

	// Mock HTML query not found, so create is needed
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, fmt.Sprintf("audit_view_%d", standardID)).
		Return(types.MaterializedHTMLQuery{}, errors.New("HTML query not found"))

	newHTMLQuery := createTestMaterializedHTMLQuery()
	suite.mockHTMLRepo.On("CreateMaterializedHTMLQuery", ctx, mock.MatchedBy(func(q types.MaterializedHTMLQuery) bool {
		return q.Name == fmt.Sprintf("audit_view_%d", standardID) && q.Version == 1
	})).Return(newHTMLQuery, nil)

	// Mock requirements view not found and created successfully
	suite.mockHTMLRepo.On("GetByNameMaterializedHTMLQuery", ctx, fmt.Sprintf("requirements_view_%d", standardID)).
		Return(types.MaterializedHTMLQuery{}, errors.New("HTML query not found"))

	newReqHTMLQuery := createTestMaterializedHTMLQuery()
	newReqHTMLQuery.Name = "requirements_view_1"
	newReqHTMLQuery.ViewPath = "/web/requirements/standard/1"
	suite.mockHTMLRepo.On("CreateMaterializedHTMLQuery", ctx, mock.MatchedBy(func(q types.MaterializedHTMLQuery) bool {
		return q.Name == fmt.Sprintf("requirements_view_%d", standardID) && q.Version == 1
	})).Return(newReqHTMLQuery, nil)

	// Act
	err := suite.service.RegenerateHTML(ctx, standardID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Run all the test suites
func TestHTMLCacheServiceSuites(t *testing.T) {
	suite.Run(t, new(HTMLCacheServiceSuccessSuite))
	suite.Run(t, new(HTMLCacheServiceErrorSuite))
}
