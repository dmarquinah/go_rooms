package models

import (
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
}

func NewUserSession(userId string, roomId string, userRole string, conn *websocket.Conn) *UserSession {
	return &UserSession{
		userId:   userId,
		roomId:   roomId,
		userRole: userRole,
		conn:     conn,
		lastSeen: time.Now(),
	}
}

func (u *UserSession) CreateSessionLoop() {

	log.Println("in session loop!")

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

func (u *UserSession) UpdateLastSeen() {
	u.lastSeen = time.Now()
}

func (u *UserSession) GetSessionRoom() string {
	return u.roomId
}

func (u *UserSession) GetSessionUser() string {
	return u.userId
}
