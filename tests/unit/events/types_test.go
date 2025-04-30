package events_test

import (
	"ISO_Auditing_Tool/pkg/events"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// EventTypeConstantsSuite tests the constants defined in the events package
type EventTypeConstantsSuite struct {
	suite.Suite
}

// EventCreationSuite tests event creation functions
type EventCreationSuite struct {
	suite.Suite
}

// PayloadExtractionSuite tests functions that extract payloads from events
type PayloadExtractionSuite struct {
	suite.Suite
}

// PayloadValidationSuite tests the ValidateEventPayload function
type PayloadValidationSuite struct {
	suite.Suite
}

// --- EventTypeConstantsSuite Tests ---

// TestEventTypeConstants_HaveExpectedValues verifies event type constants are defined as expected
func (suite *EventTypeConstantsSuite) TestEventTypeConstants_HaveExpectedValues() {
	// Data change events
	assert.Equal(suite.T(), events.EventType("data_created"), events.DataCreated)
	assert.Equal(suite.T(), events.EventType("data_updated"), events.DataUpdated)
	assert.Equal(suite.T(), events.EventType("data_deleted"), events.DataDeleted)

	// Entity change event
	assert.Equal(suite.T(), events.EventType("entity_changed"), events.EntityChanged)

	// Materialized query events
	assert.Equal(suite.T(), events.EventType("materialized_query_refresh_request"), events.MaterializedQueryRefreshRequested)
	assert.Equal(suite.T(), events.EventType("materialized_query_created"), events.MaterializedQueryCreated)
	assert.Equal(suite.T(), events.EventType("materialized_query_updated"), events.MaterializedQueryUpdated)
}

// TestEntityTypeConstants_HaveExpectedValues verifies entity type constants are defined as expected
func (suite *EventTypeConstantsSuite) TestEntityTypeConstants_HaveExpectedValues() {
	assert.Equal(suite.T(), events.EntityType("standard"), events.EntityStandard)
	assert.Equal(suite.T(), events.EntityType("requirement"), events.EntityRequirement)
	assert.Equal(suite.T(), events.EntityType("question"), events.EntityQuestion)
	assert.Equal(suite.T(), events.EntityType("evidence"), events.EntityEvidence)
}

// TestChangeTypeConstants_HaveExpectedValues verifies change type constants are defined as expected
func (suite *EventTypeConstantsSuite) TestChangeTypeConstants_HaveExpectedValues() {
	assert.Equal(suite.T(), events.ChangeType("created"), events.ChangeCreated)
	assert.Equal(suite.T(), events.ChangeType("updated"), events.ChangeUpdated)
	assert.Equal(suite.T(), events.ChangeType("deleted"), events.ChangeDeleted)
}

// --- EventCreationSuite Tests ---

// TestNewDataCreatedEvent_CreatesEventWithCorrectTypeAndPayload tests creating a DataCreated event
func (suite *EventCreationSuite) TestNewDataCreatedEvent_CreatesEventWithCorrectTypeAndPayload() {
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

// TestNewDataUpdatedEvent_CreatesEventWithCorrectTypeAndPayload tests creating a DataUpdated event
func (suite *EventCreationSuite) TestNewDataUpdatedEvent_CreatesEventWithCorrectTypeAndPayload() {
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

// TestNewDataDeletedEvent_CreatesEventWithCorrectTypeAndPayload tests creating a DataDeleted event
func (suite *EventCreationSuite) TestNewDataDeletedEvent_CreatesEventWithCorrectTypeAndPayload() {
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

// TestNewEntityChangeEvent_CreatesEventWithCorrectTypeAndPayload tests the generic entity change event creation
func (suite *EventCreationSuite) TestNewEntityChangeEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	entityType := events.EntityStandard
	entityID := 123
	changeType := events.ChangeCreated
	affectedQuery := "standards_query"
	parentType := events.EntityType("")
	parentID := any(nil)
	data := map[string]any{"name": "Standard 1"}

	// Act
	event := events.NewEntityChangeEvent(
		entityType,
		entityID,
		changeType,
		affectedQuery,
		parentType,
		parentID,
		data,
	)

	// Assert
	assert.Equal(suite.T(), events.EntityChanged, event.Type)

	payload, err := events.GetEntityChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entityType, payload.EntityType)
	assert.Equal(suite.T(), entityID, payload.EntityID)
	assert.Equal(suite.T(), changeType, payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
	assert.Equal(suite.T(), parentType, payload.ParentType)
	assert.Equal(suite.T(), parentID, payload.ParentID)
	assert.Equal(suite.T(), data, payload.Data)
}

// TestNewStandardEvent_CreatesEntityChangeEventWithStandardType tests standard events
func (suite *EventCreationSuite) TestNewStandardEvent_CreatesEntityChangeEventWithStandardType() {
	// Arrange
	standardID := 123
	changeType := events.ChangeCreated
	affectedQuery := "standards_query"
	data := map[string]any{"name": "Standard 1"}

	// Act
	event := events.NewStandardEvent(standardID, changeType, affectedQuery, data)

	// Assert
	assert.Equal(suite.T(), events.EntityChanged, event.Type)

	payload, err := events.GetEntityChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), events.EntityStandard, payload.EntityType)
	assert.Equal(suite.T(), standardID, payload.EntityID)
	assert.Equal(suite.T(), changeType, payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
	assert.Equal(suite.T(), data, payload.Data)
}

// TestNewRequirementEvent_CreatesEntityChangeEventWithRequirementType tests requirement events
func (suite *EventCreationSuite) TestNewRequirementEvent_CreatesEntityChangeEventWithRequirementType() {
	// Arrange
	requirementID := 123
	standardID := 456
	changeType := events.ChangeCreated
	affectedQuery := "requirements_query"
	data := map[string]any{"name": "Requirement 1"}

	// Act
	event := events.NewRequirementEvent(requirementID, changeType, standardID, affectedQuery, data)

	// Assert
	assert.Equal(suite.T(), events.EntityChanged, event.Type)

	payload, err := events.GetEntityChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), events.EntityRequirement, payload.EntityType)
	assert.Equal(suite.T(), requirementID, payload.EntityID)
	assert.Equal(suite.T(), events.EntityStandard, payload.ParentType)
	assert.Equal(suite.T(), standardID, payload.ParentID)
	assert.Equal(suite.T(), changeType, payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
	assert.Equal(suite.T(), data, payload.Data)
}

// TestNewQuestionEvent_CreatesEntityChangeEventWithQuestionType tests question events
func (suite *EventCreationSuite) TestNewQuestionEvent_CreatesEntityChangeEventWithQuestionType() {
	// Arrange
	questionID := 123
	requirementID := 456
	changeType := events.ChangeCreated
	affectedQuery := "questions_query"
	data := map[string]any{"text": "Question 1?"}

	// Act
	event := events.NewQuestionEvent(questionID, changeType, requirementID, affectedQuery, data)

	// Assert
	assert.Equal(suite.T(), events.EntityChanged, event.Type)

	payload, err := events.GetEntityChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), events.EntityQuestion, payload.EntityType)
	assert.Equal(suite.T(), questionID, payload.EntityID)
	assert.Equal(suite.T(), events.EntityRequirement, payload.ParentType)
	assert.Equal(suite.T(), requirementID, payload.ParentID)
	assert.Equal(suite.T(), changeType, payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
	assert.Equal(suite.T(), data, payload.Data)
}

// TestNewEvidenceEvent_CreatesEntityChangeEventWithEvidenceType tests evidence events
func (suite *EventCreationSuite) TestNewEvidenceEvent_CreatesEntityChangeEventWithEvidenceType() {
	// Arrange
	evidenceID := 123
	questionID := 456
	changeType := events.ChangeCreated
	affectedQuery := "evidence_query"
	data := map[string]any{"expected": "Evidence content"}

	// Act
	event := events.NewEvidenceEvent(evidenceID, changeType, questionID, affectedQuery, data)

	// Assert
	assert.Equal(suite.T(), events.EntityChanged, event.Type)

	payload, err := events.GetEntityChangePayload(event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), events.EntityEvidence, payload.EntityType)
	assert.Equal(suite.T(), evidenceID, payload.EntityID)
	assert.Equal(suite.T(), events.EntityQuestion, payload.ParentType)
	assert.Equal(suite.T(), questionID, payload.ParentID)
	assert.Equal(suite.T(), changeType, payload.ChangeType)
	assert.Equal(suite.T(), affectedQuery, payload.AffectedQuery)
	assert.Equal(suite.T(), data, payload.Data)
}

// TestNewMaterializedQueryCreatedEvent_CreatesEventWithCorrectTypeAndPayload tests creating a MaterializedQueryCreated event
func (suite *EventCreationSuite) TestNewMaterializedQueryCreatedEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	name := "test_query"
	sql := "SELECT * FROM test"
	data := json.RawMessage(`{"test":"data"}`)
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
	assert.Equal(suite.T(), string(data), string(payload.QueryData))
	assert.Equal(suite.T(), version, payload.QueryVersion)
	assert.Equal(suite.T(), errorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), lastError, payload.QueryLastError)
}

// TestNewMaterializedQueryUpdatedEvent_CreatesEventWithCorrectTypeAndPayload tests creating a MaterializedQueryUpdated event
func (suite *EventCreationSuite) TestNewMaterializedQueryUpdatedEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	name := "test_query"
	sql := "SELECT * FROM test"
	data := json.RawMessage(`{"test":"data"}`)
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
	assert.Equal(suite.T(), string(data), string(payload.QueryData))
	assert.Equal(suite.T(), version, payload.QueryVersion)
	assert.Equal(suite.T(), errorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), lastError, payload.QueryLastError)
}

// TestNewMaterializedQueryRefreshEvent_CreatesEventWithCorrectTypeAndPayload tests creating a MaterializedQueryRefreshRequested event
func (suite *EventCreationSuite) TestNewMaterializedQueryRefreshEvent_CreatesEventWithCorrectTypeAndPayload() {
	// Arrange
	name := "test_query"
	sql := "SELECT * FROM test"
	data := json.RawMessage(`{"test":"data"}`)
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
	assert.Equal(suite.T(), string(data), string(payload.QueryData))
	assert.Equal(suite.T(), version, payload.QueryVersion)
	assert.Equal(suite.T(), errorCount, payload.QueryErrorCount)
	assert.Equal(suite.T(), lastError, payload.QueryLastError)
}

// --- PayloadExtractionSuite Tests ---

// TestGetEntityChangePayload_WithCorrectEventType_ReturnsPayload tests extracting EntityChangePayload from events
func (suite *PayloadExtractionSuite) TestGetEntityChangePayload_WithCorrectEventType_ReturnsPayload() {
	// Arrange
	payload := events.EntityChangePayload{
		EntityType:    events.EntityStandard,
		EntityID:      123,
		ChangeType:    events.ChangeCreated,
		AffectedQuery: "standards_query",
	}
	event := events.Event{
		Type:    events.EntityChanged,
		Payload: payload,
	}

	// Act
	extractedPayload, err := events.GetEntityChangePayload(event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), payload, extractedPayload)
}

// TestGetEntityChangePayload_WithIncorrectEventType_ReturnsError tests error case for wrong event type
func (suite *PayloadExtractionSuite) TestGetEntityChangePayload_WithIncorrectEventType_ReturnsError() {
	// Arrange
	event := events.NewDataCreatedEvent("user", 123, "users_query")

	// Act
	_, err := events.GetEntityChangePayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "event type data_created does not use EntityChangePayload")
}

// TestGetEntityChangePayload_WithIncorrectPayloadType_ReturnsError tests error case for wrong payload type
func (suite *PayloadExtractionSuite) TestGetEntityChangePayload_WithIncorrectPayloadType_ReturnsError() {
	// Arrange - create an event with an incorrect payload
	event := events.Event{
		Type:    events.EntityChanged,
		Payload: "not an EntityChangePayload",
	}

	// Act
	_, err := events.GetEntityChangePayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type for event entity_changed")
}

// TestGetDataChangePayload_WithCorrectEventType_ReturnsPayload tests extracting DataChangePayload from events
func (suite *PayloadExtractionSuite) TestGetDataChangePayload_WithCorrectEventType_ReturnsPayload() {
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

// TestGetDataChangePayload_WithEntityChangedEvent_ConvertsToLegacyFormat tests conversion from EntityChangePayload
func (suite *PayloadExtractionSuite) TestGetDataChangePayload_WithEntityChangedEvent_ConvertsToLegacyFormat() {
	// Arrange
	entityPayload := events.EntityChangePayload{
		EntityType:    events.EntityStandard,
		EntityID:      123,
		ChangeType:    events.ChangeCreated,
		AffectedQuery: "standards_query",
	}
	event := events.Event{
		Type:    events.EntityChanged,
		Payload: entityPayload,
	}

	// Act
	payload, err := events.GetDataChangePayload(event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), string(entityPayload.EntityType), payload.EntityType)
	assert.Equal(suite.T(), entityPayload.EntityID, payload.EntityID)
	assert.Equal(suite.T(), string(entityPayload.ChangeType), payload.ChangeType)
	assert.Equal(suite.T(), entityPayload.AffectedQuery, payload.AffectedQuery)
}

// TestGetDataChangePayload_WithIncorrectEventType_ReturnsError tests error case for wrong event type
func (suite *PayloadExtractionSuite) TestGetDataChangePayload_WithIncorrectEventType_ReturnsError() {
	// Arrange
	event := events.NewMaterializedQueryCreatedEvent("test_query", "SELECT * FROM test", json.RawMessage(`{"test":"data"}`), 1, 0, "")

	// Act
	_, err := events.GetDataChangePayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "event type materialized_query_created does not use DataChangePayload")
}

// TestGetDataChangePayload_WithIncorrectPayloadType_ReturnsError tests error case for wrong payload type
func (suite *PayloadExtractionSuite) TestGetDataChangePayload_WithIncorrectPayloadType_ReturnsError() {
	// Arrange - create an event with an incorrect payload
	event := events.Event{
		Type:    events.DataCreated,
		Payload: "not a DataChangePayload",
	}

	// Act
	_, err := events.GetDataChangePayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type for event data_created")
}

// TestGetMaterializedQueryPayload_WithCorrectEventType_ReturnsPayload tests extracting MaterializedQueryPayload
func (suite *PayloadExtractionSuite) TestGetMaterializedQueryPayload_WithCorrectEventType_ReturnsPayload() {
	// Arrange
	event := events.NewMaterializedQueryCreatedEvent(
		"test_query",
		"SELECT * FROM test",
		json.RawMessage(`{"test":"data"}`),
		1,
		0,
		"",
	)

	// Act
	payload, err := events.GetMaterializedQueryPayload(event)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test_query", payload.QueryName)
	assert.Equal(suite.T(), "SELECT * FROM test", payload.QuerySQL)
	assert.Equal(suite.T(), `{"test":"data"}`, string(payload.QueryData))
	assert.Equal(suite.T(), 1, payload.QueryVersion)
	assert.Equal(suite.T(), 0, payload.QueryErrorCount)
	assert.Equal(suite.T(), "", payload.QueryLastError)
}

// TestGetMaterializedQueryPayload_WithIncorrectEventType_ReturnsError tests error case for wrong event type
func (suite *PayloadExtractionSuite) TestGetMaterializedQueryPayload_WithIncorrectEventType_ReturnsError() {
	// Arrange
	event := events.NewDataCreatedEvent("user", 123, "users_query")

	// Act
	_, err := events.GetMaterializedQueryPayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "event type data_created does not use MaterializedQueryPayload")
}

// TestGetMaterializedQueryPayload_WithIncorrectPayloadType_ReturnsError tests error case for wrong payload type
func (suite *PayloadExtractionSuite) TestGetMaterializedQueryPayload_WithIncorrectPayloadType_ReturnsError() {
	// Arrange - create an event with an incorrect payload
	event := events.Event{
		Type:    events.MaterializedQueryCreated,
		Payload: "not a MaterializedQueryPayload",
	}

	// Act
	_, err := events.GetMaterializedQueryPayload(event)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid payload type for event materialized_query_created")
}

// --- PayloadValidationSuite Tests ---

// TestValidateEventPayload_WithValidDataChangeEvent_ReturnsNil tests validation of DataChangePayload
func (suite *PayloadValidationSuite) TestValidateEventPayload_WithValidDataChangeEvent_ReturnsNil() {
	// Arrange
	event := events.NewDataCreatedEvent("user", 123, "users_query")

	// Act
	err := events.ValidateEventPayload(event)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestValidateEventPayload_WithValidMaterializedQueryEvent_ReturnsNil tests validation of MaterializedQueryPayload
func (suite *PayloadValidationSuite) TestValidateEventPayload_WithValidMaterializedQueryEvent_ReturnsNil() {
	// Arrange
	event := events.NewMaterializedQueryCreatedEvent(
		"test_query",
		"SELECT * FROM test",
		json.RawMessage(`{"test":"data"}`),
		1,
		0,
		"",
	)

	// Act
	err := events.ValidateEventPayload(event)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestValidateEventPayload_WithUnknownEventType_ReturnsError tests validation with unknown event type
func (suite *PayloadValidationSuite) TestValidateEventPayload_WithUnknownEventType_ReturnsError() {
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

// TestValidateEventPayload_WithIncorrectPayloadType_ReturnsError tests validation with incorrect payload type
func (suite *PayloadValidationSuite) TestValidateEventPayload_WithIncorrectPayloadType_ReturnsError() {
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

// Run all test suites
func TestTypesSuites(t *testing.T) {
	suite.Run(t, new(EventTypeConstantsSuite))
	suite.Run(t, new(EventCreationSuite))
	suite.Run(t, new(PayloadExtractionSuite))
	suite.Run(t, new(PayloadValidationSuite))
}
