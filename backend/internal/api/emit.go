package api

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

func (s *Server) emit(userID uuid.UUID, action string, data any) {
	if s.sse == nil {
		return
	}
	var payload []byte
	var err error

	switch action {
	case "event/delete", "calendar/delete":
		if id, ok := data.(uuid.UUID); ok {
			payload, err = id.MarshalText()
		}
	default:
		payload, err = json.Marshal(data)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	s.sse.Publish(userID.String(), &sse.Event{
		Event: []byte(action),
		Data:  payload,
	})
}
