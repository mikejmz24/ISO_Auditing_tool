package events_test

import (
	"ISO_Auditing_Tool/pkg/events"
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// EventTypesTestSuite is a test suite for the event types functionality
type EventTypesTestSuite struct {
	suite.Suite
}

// TestEventTypeConstants ensures event type constants are defined as expected
func (suite *EventTypesTestSuite) TestEventTypeConstants_HaveExpectedValues() {
	// Data change events
	assert.Equal(suite.T(), events.EventType("data_created"), events.DataCreated)
	assert.Equal(suite.T(), events.EventType("data_updated"), events.DataUpdated)
	assert.Equal(suite.T(), events.EventType("data_deleted"), events.DataDeleted)

	// Materialized query events
	assert.Equal(suite.T(), events.EventType("materialized_query_refresh_request"), events.MaterializedQueryRefreshRequested)
	assert.Equal(suite.T(), events.EventType("materialized_query_created"), events.MaterializedQueryCreated)
	assert.Equal(suite.T(), events.EventType("materialized_query_updated"), events.MaterializedQueryUpdated)
}

// TestNewDataCreatedEvent tests creating a DataCreated event
func (suite *EventTypesTestSuite) TestNewDataCreatedEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	entityType := "user"
	entityID := 123
	affectedQuery := "users_query"

	// Act
	event := events.NewDataCreatedEvent(entityType, entityID, affectedQuery)

	// Assert
	assert.Equal(suite.T(), events.DataCreated, event.Type)

	payload, err := events.GetDataChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entityType, payload.EntityType)
	assert.Equal(suite.T(), entityID, payload.EntityID)
	assert.Equal(suite.T(), "created", payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
}

// TestNewDataUpdatedEvent tests creating a DataUpdated event
func (suite *EventTypesTestSuite) TestNewDataUpdatedEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	entityType := "user"
	entityID := 123
	affectedQuery := "users_query"

	// Act
	event := events.NewDataUpdatedEvent(entityType, entityID, affectedQuery)

	// Assert
	assert.Equal(suite.T(), events.DataUpdated, event.Type)

	payload, err := events.GetDataChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entityType, payload.EntityType)
	assert.Equal(suite.T(), entityID, payload.EntityID)
	assert.Equal(suite.T(), "updated", payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
}

// TestNewDataDeletedEvent tests creating a DataDeleted event
func (suite *EventTypesTestSuite) TestNewDataDeletedEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	entityType := "user"
	entityID := 123
	affectedQuery := "users_query"

	// Act
	event := events.NewDataDeletedEvent(entityType, entityID, affectedQuery)

	// Assert
	assert.Equal(suite.T(), events.DataDeleted, event.Type)

	payload, err := events.GetDataChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entityType, payload.EntityType)
	assert.Equal(suite.T(), entityID, payload.EntityID)
	assert.Equal(suite.T(), "deleted", payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
}

// TestNewMaterializedQueryCreatedEvent tests creating a MaterializedQueryCreated event
func (suite *EventTypesTestSuite) TestNewMaterializedQueryCreatedEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	name := "test_query"
	sql := "SELECT * FROM test"
	data := json.RawMessage("test_data")
	version := 1
	errorCount := 0
	lastError := ""

	// Act
	event := events.NewMaterializedQueryCreatedEvent(name, sql, data, version, errorCount, lastError)

	// Assert
	assert.Equal(suite.T(), events.MaterializedQueryCreated, event.Type)

	payload, err := events.GetMaterializedQueryPayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), name, payload.QueryName)
	assert.Equal(suite.T(), sql, payload.QuerySQL)
	assert.Equal(suite.T(), data, payload.QueryData)
	assert.Equal(suite.T(), version, payload.QueryVersion)
	assert.Equal(suite.T(), errorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), lastError, payload.QueryLastError)
}

// TestNewMaterializedQueryUpdatedEvent tests creating a MaterializedQueryUpdated event
func (suite *EventTypesTestSuite) TestNewMaterializedQueryUpdatedEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	name := "test_query"
	sql := "SELECT * FROM test"
	data := json.RawMessage("test_data")
	version := 1
	errorCount := 0
	lastError := ""

	// Act
	event := events.NewMaterializedQueryUpdatedEvent(name, sql, data, version, errorCount, lastError)

	// Assert
	assert.Equal(suite.T(), events.MaterializedQueryUpdated, event.Type)

	payload, err := events.GetMaterializedQueryPayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), name, payload.QueryName)
	assert.Equal(suite.T(), sql, payload.QuerySQL)
	assert.Equal(suite.T(), data, payload.QueryData)
	assert.Equal(suite.T(), version, payload.QueryVersion)
	assert.Equal(suite.T(), errorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), lastError, payload.QueryLastError)
}

// TestNewMaterializedQueryRefreshEvent tests creating a MaterializedQueryRefreshRequested event
func (suite *EventTypesTestSuite) TestNewMaterializedQueryRefreshEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	name := "test_query"
	sql := "SELECT * FROM test"
	data := json.RawMessage("test_data")
	version := 1
	errorCount := 0
	lastError := ""

	// Act
	event := events.NewMaterializedQueryRefreshEvent(name, sql, data, version, errorCount, lastError)

	// Assert
	assert.Equal(suite.T(), events.MaterializedQueryRefreshRequested, event.Type)

	payload, err := events.GetMaterializedQueryPayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), name, payload.QueryName)
	assert.Equal(suite.T(), sql, payload.QuerySQL)
	assert.Equal(suite.T(), data, payload.QueryData)
	assert.Equal(suite.T(), version, payload.QueryVersion)
	assert.Equal(suite.T(), errorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), lastError, payload.QueryLastError)
}

// TestGetDataChangePayload tests extracting DataChangePayload from events
func (suite *EventTypesTestSuite) TestGetDataChangePayload_WithCorrectEventType_ReturnsPayload() {
	// Arrange
	event := events.NewDataCreatedEvent("user", 123, "users_query")

	// Act
	payload, err := events.GetDataChangePayload(event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "user", payload.EntityType)
	assert.Equal(suite.T(), 123, payload.EntityID)
	assert.Equal(suite.T(), "created", payload.ChangeType)
	assert.Equal(suite.T(), "users_query", payload.AffectedQuery)
}

// TestGetDataChangePayload_WithIncorrectEventType tests error case for wrong event type
func (suite *EventTypesTestSuite) TestGetDataChangePayload_WithIncorrectEventType_ReturnsError() {
	// Arrange
	event := events.NewMaterializedQueryCreatedEvent("test_query", "SELECT * FROM test", []byte("test_data"), 1, 0, "")

	// Act
	_, err := events.GetDataChangePayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "does not use DataChangePayload")
}

// TestGetDataChangePayload_WithIncorrectPayloadType tests error case for wrong payload type
func (suite *EventTypesTestSuite) TestGetDataChangePayload_WithIncorrectPayloadType_ReturnsError() {
	// Arrange - create an event with an incorrect payload
	event := events.Event{
		Type:    events.DataCreated,
		Payload: "not a DataChangePayload",
	}

	// Act
	_, err := events.GetDataChangePayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type")
}

// TestGetMaterializedQueryPayload tests extracting MaterializedQueryPayload from events
func (suite *EventTypesTestSuite) TestGetMaterializedQueryPayload_WithCorrectEventType_ReturnsPayload() {
	// Arrange
	event := events.NewMaterializedQueryCreatedEvent("test_query", "SELECT * FROM test", []byte("test_data"), 1, 0, "")

	// Act
	payload, err := events.GetMaterializedQueryPayload(event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test_query", payload.QueryName)
	assert.Equal(suite.T(), "SELECT * FROM test", payload.QuerySQL)
	// assert.Equal(suite.T(), []byte("test_data"), payload.QueryData)
	assert.Equal(suite.T(), json.RawMessage("test_data"), payload.QueryData)
	assert.Equal(suite.T(), 1, payload.QueryVersion)
	assert.Equal(suite.T(), 0, payload.QueryErrorCount)
	assert.Equal(suite.T(), "", payload.QueryLastError)
}

// TestGetMaterializedQueryPayload_WithIncorrectEventType tests error case for wrong event type
func (suite *EventTypesTestSuite) TestGetMaterializedQueryPayload_WithIncorrectEventType_ReturnsError() {
	// Arrange
	event := events.NewDataCreatedEvent("user", 123, "users_query")

	// Act
	_, err := events.GetMaterializedQueryPayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "does not use MaterializedQueryPayload")
}

// TestGetMaterializedQueryPayload_WithIncorrectPayloadType tests error case for wrong payload type
func (suite *EventTypesTestSuite) TestGetMaterializedQueryPayload_WithIncorrectPayloadType_ReturnsError() {
	// Arrange - create an event with an incorrect payload
	event := events.Event{
		Type:    events.MaterializedQueryCreated,
		Payload: "not a MaterializedQueryPayload",
	}

	// Act
	_, err := events.GetMaterializedQueryPayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type")
}

// TestValidateEventPayload tests the ValidateEventPayload function
func (suite *EventTypesTestSuite) TestValidateEventPayload_WithValidDataChangeEvent_ReturnsNil() {
	// Arrange
	event := events.NewDataCreatedEvent("user", 123, "users_query")

	// Act
	err := events.ValidateEventPayload(event)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestValidateEventPayload with a MaterializedQueryEvent
func (suite *EventTypesTestSuite) TestValidateEventPayload_WithValidMaterializedQueryEvent_ReturnsNil() {
	// Arrange
	event := events.NewMaterializedQueryCreatedEvent("test_query", "SELECT * FROM test", []byte("test_data"), 1, 0, "")

	// Act
	err := events.ValidateEventPayload(event)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestValidateEventPayload with an unknown event type
func (suite *EventTypesTestSuite) TestValidateEventPayload_WithUnknownEventType_ReturnsError() {
	// Arrange
	event := events.Event{
		Type:    "unknown_event_type",
		Payload: "some payload",
	}

	// Act
	err := events.ValidateEventPayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown event type")
}

// TestValidateEventPayload with an incorrect payload type
func (suite *EventTypesTestSuite) TestValidateEventPayload_WithIncorrectPayloadType_ReturnsError() {
	// Arrange
	event := events.Event{
		Type:    events.DataCreated,
		Payload: "not a DataChangePayload",
	}

	// Act
	err := events.ValidateEventPayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type")
}

// Run the test suite
func TestEventTypesSuite(t *testing.T) {
	suite.Run(t, new(EventTypesTestSuite))
}
