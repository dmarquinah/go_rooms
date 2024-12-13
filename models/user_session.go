package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	socketmanager "github.com/dmarquinah/go_rooms/pkg/socket_manager"
	"github.com/gorilla/websocket"
)

type UserSession struct {
	userId   string
	roomId   string
	userRole string
	conn     *websocket.Conn
	lastSeen time.Time
	mutex    sync.Mutex
	//active   chan bool
}

func NewUserSession(userId string, userRole string, roomId string, conn *websocket.Conn) *UserSession {
	return &UserSession{
		userId:   userId,
		roomId:   roomId,
		userRole: userRole,
		conn:     conn,
		lastSeen: time.Now(),
	}
}

func (u *UserSession) CreateSessionLoop(roomSession *RoomSession) {

	// Send connection confirmation message to all session users
	go u.broadcastNewConnectionMsg(roomSession)

	// Creating a loop to keep hearing for messages
	u.createReadLoop()

}

func (u *UserSession) broadcastNewConnectionMsg(roomSession *RoomSession) {
	u.broadcastMessageToRoom(roomSession, fmt.Sprintf("The user %s has just joined the Room #%s", u.userId, u.roomId))
}

func (u *UserSession) broadcastMessageToRoom(roomSession *RoomSession, message string) {
	currentSessions := roomSession.getSessions()

	for userId, session := range currentSessions {

		// Not necessary show message to self
		if userId == u.userId {
			continue
		}

		// Protecting websocket WriteMessage
		session.mutex.Lock()
		defer session.mutex.Unlock()

		err := session.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending connection message")
			continue
		}
	}
}

func (u *UserSession) createReadLoop() {

	for {
		messageType, msgBytes, err := u.conn.ReadMessage()

		if err != nil {
			break
		}

		message := string(msgBytes)

		if messageType == websocket.TextMessage {
			if err := u.HandleSocketTextMessage(message); err != nil {
				if err != websocket.ErrReadLimit {
					// If read limit error, keep loop alive
					continue
				}
				break // Close loop any other case
			}
		}
	}
}

func (u *UserSession) HandleSocketTextMessage(command string) error {
	// Protecting websocket WriteMessage
	u.mutex.Lock()
	defer u.mutex.Unlock()

	err := u.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("thanks for the message: %s", command)))
	if err != nil {
		if err != websocket.ErrReadLimit {
			return err
		}
		log.Printf("encountered error: %v", err)
	}

	switch command {
	case string(socketmanager.MessageTypeRoomNextTrack):
		return nil
	default:
		return nil
	}

	//return nil
}

func (u *UserSession) UpdateLastSeen() {
	u.lastSeen = time.Now()
}

func (u *UserSession) GetSessionRoom() string {
	return u.roomId
}

func (u *UserSession) GetSessionUser() string {
	return u.userId
}
