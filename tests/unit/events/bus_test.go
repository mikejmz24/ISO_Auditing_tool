package events_test

import (
	"ISO_Auditing_Tool/pkg/events"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// EventBusHappyPathSuite is a test suite for successful EventBus operations
type EventBusHappyPathSuite struct {
	suite.Suite
	bus *events.EventBus
}

// EventBusErrorHandlingSuite is a test suite for EventBus error handling
type EventBusErrorHandlingSuite struct {
	suite.Suite
	bus *events.EventBus
}

// EventBusAsyncSuite is a test suite for asynchronous EventBus operations
type EventBusAsyncSuite struct {
	suite.Suite
	bus *events.EventBus
}

// SetupTest initializes a new EventBus before each test
func (suite *EventBusHappyPathSuite) SetupTest() {
	suite.bus = events.NewEventBus()
}

// SetupTest initializes a new EventBus before each test
func (suite *EventBusErrorHandlingSuite) SetupTest() {
	suite.bus = events.NewEventBus()
}

// SetupTest initializes a new EventBus before each test
func (suite *EventBusAsyncSuite) SetupTest() {
	suite.bus = events.NewEventBus()
}

// --- Happy Path Tests ---

// TestNewEventBus_ReturnsInitializedInstance tests the constructor function
func (suite *EventBusHappyPathSuite) TestNewEventBus_ReturnsInitializedInstance() {
	// Act
	bus := events.NewEventBus()

	// Assert
	assert.NotNil(suite.T(), bus)
}

// TestSubscribe_AddsHandlerToRegistry tests that Subscribe adds a handler to the registry
func (suite *EventBusHappyPathSuite) TestSubscribe_AddsHandlerToRegistry() {
	// Arrange
	handlerCalled := false
	handler := func(ctx context.Context, event events.Event) error {
		handlerCalled = true
		return nil
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, handler)
	err := suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), handlerCalled)
}

// TestSubscribe_MultipleHandlers_AllReceiveEvent tests that multiple handlers for an event all get called
func (suite *EventBusHappyPathSuite) TestSubscribe_MultipleHandlers_AllReceiveEvent() {
	// Arrange
	firstHandlerCalled := false
	secondHandlerCalled := false

	firstHandler := func(ctx context.Context, event events.Event) error {
		firstHandlerCalled = true
		return nil
	}

	secondHandler := func(ctx context.Context, event events.Event) error {
		secondHandlerCalled = true
		return nil
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, firstHandler)
	suite.bus.Subscribe(events.DataCreated, secondHandler)
	err := suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), firstHandlerCalled)
	assert.True(suite.T(), secondHandlerCalled)
}

// TestSubscribeAll_RegistersHandlerForAllExistingEventTypes tests SubscribeAll functionality
func (suite *EventBusHappyPathSuite) TestSubscribeAll_RegistersHandlerForAllExistingEventTypes() {
	// Arrange
	// Create event type registrations first
	emptyHandler := func(ctx context.Context, event events.Event) error { return nil }
	suite.bus.Subscribe(events.DataCreated, emptyHandler)
	suite.bus.Subscribe(events.DataUpdated, emptyHandler)

	// Track handler calls
	handlerCalls := 0
	handler := func(ctx context.Context, event events.Event) error {
		handlerCalls++
		return nil
	}

	// Act
	suite.bus.SubscribeAll(handler)

	// Publish both event types
	suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))
	suite.bus.Publish(context.Background(), events.NewDataUpdatedEvent("test", 1, ""))

	// Assert
	assert.Equal(suite.T(), 2, handlerCalls)
}

// TestPublish_NoHandlers_ReturnsNil tests that Publish returns nil when no handlers exist
func (suite *EventBusHappyPathSuite) TestPublish_NoHandlers_ReturnsNil() {
	// Act
	err := suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.NoError(suite.T(), err)
}

// TestPublish_PropagatesContextToHandlers tests context propagation in handlers
func (suite *EventBusHappyPathSuite) TestPublish_PropagatesContextToHandlers() {
	// Arrange
	contextReceived := false

	handler := func(ctx context.Context, event events.Event) error {
		// Check if context has our test value
		value := ctx.Value("test").(string)
		if value == "test_value" {
			contextReceived = true
		}
		return nil
	}

	// Create a context with a test value
	ctx := context.WithValue(context.Background(), "test", "test_value")

	// Act
	suite.bus.Subscribe(events.DataCreated, handler)
	err := suite.bus.Publish(ctx, events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), contextReceived)
}

// TestPublish_ExecutesHandlersInSubscriptionOrder tests handler execution order
func (suite *EventBusHappyPathSuite) TestPublish_ExecutesHandlersInSubscriptionOrder() {
	// Arrange
	executionOrder := []int{}

	handler1 := func(ctx context.Context, event events.Event) error {
		executionOrder = append(executionOrder, 1)
		return nil
	}

	handler2 := func(ctx context.Context, event events.Event) error {
		executionOrder = append(executionOrder, 2)
		return nil
	}

	handler3 := func(ctx context.Context, event events.Event) error {
		executionOrder = append(executionOrder, 3)
		return nil
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, handler1)
	suite.bus.Subscribe(events.DataCreated, handler2)
	suite.bus.Subscribe(events.DataCreated, handler3)

	err := suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), []int{1, 2, 3}, executionOrder)
}

// TestPublish_OnlyCallsHandlersForMatchingEventType tests event type isolation
func (suite *EventBusHappyPathSuite) TestPublish_OnlyCallsHandlersForMatchingEventType() {
	// Arrange
	dataCreatedCalled := false
	dataUpdatedCalled := false

	dataCreatedHandler := func(ctx context.Context, event events.Event) error {
		dataCreatedCalled = true
		return nil
	}

	dataUpdatedHandler := func(ctx context.Context, event events.Event) error {
		dataUpdatedCalled = true
		return nil
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, dataCreatedHandler)
	suite.bus.Subscribe(events.DataUpdated, dataUpdatedHandler)

	// Publish a DataCreated event
	err := suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), dataCreatedCalled)
	assert.False(suite.T(), dataUpdatedCalled)
}

// --- Error Handling Tests ---

// TestPublish_HandlerReturnsError_ReturnsFirstError tests error propagation from handlers
func (suite *EventBusErrorHandlingSuite) TestPublish_HandlerReturnsError_ReturnsFirstError() {
	// Arrange
	expectedErr := errors.New("handler error")
	handler := func(ctx context.Context, event events.Event) error {
		return expectedErr
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, handler)
	err := suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), expectedErr.Error())
}

// TestPublish_MultipleHandlersReturnErrors_ReturnsFirstError tests error handling with multiple errors
func (suite *EventBusErrorHandlingSuite) TestPublish_MultipleHandlersReturnErrors_ReturnsFirstError() {
	// Arrange
	expectedErr := errors.New("first handler error")

	firstHandler := func(ctx context.Context, event events.Event) error {
		return expectedErr
	}

	secondHandler := func(ctx context.Context, event events.Event) error {
		return errors.New("second handler error")
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, firstHandler)
	suite.bus.Subscribe(events.DataCreated, secondHandler)
	err := suite.bus.Publish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), expectedErr.Error())
}

// --- Async Tests ---

// TestAsyncPublish_ExecutesHandlerAsynchronously tests async event publishing
func (suite *EventBusAsyncSuite) TestAsyncPublish_ExecutesHandlerAsynchronously() {
	// Arrange
	var wg sync.WaitGroup
	wg.Add(1)

	handler := func(ctx context.Context, event events.Event) error {
		wg.Done()
		return nil
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, handler)
	suite.bus.AsyncPublish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// Assert - wait with timeout
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// Success - handler was called
	case <-time.After(time.Second):
		suite.T().Fatal("Handler was not called within timeout")
	}
}

// TestAsyncPublishWithCallback_HandlerError_CallsCallback tests async publishing with error callback
func (suite *EventBusAsyncSuite) TestAsyncPublishWithCallback_HandlerError_CallsCallback() {
	// Arrange
	expectedErr := errors.New("handler error")

	var wg sync.WaitGroup
	wg.Add(1)

	callbackCalled := false
	callbackEventType := events.EventType("")
	callbackErr := error(nil)

	handler := func(ctx context.Context, event events.Event) error {
		return expectedErr
	}

	callback := func(eventType events.EventType, err error) {
		callbackCalled = true
		callbackEventType = eventType
		callbackErr = err
		wg.Done()
	}

	// Act
	suite.bus.Subscribe(events.DataCreated, handler)
	suite.bus.AsyncPublishWithCallback(
		context.Background(),
		events.NewDataCreatedEvent("test", 1, ""),
		callback,
	)

	// Assert - wait with timeout
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// Continue with assertions
	case <-time.After(time.Second):
		suite.T().Fatal("Callback was not called within timeout")
	}

	assert.True(suite.T(), callbackCalled)
	assert.Equal(suite.T(), events.DataCreated, callbackEventType)
	assert.Contains(suite.T(), callbackErr.Error(), expectedErr.Error())
}

// TestAsyncPublishWithContext_RespectsContext tests that context is propagated correctly in async handlers
func (suite *EventBusAsyncSuite) TestAsyncPublishWithContext_RespectsContext() {
	// Arrange
	var wg sync.WaitGroup
	wg.Add(1)

	handler := func(ctx context.Context, event events.Event) error {
		// Check if context has our test value
		value := ctx.Value("test")
		assert.Equal(suite.T(), "value", value)
		wg.Done()
		return nil
	}

	// Add a value to the context
	ctx := context.WithValue(context.Background(), "test", "value")

	// Act
	suite.bus.Subscribe(events.DataCreated, handler)
	suite.bus.AsyncPublishWithContext(ctx, events.NewDataCreatedEvent("test", 1, ""), nil)

	// Assert - wait with timeout
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// Success - handler was called with context
	case <-time.After(time.Second):
		suite.T().Fatal("Handler was not called within timeout")
	}
}

// TestAsyncPublish_UsesDefaultErrorCallback tests the default error callback when none is provided
func (suite *EventBusAsyncSuite) TestAsyncPublish_UsesDefaultErrorCallback() {
	// Arrange
	expectedErr := errors.New("handler error")

	handler := func(ctx context.Context, event events.Event) error {
		return expectedErr
	}

	// This test is a bit tricky as it relies on log output which we're not capturing here
	// In a real implementation, you might want to inject a mock logger or capture output

	// Act
	suite.bus.Subscribe(events.DataCreated, handler)
	suite.bus.AsyncPublish(context.Background(), events.NewDataCreatedEvent("test", 1, ""))

	// We need a small delay to allow async processing to complete
	time.Sleep(100 * time.Millisecond)

	// Assert - we can't directly assert on log output here, but we've at least verified
	// that the code doesn't panic when using the default error callback
	// A more thorough test would capture log output and verify it
}

// TestEventBusSuites runs all the test suites
func TestEventBusSuites(t *testing.T) {
	suite.Run(t, new(EventBusHappyPathSuite))
	suite.Run(t, new(EventBusErrorHandlingSuite))
	suite.Run(t, new(EventBusAsyncSuite))
}
