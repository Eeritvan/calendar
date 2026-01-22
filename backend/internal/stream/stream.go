package stream

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/r3labs/sse/v2"
)

type SSEHandler struct {
	SSEServer *sse.Server
}

// TODO:
// - max 5? connections per userId
// - bug: disconnects all clients at once if one disconnects
func (h *SSEHandler) HandleSSE(c *echo.Context) error {
	userId, ok := c.Get("userId").(uuid.UUID)
	if !ok {
		return nil // TODO: error handling i guess
	}

	userIdStr := userId.String()

	fmt.Println("user connected", userIdStr)

	h.SSEServer.CreateStream(userIdStr)
	go func() {
		<-c.Request().Context().Done()
		fmt.Println("user disconnected", userIdStr)
		h.SSEServer.RemoveStream(userIdStr)
	}()

	h.SSEServer.ServeHTTP(c.Response(), c.Request())

	return nil
}
