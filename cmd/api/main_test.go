package main

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketConnection(t *testing.T) {
	// Connect ke WebSocket
	u := url.URL{Scheme: "ws", Host: "localhost:1000", Path: "/ws/transactions"}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer ws.Close()

	// Send start command
	msg := map[string]interface{}{
		"command":  "start",
		"interval": 5,
	}
	err = ws.WriteJSON(msg)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Listen untuk messages
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			var response interface{}
			err := ws.ReadJSON(&response)
			if err != nil {
				fmt.Printf("Error reading: %v\n", err)
				return
			}
			fmt.Printf("Received: %v\n", response)
		}
	}()

	// Wait 30 seconds then close
	select {
	case <-time.After(30 * time.Second):
	case <-done:
	}
}
