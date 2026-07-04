package eventRecorder

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/snowflake"
	"github.com/jackc/pgx/v5"
)

type EventRecorder struct {
	DB                 db.Querier
	SnowflakeGenerator snowflake.IDGenerator
}

// Event describes a domain event to be appended to the event store (a
// transactional outbox). It is recorded inside the caller's transaction so the
// event is committed atomically with the business change it describes.
type Event struct {
	// AggregateID identifies the domain entity the event applies to (e.g. a user's id).
	AggregateID int64
	// AggregateType is the kind of entity the event applies to (e.g. "User").
	AggregateType string
	// EventType is the name of the fact that occurred (e.g. "UserCreated").
	EventType string
	// Payload is the event body; it is serialized to JSON for storage.
	Payload any
}

// RecordEvent appends an event to the event store. The event's version is a
// per-aggregate, monotonically increasing sequence: each new event for the same
// aggregate (AggregateID + AggregateType) receives the next version. This
// preserves the order in which events were applied to a single entity and
// enables state replay and optimistic-concurrency checks.
func (e *EventRecorder) RecordEvent(ctx context.Context, tx pgx.Tx, event Event) error {
	jsonPayload, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("error marshalling event payload: %w", err)
	}

	// Generate snowflake ID for the event
	eventID, err := e.SnowflakeGenerator.NextID()
	if err != nil {
		return fmt.Errorf("error generating snowflake ID for event: %w", err)
	}

	// Read within the caller's transaction so the version is consistent with any
	// events this transaction may have already appended for the same aggregate.
	latestVersion, err := e.DB.GetLatestEventVersionForAggregate(ctx, tx.Conn(), db.GetLatestEventVersionForAggregateParams{
		AggregateID:   event.AggregateID,
		AggregateType: event.AggregateType,
	})
	if err != nil {
		return fmt.Errorf("error getting latest event version for aggregate %s/%d: %w", event.AggregateType, event.AggregateID, err)
	}

	createEventParams := db.CreateEventWithIDParams{
		ID:            eventID,
		AggregateID:   event.AggregateID,
		AggregateType: event.AggregateType,
		EventType:     event.EventType,
		Payload:       jsonPayload,
		Version:       latestVersion + 1,
	}

	if err := e.DB.CreateEventWithID(ctx, tx.Conn(), createEventParams); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	return nil
}
