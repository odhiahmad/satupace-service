package websocket

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to rooms.
type Hub struct {
	// Map of roomID -> set of clients in that room
	rooms map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Inbound messages from clients to broadcast
	broadcast chan *BroadcastMessage

	mu sync.RWMutex
}

// BroadcastMessage wraps a message with its target room.
type BroadcastMessage struct {
	RoomID  string
	Message []byte
	Sender  *Client // sender to optionally exclude
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Run starts the hub's main loop. Should be called as a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			count := len(h.rooms[client.RoomID])
			h.mu.Unlock()
			log.Printf("🟢 WS: User %s joined room %s (%d online)", client.UserID, client.RoomID, count)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("🔴 WS: User %s left room %s", client.UserID, client.RoomID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.rooms[msg.RoomID]
			h.mu.RUnlock()

			for client := range clients {
				select {
				case client.Send <- msg.Message:
				default:
					// Client buffer full, disconnect
					h.mu.Lock()
					delete(h.rooms[msg.RoomID], client)
					close(client.Send)
					h.mu.Unlock()
				}
			}
		}
	}
}

// BroadcastToRoom sends a message to all clients in a room.
func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	h.broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: message,
	}
}

// GetOnlineCount returns the number of online users in a room.
func (h *Hub) GetOnlineCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}

// Register sends a client to the register channel.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister sends a client to the unregister channel.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
