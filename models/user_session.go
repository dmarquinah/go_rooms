package models

import (
	"fmt"
	"log"
	"sync"
	"time"

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
	go u.broadcastMessageToRoom(roomSession)

	for {
		messageType, msg, err := u.conn.ReadMessage()

		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			log.Printf("message recieved from client: %s\n", string(msg))

			// Protecting websocket WriteMessage
			u.mutex.Lock()
			defer u.mutex.Unlock()

			err := u.conn.WriteMessage(websocket.TextMessage, []byte("thanks for the data as message!"))
			if err != nil {
				if err != websocket.ErrReadLimit {
					continue
				}
				log.Printf("encountered error: %v", err)
			}
		}

	}
}

func (u *UserSession) broadcastMessageToRoom(roomSession *RoomSession) {
	currentSessions := roomSession.getSessions()

	for userId, session := range currentSessions {

		// Not necessary show message to self
		if userId == u.userId {
			continue
		}

		// Protecting websocket WriteMessage
		session.mutex.Lock()
		defer session.mutex.Unlock()

		err := session.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("The user %s has just joined the Room #%s", u.userId, u.roomId)))
		if err != nil {
			log.Println("Error sending connection message")
			continue
		}
	}
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
