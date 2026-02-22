package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Hub maintains the set of active clients and broadcasts messages to rooms.
// It uses Redis Pub/Sub so that messages are distributed across multiple
// server instances (horizontal scaling).
type Hub struct {
	// Map of roomID -> set of clients in that room
	rooms map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// localDeliver delivers a message received from Redis to local clients.
	localDeliver chan *BroadcastMessage

	// Redis client for Pub/Sub
	rdb *redis.Client
	ctx context.Context

	// Active Redis subscriptions per room (so we subscribe only once)
	subscriptions map[string]context.CancelFunc

	mu sync.RWMutex
}

// BroadcastMessage wraps a message with its target room.
type BroadcastMessage struct {
	RoomID  string
	Message []byte
	Sender  *Client // sender to optionally exclude
}

// NewHub creates a new Hub instance with Redis Pub/Sub support.
func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		rooms:         make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		localDeliver:  make(chan *BroadcastMessage, 256),
		rdb:           rdb,
		ctx:           context.Background(),
		subscriptions: make(map[string]context.CancelFunc),
	}
}

// Run starts the hub's main loop. Should be called as a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			firstInRoom := h.rooms[client.RoomID] == nil || len(h.rooms[client.RoomID]) == 0
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			count := len(h.rooms[client.RoomID])
			h.mu.Unlock()

			// Subscribe to Redis channel when the first client joins a room
			if firstInRoom {
				h.subscribeRoom(client.RoomID)
			}

			log.Printf("🟢 WS: User %s joined room %s (%d online)", client.UserID, client.RoomID, count)

			// Broadcast system message: user joined
			h.broadcastSystemMessage(client.RoomID, client.UserID, client.UserName, "joined", count)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					remaining := len(clients)
					if remaining == 0 {
						delete(h.rooms, client.RoomID)
						// Unsubscribe from Redis channel when no more local clients
						if cancel, ok := h.subscriptions[client.RoomID]; ok {
							cancel()
							delete(h.subscriptions, client.RoomID)
							log.Printf("📡 Redis: Unsubscribed from channel %s", redisChannel(client.RoomID))
						}
					}
					// Broadcast system message: user left (only if room still has clients)
					if remaining > 0 {
						h.mu.Unlock()
						h.broadcastSystemMessage(client.RoomID, client.UserID, client.UserName, "left", remaining)
						log.Printf("🔴 WS: User %s left room %s", client.UserID, client.RoomID)
						continue
					}
				}
			}
			h.mu.Unlock()
			log.Printf("🔴 WS: User %s left room %s", client.UserID, client.RoomID)

		case msg := <-h.localDeliver:
			// Deliver a message (received from Redis) to local WebSocket clients
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

// BroadcastToRoom publishes a message to the Redis Pub/Sub channel for the
// given room. All server instances subscribed to that channel will receive
// the message and deliver it to their local WebSocket clients.
func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	channel := redisChannel(roomID)
	if err := h.rdb.Publish(h.ctx, channel, message).Err(); err != nil {
		log.Printf("❌ Redis Publish error (channel %s): %v", channel, err)
		// Fallback: deliver locally so the current instance still works
		h.localDeliver <- &BroadcastMessage{RoomID: roomID, Message: message}
	}
}

// subscribeRoom starts a goroutine that subscribes to a Redis Pub/Sub
// channel and forwards incoming messages to localDeliver.
func (h *Hub) subscribeRoom(roomID string) {
	channel := redisChannel(roomID)
	subCtx, cancel := context.WithCancel(h.ctx)

	h.mu.Lock()
	h.subscriptions[roomID] = cancel
	h.mu.Unlock()

	pubsub := h.rdb.Subscribe(subCtx, channel)

	go func() {
		defer pubsub.Close()
		log.Printf("📡 Redis: Subscribed to channel %s", channel)

		ch := pubsub.Channel()
		for {
			select {
			case <-subCtx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				h.localDeliver <- &BroadcastMessage{
					RoomID:  roomID,
					Message: []byte(msg.Payload),
				}
			}
		}
	}()
}

// SystemMessage is broadcast when a user joins or leaves a room.
type SystemMessage struct {
	Type        string `json:"type"`   // "system"
	Action      string `json:"action"` // "joined" or "left"
	SenderId    string `json:"sender_id"`
	SenderName  string `json:"sender_name"`
	RoomID      string `json:"room_id"`
	OnlineCount int    `json:"online_count"`
	CreatedAt   string `json:"created_at"`
}

// broadcastSystemMessage publishes a join/leave system message via Redis.
func (h *Hub) broadcastSystemMessage(roomID, userID, userName, action string, onlineCount int) {
	msg := SystemMessage{
		Type:        "system",
		Action:      action,
		SenderId:    userID,
		SenderName:  userName,
		RoomID:      roomID,
		OnlineCount: onlineCount,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	data, _ := json.Marshal(msg)
	h.BroadcastToRoom(roomID, data)
}

// redisChannel returns the Redis Pub/Sub channel name for a given room.
// Room IDs are already prefixed with "direct:" or "group:".
func redisChannel(roomID string) string {
	return "chat:" + roomID
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
