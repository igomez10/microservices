package eventRecorder

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/jackc/pgx/v5"
)

type EventRecorder struct {
	DB db.Querier
}

func (e *EventRecorder) RecordEvent(ctx context.Context, tx pgx.Tx, rawEvent interface{}, id int64) error {
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

	if err := e.DB.CreateEvent(ctx, tx.Conn(), createEventParams); err != nil {
		return fmt.Errorf("Error creating event: %v", err)
	}
	return nil
}
