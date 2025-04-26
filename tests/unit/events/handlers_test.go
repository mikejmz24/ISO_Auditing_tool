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
	"github.com/stretchr/testify/suite"
)

// EventHandlersTestSuite is a test suite for event handlers functionality
type EventHandlersTestSuite struct {
	suite.Suite
	logOutput *bytes.Buffer
	oldOutput io.Writer // Changed from log.Writer to io.Writer
}

// SetupTest initializes the test environment before each test
func (suite *EventHandlersTestSuite) SetupTest() {
	// Capture log output
	suite.logOutput = new(bytes.Buffer)
	suite.oldOutput = log.Writer()
	log.SetOutput(suite.logOutput)
}

// TearDownTest cleans up after each test
func (suite *EventHandlersTestSuite) TearDownTest() {
	// Restore log output
	log.SetOutput(suite.oldOutput)
}

// TestLoggingHandler_LogsEventInformation tests the LoggingHandler function
func (suite *EventHandlersTestSuite) TestLoggingHandler_LogsEventInformation() {
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

// TestTypedEventHandler_DataChangeHandler tests the TypedEventHandler with DataChangeHandler
func (suite *EventHandlersTestSuite) TestTypedEventHandler_DataChangeHandler_CallsHandlerWithTypedPayload() {
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

// TestTypedEventHandler_MaterializedQueryHandler tests the TypedEventHandler with MaterializedQueryHandler
func (suite *EventHandlersTestSuite) TestTypedEventHandler_MaterializedQueryHandler_CallsHandlerWithTypedPayload() {
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

	event := events.NewMaterializedQueryCreatedEvent("test_query", "SELECT * FROM test", []byte("test_data"), 1, 0, "")
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

// TestTypedEventHandler_FallbackHandler tests the fallback handler for unknown event types
func (suite *EventHandlersTestSuite) TestTypedEventHandler_FallbackHandler_CalledForUnhandledEventTypes() {
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

// TestTypedEventHandler_DataChangeHandlerReturnsError tests error propagation from DataChangeHandler
func (suite *EventHandlersTestSuite) TestTypedEventHandler_DataChangeHandlerReturnsError_PropagatesError() {
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

// TestTypedEventHandler_AsHandler tests conversion of TypedEventHandler to Handler
func (suite *EventHandlersTestSuite) TestTypedEventHandler_AsHandler_ReturnsHandlerFunction() {
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

// TestTypedEventHandler_InvalidPayload tests handling of invalid payloads
func (suite *EventHandlersTestSuite) TestTypedEventHandler_InvalidPayload_LogsErrorButDoesNotCallHandler() {
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

// TestTypedEventHandler_NoHandlersRegistered tests behavior with no handlers
func (suite *EventHandlersTestSuite) TestTypedEventHandler_NoHandlersRegistered_LogsAndReturnsNil() {
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

// TestCreateMaterializedQueryEvent_CreatesEventFromMaterializedQuery tests helper function for MaterializedQuery events
func (suite *EventHandlersTestSuite) TestCreateMaterializedQueryEvent_CreatesEventFromMaterializedQuery() {
	// Create a real MaterializedQuery with json.RawMessage
	materializedQuery := types.MaterializedQuery{
		Name:       "test_query",
		Definition: "SELECT * FROM test",
		Data:       json.RawMessage(`test_data`), // Use json.RawMessage for JSON data
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
	assert.Equal(suite.T(), materializedQuery.Definition, payload.QuerySQL)
	assert.Equal(suite.T(), materializedQuery.Data, payload.QueryData) // Compare json.RawMessage values
	assert.Equal(suite.T(), materializedQuery.Version, payload.QueryVersion)
	assert.Equal(suite.T(), materializedQuery.ErrorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), materializedQuery.LastError, payload.QueryLastError)
}

// TestRefreshMaterializedQueryEvent tests helper function for refresh events
func (suite *EventHandlersTestSuite) TestRefreshMaterializedQueryEvent_CreatesEventFromMaterializedQuery() {
	// Arrange
	materializedQuery := types.MaterializedQuery{
		Name:       "test_query",
		Definition: "SELECT * FROM test",
		Data:       []byte("test_data"),
		Version:    1,
		ErrorCount: 2,
		LastError:  "test error",
	}

	// Act
	event := events.RefreshMaterializedQueryEvent(materializedQuery)

	// Assert
	assert.Equal(suite.T(), events.MaterializedQueryRefreshRequested, event.Type)

	payload, err := events.GetMaterializedQueryPayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), materializedQuery.Name, payload.QueryName)
}

// TestUpdateMaterializedQueryEvent tests helper function for update events
func (suite *EventHandlersTestSuite) TestUpdateMaterializedQueryEvent_CreatesEventFromMaterializedQuery() {
	// Arrange
	materializedQuery := types.MaterializedQuery{
		Name:       "test_query",
		Definition: "SELECT * FROM test",
		Data:       []byte("test_data"),
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

// Run the test suite
func TestEventHandlersSuite(t *testing.T) {
	suite.Run(t, new(EventHandlersTestSuite))
}
