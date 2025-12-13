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

func (e *EventRecorder) RecordEvent(ctx context.Context, tx pgx.Tx, rawEvent interface{}, aggregateID int64) error {
	jsonEvent, err := json.Marshal(rawEvent)
	if err != nil {
		return fmt.Errorf("Error marshalling event payload: %v", err)
	}

	// Generate snowflake ID for the event
	eventID, err := e.SnowflakeGenerator.NextID()
	if err != nil {
		return fmt.Errorf("Error generating snowflake ID for event: %v", err)
	}

	createEventParams := db.CreateEventWithIDParams{
		ID:            eventID,
		EventType:     fmt.Sprintf("%T", rawEvent),
		Payload:       jsonEvent,
		AggregateType: fmt.Sprintf("%T", rawEvent),
		AggregateID:   aggregateID,
		Version:       1,
	}

	if err := e.DB.CreateEventWithID(ctx, tx.Conn(), createEventParams); err != nil {
		return fmt.Errorf("Error creating event: %v", err)
	}
	return nil
}
