package event

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Event struct {
	Id         string      `json:"id"`
	RecordID   string      `json:"recordId"`
	RecordType string      `json:"recordType"`
	Record     interface{} `json:"record"`
	EventType  string      `json:"eventType"`
}

type EventRecord interface {
	RecordId() string
	RecordType() string
}

func NewEvent(record EventRecord, eventType string) Event {
	return Event{
		Id:         generateId(),
		RecordID:   record.RecordId(),
		RecordType: record.RecordType(),
		Record:     record,
		EventType:  eventType,
	}
}

func generateId() string {
	return uuid.New().String()
}

func (event Event) Raw() []byte {
	data, _ := json.Marshal(event)
	return data
}
