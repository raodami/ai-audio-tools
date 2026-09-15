package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// WSClient represents a WebSocket client connection
type WSClient struct {
	conn *websocket.Conn
	send chan []byte
	id   string
}

// WSServer manages WebSocket connections
type WSServer struct {
	clients    map[string]*WSClient
	mu         sync.RWMutex
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan []byte
	clientCount atomic.Int32
}

// Upgrader for WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// NewWSServer creates a new WebSocket server
func NewWSServer() *WSServer {
	return &WSServer{
		clients:    make(map[string]*WSClient),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		broadcast:  make(chan []byte, 256),
	}
}

// Run starts the WebSocket server
func (s *WSServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.id] = client
			s.mu.Unlock()
			s.clientCount.Add(1)
			log.Printf("Client connected. Total: %d", s.clientCount.Load())

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.id]; ok {
				delete(s.clients, client.id)
				close(client.send)
				s.clientCount.Add(-1)
			}
			s.mu.Unlock()
			log.Printf("Client disconnected. Total: %d", s.clientCount.Load())

		case message := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client.id)
					s.clientCount.Add(-1)
				}
			}
			s.mu.RUnlock()
		}
	}
}

// ServeHTTP handles WebSocket upgrades
func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &WSClient{
		conn: conn,
		send: make(chan []byte, 256),
		id:   time.Now().Format(time.RFC3339Nano),
	}

	s.register <- client

	go s.writePump(client)
	go s.readPump(client)
}

func (s *WSServer) writePump(client *WSClient) {
	defer func() {
		client.conn.Close()
	}()

	for message := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func (s *WSServer) readPump(client *WSClient) {
	defer func() {
		s.unregister <- client
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}

		// Handle incoming messages
		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err == nil {
			s.handleMessage(client, &wsMsg)
		}
	}
}

func (s *WSServer) handleMessage(client *WSClient, msg *WSMessage) {
	switch msg.Type {
	case "ping":
		s.SendToClient(client.id, &WSMessage{
			Type:      "pong",
			Timestamp: time.Now().UnixMilli(),
		})
	case "subscribe":
		// Client can subscribe to specific channels
		log.Printf("Client %s subscribed to %v", client.id, msg.Payload)
	}
}

// SendToClient sends a message to a specific client
func (s *WSServer) SendToClient(clientID string, msg *WSMessage) {
	s.mu.RLock()
	client, ok := s.clients[clientID]
	s.mu.RUnlock()

	if !ok {
		return
	}

	data, _ := json.Marshal(msg)
	select {
	case client.send <- data:
	default:
		log.Printf("Failed to send to client %s", clientID)
	}
}

// Broadcast sends a message to all connected clients
func (s *WSServer) Broadcast(msg *WSMessage) {
	data, _ := json.Marshal(msg)
	select {
	case s.broadcast <- data:
	default:
		log.Printf("Broadcast queue full")
	}
}

// GetClientCount returns the number of connected clients
func (s *WSServer) GetClientCount() int32 {
	return s.clientCount.Load()
}
