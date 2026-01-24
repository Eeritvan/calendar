package stream

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/r3labs/sse/v2"
)

type SSEHandler struct {
	SSEServer   *sse.Server
	UserClients map[uuid.UUID]map[string]struct{}
}

// TODO:
// - max connections per userId
func (h *SSEHandler) HandleSSE(c *echo.Context) error {
	userId, ok := c.Get("userId").(uuid.UUID)
	if !ok {
		return nil // TODO: error handling i guess
	}

	streamToken := c.QueryParam("stream")
	if streamToken == "" {
		return c.JSON(http.StatusBadRequest, "stream token missing")
	}

	if h.UserClients[userId] == nil {
		h.UserClients[userId] = make(map[string]struct{})
	}
	h.UserClients[userId][streamToken] = struct{}{}

	h.SSEServer.CreateStream(streamToken)

	go func() {
		<-c.Request().Context().Done()
		delete(h.UserClients[userId], streamToken)
		if len(h.UserClients[userId]) == 0 {
			delete(h.UserClients, userId)
		}
		h.SSEServer.RemoveStream(streamToken)
	}()

	h.SSEServer.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *SSEHandler) Emit(userID uuid.UUID, action string, data any) {
	if h.SSEServer == nil || h.UserClients == nil {
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

	streamTokens, ok := h.UserClients[userID]
	if !ok {
		return
	}
	for streamToken := range streamTokens {
		h.SSEServer.Publish(streamToken, &sse.Event{
			Event: []byte(action),
			Data:  payload,
		})
	}
}
