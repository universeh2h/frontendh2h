package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/universeh2h/report/internal/services"
)

type TransactionHandler struct {
	service *services.TransactionService
	upgrade websocket.Upgrader
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		upgrade: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// CheckTransactionsRealTime handle WebSocket connection
func (h *TransactionHandler) CheckTransactionsRealTime(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrade.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	client := services.NewClient(conn, r.Header.Get("Sec-WebSocket-Key"))
	h.service.RegisterClient(client)

	// Start monitoring goroutinh.service.RegisterClient(client)
	go h.handleClientMessages(conn, client)
	go h.handleClientSend(conn, client)
}

// handleClientMessages baca pesan dari client
func (h *TransactionHandler) handleClientMessages(conn *websocket.Conn, client services.Client) {
	defer func() {
		h.service.UnregisterClient(client)
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Log error
			}
			break
		}

		// Handle client messages jika diperlukan
		if cmd, ok := msg["command"].(string); ok && cmd == "start" {
			interval := 5 * time.Second // Default interval
			if v, ok := msg["interval"].(float64); ok {
				interval = time.Duration(v) * time.Second
			}

			go h.service.CheckTransactionsRealTime(
				context.Background(),
				client,
				interval,
			)
		}
	}
}

// handleClientSend kirim pesan ke client
func (h *TransactionHandler) handleClientSend(conn *websocket.Conn, client services.Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
