package services

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/universeh2h/report/internal/repositories"
)

type TransactionService struct {
	repo       *repositories.TransactionRepository
	clients    map[Client]bool
	broadcast  chan interface{}
	register   chan Client
	unregister chan Client
	mu         sync.RWMutex
}

type Client struct {
	conn interface{}
	Send chan interface{}
	ID   string
}

type TransactionUpdate struct {
	Transactions []repositories.TransactionData `json:"transactions"`
	Timestamp    time.Time                      `json:"timestamp"`
}

func NewClient(conn *websocket.Conn, id string) Client {
	return Client{
		conn: conn,
		Send: make(chan interface{}, 256),
		ID:   id,
	}
}

func NewTransactionsService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{
		repo:       repo,
		clients:    make(map[Client]bool),
		broadcast:  make(chan interface{}, 256),
		register:   make(chan Client),
		unregister: make(chan Client),
	}
}

// Start menjalankan hub untuk broadcast
func (s *TransactionService) Start() {
	go func() {
		for {
			select {
			case client := <-s.register:
				s.mu.Lock()
				s.clients[client] = true
				s.mu.Unlock()

			case client := <-s.unregister:
				s.mu.Lock()
				if _, ok := s.clients[client]; ok {
					delete(s.clients, client)
					close(client.Send)
				}
				s.mu.Unlock()

			case message := <-s.broadcast:
				s.mu.RLock()
				for client := range s.clients {
					select {
					case client.Send <- message:
					default:
						// Client's send channel is full, skip
					}
				}
				s.mu.RUnlock()
			}
		}
	}()
}

// CheckTransactions mengambil data sekali saja
func (s *TransactionService) CheckTransactions(c context.Context) ([]repositories.TransactionData, error) {
	return s.repo.GetTransactions(c)
}

// CheckTransactionsRealTime memulai polling untuk update real-time
func (s *TransactionService) CheckTransactionsRealTime(ctx context.Context, client Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.unregister <- client
			return

		case <-ticker.C:
			transactions, err := s.repo.GetTransactions(ctx)
			if err != nil {
				// Kirim error ke client
				s.broadcast <- map[string]interface{}{
					"error": err.Error(),
				}
				continue
			}

			update := TransactionUpdate{
				Transactions: transactions,
				Timestamp:    time.Now(),
			}

			s.broadcast <- update
		}
	}
}

// RegisterClient menambahkan client baru
func (s *TransactionService) RegisterClient(client Client) {
	s.register <- client
}

// UnregisterClient menghapus client
func (s *TransactionService) UnregisterClient(client Client) {
	s.unregister <- client
}
