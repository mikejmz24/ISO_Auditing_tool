package events

import "context"

type EventType string

const (
	DataCreated EventType = "data_created"
	DataUpdated EventType = "data_updated"
	DataDeleted EventType = "data_deleted"

	MaterializedQueryRefreshRequested EventType = "materialized_query_refresh_request"
	MaterializedQueryCreated          EventType = "materialized_query_created"
	MaterializedQueryUpdated          EventType = "materiazlied_query_updated"
)

type Event struct {
	Type    EventType
	Payload any
}

type Handler func(ctx context.Context, event Event) error

type Subscriber interface {
	HandleEvent(ctx context.Context, event Event) (any, error)
}

type DataChangePayload struct {
	EntityType    string
	EntityID      any
	ChangeType    string
	AffectedQuery string
}
