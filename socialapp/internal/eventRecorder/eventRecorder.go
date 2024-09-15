package eventRecorder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/igomez10/microservices/socialapp/pkg/db"
)

type EventRecorder struct {
	DB db.Querier
}

func (e *EventRecorder) RecordEvent(ctx context.Context, dbConn *sql.Tx, rawEvent interface{}, id int64) error {
	jsonEvent, err := json.Marshal(rawEvent)
	if err != nil {
		return fmt.Errorf("Error marshalling event payload: %v", err)
	}

	createEventParams := db.CreateEventParams{
		EventType:     fmt.Sprintf("%T", rawEvent),
		Payload:       jsonEvent,
		AggregateType: fmt.Sprintf("%T", rawEvent),
		AggregateID:   id,
		Version:       1,
	}

	if err := e.DB.CreateEvent(ctx, dbConn, createEventParams); err != nil {
		return fmt.Errorf("Error creating event: %v", err)
	}
	return nil
}
