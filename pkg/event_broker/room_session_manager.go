package eventbroker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dmarquinah/go_rooms/models"
	socketmanager "github.com/dmarquinah/go_rooms/pkg/socket_manager"
	"github.com/dmarquinah/go_rooms/repositories"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// RoomSessionManager will manage the broker connection to handle events/messages
type RoomSessionManager struct {
	PubSubClient  *redis.Client // Broker for distributed communication
	RedisClient   *redis.Client // Data persistance service needed for scaling
	channelPrefix string        // Value to differentiate which broker channel is being used by the application
	NodeId        string
	RoomsMap      map[string]*models.RoomSession // In-memory map for collection of rooms
	Mutex         sync.Mutex                     // Protect in-memory mapper
}

func NewRoomSessionManager(redisClient *redis.Client, pubSubClient *redis.Client, channelPrefix string) *RoomSessionManager {
	return &RoomSessionManager{
		PubSubClient:  pubSubClient,
		RedisClient:   redisClient,
		channelPrefix: channelPrefix,
		NodeId:        generateNodeId(),
		RoomsMap:      make(map[string]*models.RoomSession),
	}
}

// generateNodeId creates a unique identifier for this node
func generateNodeId() string {
	return fmt.Sprintf("node-%s", uuid.New().String())
}

// GetRoomById returns an instance of the RoomSession model saved in memory or cache
// or returns a new instance if it doesn't exists
func (rm *RoomSessionManager) GetRoomSessionById(roomId string, roomPassword string, roomRepository *repositories.RoomRepository) (*models.RoomSession, error) {

	room, exists := rm.RoomsMap[roomId] // Checks memory first
	if !exists {
		// Check Redis for existing room
		roomKey := fmt.Sprintf("room:session:%s", roomId)
		roomData, err := rm.RedisClient.Get(context.Background(), roomKey).Bytes()
		if err == redis.Nil {
			// Room doesn't exist in cache

			// TODO: Query database
			roomInstance, err := roomRepository.FindRoomById(roomId)
			if err != nil {
				return nil, fmt.Errorf("error retrieving room information")
			}

			if roomInstance.RoomCode != roomPassword {
				return nil, fmt.Errorf("error validating room id/password")
			}

			// Creating value to save in cache and also returning it
			room = &models.RoomSession{
				RoomId:     roomId,
				PwdCode:    roomInstance.RoomCode,
				SessionMap: make(map[string]*models.UserSession),
			}

			// Marshaling to save in cache
			roomJSON, err := json.Marshal(room)
			if err != nil {
				return nil, fmt.Errorf("error marshaling data")
			}

			// Store in cache
			err = rm.RedisClient.Set(context.Background(), roomKey, roomJSON, 24*time.Hour).Err()
			if err != nil {
				return nil, fmt.Errorf("error storing room in cache")
			}
		} else if err != nil {
			return nil, fmt.Errorf("error checking room in cache")
		} else {
			// Room exists in Redis, unmarshal it
			room = &models.RoomSession{}
			if err := json.Unmarshal(roomData, room); err != nil {
				return nil, fmt.Errorf("error unmarshaling room")
			}

			if room.PwdCode != roomPassword {
				return nil, fmt.Errorf("error validating room id/password")
			}

			return room, nil
		}
		return models.NewRoomSession(roomId), nil
	}

	return room, nil
}

// Store the roomSession instance in memory if it doesn't exist yet
func (m *RoomSessionManager) StoreRoomSession(
	roomSession *models.RoomSession,
) *models.RoomSession {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var session *models.RoomSession
	var exists bool
	session, exists = m.RoomsMap[roomSession.RoomId]

	if !exists {
		// Store connection in the RoomSession mapper
		m.RoomsMap[roomSession.RoomId] = roomSession
		session = m.RoomsMap[roomSession.RoomId]
	}

	return session
}

// Remove the Room session instance from memory if error ocurred
func (m *RoomSessionManager) RemoveRoomSession(
	roomId string,
) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	_, exists := m.RoomsMap[roomId]

	if exists {
		// Remove roomSession from memory
		delete(m.RoomsMap, roomId)
	}

	// Then we need to send a message to broker informing the removal of room session

	return nil
}

// TODO: Save the User session information into the database
func (m *RoomSessionManager) SaveUserSession(roomSession *models.RoomSession, userSession *models.UserSession, roomRepository *repositories.RoomRepository) error {
	// TODO: Save into RoomParticipant
	// TODO: Keep returning nil if the session already exists in the database

	return nil
}

// Remove the Room session instance from memory if error ocurred
func (m *RoomSessionManager) RemoveUserFromRoomSession(
	roomSession *models.RoomSession,
	user *models.User,
) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	session, exists := m.RoomsMap[roomSession.RoomId]

	if exists {
		// Remove roomSession from memory
		delete(session.SessionMap, user.UserId)
	}

	// TODO: Then we need to send a message to broker informing the removal of room session

	return nil
}

func (m *RoomSessionManager) StoreUserSession(roomSession *models.RoomSession, user *models.User, conn *websocket.Conn) *models.UserSession {
	var session *models.UserSession
	var exists bool

	session, exists = roomSession.SessionMap[user.UserId]
	if !exists {
		session = models.NewUserSession(user.UserId, user.UserRole, roomSession.RoomId, conn)
	}

	// If new or existing instance, always update LastSeen
	session.UpdateLastSeen()

	return session
}

func (m *RoomSessionManager) getClusterChannelRoom(userSession *models.UserSession) string {
	return fmt.Sprintf("%s:session:%s", m.channelPrefix, userSession.GetSessionRoom())
}

// publishMessage sends a message to the message broker
func (m *RoomSessionManager) PublishMessage(
	userSession *models.UserSession,
	msg *socketmanager.WebSocketMessage,
) error {
	// Encode message
	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	// Publish to Redis channel
	channelRoom := m.getClusterChannelRoom(userSession)
	return m.PubSubClient.Publish(context.Background(), channelRoom, encodedMsg).Err()
}

func (m *RoomSessionManager) SubscribeToSessionEvents(userSession *models.UserSession) {
	channel := m.getClusterChannelRoom(userSession)
	pubsub := m.PubSubClient.Subscribe(context.Background(), channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		m.handleChannelMessage(msg.Payload)
	}
}

// handleChannelMessage processes incoming messages from message broker
func (m *RoomSessionManager) handleChannelMessage(payload string) {
	// Decode message
	msg, err := socketmanager.DecodeWebSocketMessage([]byte(payload))
	if err != nil {
		log.Printf("Error decoding message: %v", err)
		return
	}

	// Ignore messages from same node
	if msg.NodeId == m.NodeId {
		return
	}

	switch msg.Type {
	case socketmanager.MessageTypeUserConnected:
		m.handleUserConnected(msg)
	case socketmanager.MessageTypeUserDisconnected:
		m.handleUserDisconnected(msg)
	case socketmanager.MessageTypeRoomBroadcast:
		m.handleRoomNextQueueTrack(msg)
	}
}

// Helper methods for message handling
func (m *RoomSessionManager) handleUserConnected(msg *socketmanager.WebSocketMessage) {
	//m.handleRoomBroadcast(msg)
}

func (m *RoomSessionManager) handleUserDisconnected(msg *socketmanager.WebSocketMessage) {
	//m.handleRoomBroadcast(msg)
}

func (m *RoomSessionManager) handleRoomNextQueueTrack(msg *socketmanager.WebSocketMessage) {
	// Broadcast to local connections in the specified room
	/* m.Connections.Range(func(key, value any) bool {
		connectionKey := key.(string)
		conn := value.(*websocket.Conn)

		// Extract room and user from connection key
		roomID, userID, err := m.parseConnectionKey(connectionKey)
		if err != nil || roomID != msg.RoomId {
			return true
		}

		// Send message to connection
		err = conn.WriteMessage(websocket.TextMessage, msg.Payload)
		if err != nil {
			log.Printf("Error broadcasting to %s: %v", userID, err)
		}

		return true
	}) */
}

// startHeartbeat sends periodic heartbeats to cluster
func (m *RoomSessionManager) StartHeartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		heartbeatMsg := &socketmanager.WebSocketMessage{
			Type:      socketmanager.MessageTypeNodeHeartbeat,
			NodeId:    m.NodeId,
			Timestamp: time.Now(),
		}
		err := m.PublishMessage(nil, heartbeatMsg)
		if err != nil {
			return
		}
	}
}
