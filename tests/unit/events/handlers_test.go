package events_test

import (
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/types"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock service implementations
type MockEntityService struct {
	mock.Mock
}

func (m *MockEntityService) HandleEntityChange(ctx context.Context, payload events.EntityChangePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

type MockMaterializedQueryService struct {
	mock.Mock
}

func (m *MockMaterializedQueryService) HandleQueryEvent(ctx context.Context, eventType events.EventType, payload events.MaterializedQueryPayload) error {
	args := m.Called(ctx, eventType, payload)
	return args.Error(0)
}

type MockHTMLCacheService struct {
	mock.Mock
}

func (m *MockHTMLCacheService) RefreshHTMLForQuery(ctx context.Context, queryName string) error {
	args := m.Called(ctx, queryName)
	return args.Error(0)
}

// Suite for testing logging handlers
type LoggingHandlerSuite struct {
	suite.Suite
	logOutput *bytes.Buffer
	oldOutput io.Writer
}

// Suite for testing entity change handlers
type EntityChangeHandlerSuite struct {
	suite.Suite
	mockEntityService *MockEntityService
}

// Suite for testing materialized query handlers
type MaterializedQueryHandlerSuite struct {
	suite.Suite
	mockQueryService *MockMaterializedQueryService
}

// Suite for testing HTML cache handlers
type HTMLCacheHandlerSuite struct {
	suite.Suite
	mockHTMLCacheService *MockHTMLCacheService
}

// Suite for testing typed event handlers
type TypedEventHandlerSuite struct {
	suite.Suite
	logOutput *bytes.Buffer
	oldOutput io.Writer
}

// SetupTest initializes the test environment before each test
func (suite *LoggingHandlerSuite) SetupTest() {
	// Capture log output
	suite.logOutput = new(bytes.Buffer)
	suite.oldOutput = log.Writer()
	log.SetOutput(suite.logOutput)
}

// TearDownTest cleans up after each test
func (suite *LoggingHandlerSuite) TearDownTest() {
	// Restore log output
	log.SetOutput(suite.oldOutput)
}

// SetupTest for EntityChangeHandlerSuite
func (suite *EntityChangeHandlerSuite) SetupTest() {
	suite.mockEntityService = new(MockEntityService)
}

// SetupTest for MaterializedQueryHandlerSuite
func (suite *MaterializedQueryHandlerSuite) SetupTest() {
	suite.mockQueryService = new(MockMaterializedQueryService)
}

// SetupTest for HTMLCacheHandlerSuite
func (suite *HTMLCacheHandlerSuite) SetupTest() {
	suite.mockHTMLCacheService = new(MockHTMLCacheService)
}

// SetupTest initializes the test environment before TypedEventHandlerSuite tests
func (suite *TypedEventHandlerSuite) SetupTest() {
	// Capture log output
	suite.logOutput = new(bytes.Buffer)
	suite.oldOutput = log.Writer()
	log.SetOutput(suite.logOutput)
}

// TearDownTest cleans up after TypedEventHandlerSuite tests
func (suite *TypedEventHandlerSuite) TearDownTest() {
	// Restore log output
	log.SetOutput(suite.oldOutput)
}

// --- Logging Handler Tests ---

// TestLoggingHandler_WhenEventReceived_LogsEventInformation tests the LoggingHandler function
func (suite *LoggingHandlerSuite) TestLoggingHandler_WhenEventReceived_LogsEventInformation() {
	// Arrange
	handler := events.LoggingHandler()
	event := events.NewDataCreatedEvent("user", 123, "users_query")
	ctx := context.Background()

	// Act
	err := handler(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	logOutput := suite.logOutput.String()
	assert.Contains(suite.T(), logOutput, "Event: data_created")
	assert.Contains(suite.T(), logOutput, "Payload type: events.DataChangePayload")
}

// --- Entity Change Handler Tests ---

// TestEntityChangeHandler_WhenEntityChangedEvent_CallsServiceWithPayload tests handling of entity change events
func (suite *EntityChangeHandlerSuite) TestEntityChangeHandler_WhenEntityChangedEvent_CallsServiceWithPayload() {
	// Arrange
	ctx := context.Background()
	payload := events.EntityChangePayload{
		EntityType: events.EntityStandard,
		EntityID:   123,
		ChangeType: events.ChangeCreated,
	}
	event := events.Event{
		Type:    events.EntityChanged,
		Payload: payload,
	}

	suite.mockEntityService.On("HandleEntityChange", ctx, payload).Return(nil)
	handler := events.EntityChangeHandler(suite.mockEntityService)

	// Act
	err := handler(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockEntityService.AssertExpectations(suite.T())
}

// TestEntityChangeHandler_WhenLegacyDataCreatedEvent_ConvertsAndCallsService tests handling legacy event types
func (suite *EntityChangeHandlerSuite) TestEntityChangeHandler_WhenLegacyDataCreatedEvent_ConvertsAndCallsService() {
	// Arrange
	ctx := context.Background()
	event := events.NewDataCreatedEvent("standard", 123, "standards_query")

	// The payload that should be converted to by the handler
	expectedPayload := events.EntityChangePayload{
		EntityType:    events.EntityType("standard"),
		EntityID:      123,
		ChangeType:    events.ChangeType("created"),
		AffectedQuery: "standards_query",
	}

	suite.mockEntityService.On("HandleEntityChange", ctx, mock.MatchedBy(func(p events.EntityChangePayload) bool {
		return p.EntityType == expectedPayload.EntityType &&
			p.EntityID == expectedPayload.EntityID &&
			p.ChangeType == expectedPayload.ChangeType &&
			p.AffectedQuery == expectedPayload.AffectedQuery
	})).Return(nil)

	handler := events.EntityChangeHandler(suite.mockEntityService)

	// Act
	err := handler(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockEntityService.AssertExpectations(suite.T())
}

// TestEntityChangeHandler_WhenInvalidPayload_ReturnsError tests handling of invalid payloads
func (suite *EntityChangeHandlerSuite) TestEntityChangeHandler_WhenInvalidPayload_ReturnsError() {
	// Arrange
	ctx := context.Background()
	invalidEvent := events.Event{
		Type:    events.EntityChanged,
		Payload: "invalid payload",
	}

	handler := events.EntityChangeHandler(suite.mockEntityService)

	// Act
	err := handler(ctx, invalidEvent)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type")
	// Service should not be called with invalid payload
	suite.mockEntityService.AssertNotCalled(suite.T(), "HandleEntityChange")
}

// TestEntityChangeHandler_WhenServiceReturnsError_PropagatesError tests error propagation
func (suite *EntityChangeHandlerSuite) TestEntityChangeHandler_WhenServiceReturnsError_PropagatesError() {
	// Arrange
	ctx := context.Background()
	expectedError := errors.New("service error")
	payload := events.EntityChangePayload{
		EntityType: events.EntityStandard,
		EntityID:   123,
		ChangeType: events.ChangeCreated,
	}
	event := events.Event{
		Type:    events.EntityChanged,
		Payload: payload,
	}

	suite.mockEntityService.On("HandleEntityChange", ctx, payload).Return(expectedError)
	handler := events.EntityChangeHandler(suite.mockEntityService)

	// Act
	err := handler(ctx, event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockEntityService.AssertExpectations(suite.T())
}

// --- Materialized Query Handler Tests ---

// TestMaterializedQueryHandler_WhenValidEvent_CallsServiceWithPayload tests materialized query handler
func (suite *MaterializedQueryHandlerSuite) TestMaterializedQueryHandler_WhenValidEvent_CallsServiceWithPayload() {
	// Arrange
	ctx := context.Background()
	payload := events.MaterializedQueryPayload{
		QueryName: "test_query",
		QuerySQL:  "SELECT * FROM test",
	}
	event := events.Event{
		Type:    events.MaterializedQueryCreated,
		Payload: payload,
	}

	suite.mockQueryService.On("HandleQueryEvent", ctx, events.MaterializedQueryCreated, payload).Return(nil)
	handler := events.MaterializedQueryHandler(suite.mockQueryService)

	// Act
	err := handler(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockQueryService.AssertExpectations(suite.T())
}

// TestMaterializedQueryHandler_WhenInvalidPayload_ReturnsError tests handling of invalid payloads
func (suite *MaterializedQueryHandlerSuite) TestMaterializedQueryHandler_WhenInvalidPayload_ReturnsError() {
	// Arrange
	ctx := context.Background()
	invalidEvent := events.Event{
		Type:    events.MaterializedQueryCreated,
		Payload: "invalid payload",
	}

	handler := events.MaterializedQueryHandler(suite.mockQueryService)

	// Act
	err := handler(ctx, invalidEvent)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type")
	// Service should not be called with invalid payload
	suite.mockQueryService.AssertNotCalled(suite.T(), "HandleQueryEvent")
}

// --- HTML Cache Handler Tests ---

// TestHTMLCacheHandler_WhenValidEvent_CallsServiceWithQueryName tests HTML cache handler
func (suite *HTMLCacheHandlerSuite) TestHTMLCacheHandler_WhenValidEvent_CallsServiceWithQueryName() {
	// Arrange
	ctx := context.Background()
	payload := events.MaterializedQueryPayload{
		QueryName: "test_query",
		QuerySQL:  "SELECT * FROM test",
	}
	event := events.Event{
		Type:    events.MaterializedQueryUpdated,
		Payload: payload,
	}

	suite.mockHTMLCacheService.On("RefreshHTMLForQuery", ctx, "test_query").Return(nil)
	handler := events.HTMLCacheHandler(suite.mockHTMLCacheService)

	// Act
	err := handler(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockHTMLCacheService.AssertExpectations(suite.T())
}

// TestHTMLCacheHandler_WhenInvalidPayload_ReturnsError tests handling of invalid payloads
func (suite *HTMLCacheHandlerSuite) TestHTMLCacheHandler_WhenInvalidPayload_ReturnsError() {
	// Arrange
	ctx := context.Background()
	invalidEvent := events.Event{
		Type:    events.MaterializedQueryUpdated,
		Payload: "invalid payload",
	}

	handler := events.HTMLCacheHandler(suite.mockHTMLCacheService)

	// Act
	err := handler(ctx, invalidEvent)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type")
	// Service should not be called with invalid payload
	suite.mockHTMLCacheService.AssertNotCalled(suite.T(), "RefreshHTMLForQuery")
}

// --- Event Creation Helper Function Tests ---

// TestCreateMaterializedQueryEvent_WhenGivenMaterializedQuery_CreatesCorrectEvent tests event creation from MaterializedQuery
func (suite *TypedEventHandlerSuite) TestCreateMaterializedQueryEvent_WhenGivenMaterializedQuery_CreatesCorrectEvent() {
	// Arrange
	materializedQuery := types.MaterializedJSONQuery{
		Name:       "test_query",
		Definition: "SELECT * FROM test",
		Data:       json.RawMessage(`{"test":"data"}`),
		Version:    1,
		ErrorCount: 2,
		LastError:  "test error",
	}

	// Act
	event := events.CreateMaterializedQueryEvent(materializedQuery)

	// Assert
	assert.Equal(suite.T(), events.MaterializedQueryCreated, event.Type)

	payload, err := events.GetMaterializedQueryPayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), materializedQuery.Name, payload.QueryName)
	assert.Equal(suite.T(), materializedQuery.Data, payload.QueryData)
	assert.Equal(suite.T(), materializedQuery.Version, payload.QueryVersion)
	assert.Equal(suite.T(), materializedQuery.ErrorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), materializedQuery.LastError, payload.QueryLastError)
}

// TestUpdateMaterializedQueryEvent_WhenGivenMaterializedQuery_CreatesCorrectEvent tests update event creation
func (suite *TypedEventHandlerSuite) TestUpdateMaterializedQueryEvent_WhenGivenMaterializedQuery_CreatesCorrectEvent() {
	// Arrange
	materializedQuery := types.MaterializedJSONQuery{
		Name:       "test_query",
		Definition: "SELECT * FROM test",
		Data:       json.RawMessage(`{"test":"data"}`),
		Version:    1,
		ErrorCount: 2,
		LastError:  "test error",
	}

	// Act
	event := events.UpdateMaterializedQueryEvent(materializedQuery)

	// Assert
	assert.Equal(suite.T(), events.MaterializedQueryUpdated, event.Type)

	payload, err := events.GetMaterializedQueryPayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), materializedQuery.Name, payload.QueryName)
}

// --- Typed Event Handler Tests ---

// TestTypedEventHandler_WhenDataChangeHandler_CallsHandlerWithTypedPayload tests typed event handler with DataChangeHandler
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenDataChangeHandler_CallsHandlerWithTypedPayload() {
	// Arrange
	handlerCalled := false
	capturedType := events.EventType("")
	capturedEntityType := ""
	capturedEntityID := 0

	handler := events.TypedEventHandler{
		DataChangeHandler: func(ctx context.Context, eventType events.EventType, payload events.DataChangePayload) error {
			handlerCalled = true
			capturedType = eventType
			capturedEntityType = payload.EntityType
			capturedEntityID = payload.EntityID.(int)
			return nil
		},
	}

	event := events.NewDataCreatedEvent("user", 123, "users_query")
	ctx := context.Background()

	// Act
	err := handler.HandleEvent(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), handlerCalled)
	assert.Equal(suite.T(), events.DataCreated, capturedType)
	assert.Equal(suite.T(), "user", capturedEntityType)
	assert.Equal(suite.T(), 123, capturedEntityID)
}

// TestTypedEventHandler_WhenEntityChangeHandler_CallsHandlerWithTypedPayload tests entity change handler
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenEntityChangeHandler_CallsHandlerWithTypedPayload() {
	// Arrange
	handlerCalled := false
	capturedType := events.EventType("")
	capturedEntityType := events.EntityType("")
	capturedEntityID := 0

	handler := events.TypedEventHandler{
		EntityChangeHandler: func(ctx context.Context, eventType events.EventType, payload events.EntityChangePayload) error {
			handlerCalled = true
			capturedType = eventType
			capturedEntityType = payload.EntityType
			capturedEntityID = payload.EntityID.(int)
			return nil
		},
	}

	payload := events.EntityChangePayload{
		EntityType: events.EntityStandard,
		EntityID:   123,
		ChangeType: events.ChangeCreated,
	}
	event := events.Event{
		Type:    events.EntityChanged,
		Payload: payload,
	}
	ctx := context.Background()

	// Act
	err := handler.HandleEvent(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), handlerCalled)
	assert.Equal(suite.T(), events.EntityChanged, capturedType)
	assert.Equal(suite.T(), events.EntityStandard, capturedEntityType)
	assert.Equal(suite.T(), 123, capturedEntityID)
}

// TestTypedEventHandler_WhenMaterializedQueryHandler_CallsHandlerWithTypedPayload tests materialized query handler
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenMaterializedQueryHandler_CallsHandlerWithTypedPayload() {
	// Arrange
	handlerCalled := false
	capturedType := events.EventType("")
	capturedQueryName := ""
	capturedQuerySQL := ""

	handler := events.TypedEventHandler{
		MaterializedQueryHandler: func(ctx context.Context, eventType events.EventType, payload events.MaterializedQueryPayload) error {
			handlerCalled = true
			capturedType = eventType
			capturedQueryName = payload.QueryName
			capturedQuerySQL = payload.QuerySQL
			return nil
		},
	}

	event := events.NewMaterializedQueryCreatedEvent("test_query", "SELECT * FROM test", json.RawMessage(`{"test":"data"}`), 1, 0, "")
	ctx := context.Background()

	// Act
	err := handler.HandleEvent(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), handlerCalled)
	assert.Equal(suite.T(), events.MaterializedQueryCreated, capturedType)
	assert.Equal(suite.T(), "test_query", capturedQueryName)
	assert.Equal(suite.T(), "SELECT * FROM test", capturedQuerySQL)
}

// TestTypedEventHandler_WhenFallbackHandler_CalledForUnhandledEventTypes tests the fallback handler
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenFallbackHandler_CalledForUnhandledEventTypes() {
	// Arrange
	dataChangeCalled := false
	materializedQueryCalled := false
	fallbackCalled := false
	capturedEventType := events.EventType("")

	handler := events.TypedEventHandler{
		DataChangeHandler: func(ctx context.Context, eventType events.EventType, payload events.DataChangePayload) error {
			dataChangeCalled = true
			return nil
		},
		MaterializedQueryHandler: func(ctx context.Context, eventType events.EventType, payload events.MaterializedQueryPayload) error {
			materializedQueryCalled = true
			return nil
		},
		FallbackHandler: func(ctx context.Context, event events.Event) error {
			fallbackCalled = true
			capturedEventType = event.Type
			return nil
		},
	}

	// Create a custom event not handled by typed handlers
	event := events.Event{
		Type:    "custom_event",
		Payload: "custom_payload",
	}
	ctx := context.Background()

	// Act
	err := handler.HandleEvent(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), dataChangeCalled)
	assert.False(suite.T(), materializedQueryCalled)
	assert.True(suite.T(), fallbackCalled)
	assert.Equal(suite.T(), events.EventType("custom_event"), capturedEventType)
}

// TestTypedEventHandler_WhenDataChangeHandlerReturnsError_PropagatesError tests error propagation
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenDataChangeHandlerReturnsError_PropagatesError() {
	// Arrange
	expectedErr := errors.New("data handler error")

	handler := events.TypedEventHandler{
		DataChangeHandler: func(ctx context.Context, eventType events.EventType, payload events.DataChangePayload) error {
			return expectedErr
		},
	}

	event := events.NewDataCreatedEvent("user", 123, "users_query")
	ctx := context.Background()

	// Act
	err := handler.HandleEvent(ctx, event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
}

// TestTypedEventHandler_WhenAsHandlerCalled_ReturnsHandlerFunction tests conversion to Handler
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenAsHandlerCalled_ReturnsHandlerFunction() {
	// Arrange
	handlerCalled := false

	typedHandler := events.TypedEventHandler{
		DataChangeHandler: func(ctx context.Context, eventType events.EventType, payload events.DataChangePayload) error {
			handlerCalled = true
			return nil
		},
	}

	// Act
	handler := typedHandler.AsHandler()

	// Assert
	assert.NotNil(suite.T(), handler)

	// Test the returned handler
	event := events.NewDataCreatedEvent("user", 123, "users_query")
	ctx := context.Background()

	err := handler(ctx, event)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), handlerCalled)
}

// TestTypedEventHandler_WhenInvalidPayload_LogsErrorButDoesNotPanic tests handling of invalid payloads
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenInvalidPayload_LogsErrorButDoesNotPanic() {
	// Arrange
	handlerCalled := false

	handler := events.TypedEventHandler{
		DataChangeHandler: func(ctx context.Context, eventType events.EventType, payload events.DataChangePayload) error {
			handlerCalled = true
			return nil
		},
	}

	// Create an event with invalid payload
	event := events.Event{
		Type:    events.DataCreated,
		Payload: "not a DataChangePayload",
	}
	ctx := context.Background()

	// Act
	err := handler.HandleEvent(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), handlerCalled)
	logOutput := suite.logOutput.String()
	assert.Contains(suite.T(), logOutput, "Error extracting DataChangePayload")
}

// TestTypedEventHandler_WhenNoHandlersRegistered_LogsAndReturnsNil tests behavior with no handlers
func (suite *TypedEventHandlerSuite) TestTypedEventHandler_WhenNoHandlersRegistered_LogsAndReturnsNil() {
	// Arrange
	handler := events.TypedEventHandler{}
	event := events.NewDataCreatedEvent("user", 123, "users_query")
	ctx := context.Background()

	// Act
	err := handler.HandleEvent(ctx, event)

	// Assert
	assert.NoError(suite.T(), err)
	logOutput := suite.logOutput.String()
	assert.Contains(suite.T(), logOutput, "No handler for event type")
}

// Run all test suites
func TestHandlersSuites(t *testing.T) {
	suite.Run(t, new(LoggingHandlerSuite))
	suite.Run(t, new(EntityChangeHandlerSuite))
	suite.Run(t, new(MaterializedQueryHandlerSuite))
	suite.Run(t, new(HTMLCacheHandlerSuite))
	suite.Run(t, new(TypedEventHandlerSuite))
}
