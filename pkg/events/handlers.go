package events

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"fmt"
	"log"

	"github.com/goccy/go-json"
)

func LoggingHandler() Handler {
	return func(ctx context.Context, event Event) error {
		log.Printf("Event: %s, Payload: %+v", event.Type, event.Payload)
		return nil
	}
}

// TODO: Review is the Payload makes sense for the uses cases of the app
type MaterializedQueryPayload struct {
	QueryName       string          `json:"query_name"`
	QuerySQL        string          `json:"query_definition"`
	QueryData       json.RawMessage `json:"data"`
	QueryVersion    int             `json:"version"`
	QueryErrorCount int             `json:"error_count"`
	QueryLastError  string          `json:"last_error"`
}

func CreateMaterializedQueryEvent(materializedQuery types.MaterializedQuery) Event {
	return Event{
		Type: MaterializedQueryCreated,
		Payload: MaterializedQueryPayload{
			QueryName:       materializedQuery.Name,
			QuerySQL:        materializedQuery.Definition,
			QueryData:       materializedQuery.Data,
			QueryVersion:    materializedQuery.Version,
			QueryErrorCount: materializedQuery.ErrorCount,
			QueryLastError:  materializedQuery.LastError,
		},
	}
}

func RefreshMaterializedQueryEvent(materializedQuery types.MaterializedQuery) Event {
	return Event{
		Type: MaterializedQueryRefreshRequested,
		Payload: MaterializedQueryPayload{
			QueryName:       materializedQuery.Name,
			QuerySQL:        materializedQuery.Definition,
			QueryData:       materializedQuery.Data,
			QueryVersion:    materializedQuery.Version,
			QueryErrorCount: materializedQuery.ErrorCount,
			QueryLastError:  materializedQuery.LastError,
		},
	}
}

func GetMaterializedQueryPayload(event Event) (MaterializedQueryPayload, error) {
	payload, ok := event.Payload.(MaterializedQueryPayload)
	if !ok {
		return MaterializedQueryPayload{}, fmt.Errorf("Invalid Payload type for event %s: expected MaterializedQueryPayload, got %T",
			event.Type, event.Payload)
	}
	return payload, nil
}

func GetDataChangePayload(event Event) (DataChangePayload, error) {
	payload, ok := event.Payload.(DataChangePayload)
	if !ok {
		return DataChangePayload{}, fmt.Errorf("invalid payload type for event %s: expected DataChangePayload, got %T",
			event.Type, event.Payload)
	}
	return payload, nil
}
