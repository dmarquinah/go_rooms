package socketmanager

import (
	"encoding/json"
	"time"
)

// WebSocketMessage represents a message for distributed WebSocket communication
type WebSocketMessage struct {
	Type      MessageType `json:"type"`
	RoomId    string      `json:"room_id"`
	UserId    string      `json:"user_id"`
	Payload   string      `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	NodeId    string      `json:"node_id"`
}

// MessageType defines different types of distributed messages
type MessageType string

const (
	// User connection events
	MessageTypeUserConnected    MessageType = "user_connected"
	MessageTypeUserDisconnected MessageType = "user_disconnected"

	// Room-specific messages
	MessageTypeRoomBroadcast MessageType = "room_queue_next_track"

	// Control messages
	MessageTypeNodeHeartbeat MessageType = "node_heartbeat"
)

// Encoder for WebSocket messages
func (m *WebSocketMessage) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decoder for WebSocket messages
func DecodeWebSocketMessage(data []byte) (*WebSocketMessage, error) {
	var msg WebSocketMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
