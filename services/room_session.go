package services

import (
	"context"

	"github.com/dmarquinah/go_rooms/models"
	"github.com/dmarquinah/go_rooms/repositories"

	eventbroker "github.com/dmarquinah/go_rooms/pkg/event_broker"
	"github.com/gorilla/websocket"
)

type RoomSessionService struct {
	manager  *eventbroker.RoomSessionManager
	RoomRepo *repositories.RoomRepository
}

func NewRoomSessionService(
	roomSessionManager *eventbroker.RoomSessionManager,
	roomRepository *repositories.RoomRepository,
) *RoomSessionService {
	return &RoomSessionService{
		manager:  roomSessionManager,
		RoomRepo: roomRepository,
	}
}

// Creates or retrieves the roomSession stored, or save in memory if it doesn't exist yet
func (r *RoomSessionService) GetOrCreateRoomSession(context context.Context, roomId string, roomPassword string) (*models.RoomSession, error) {

	// Retrieving the Room session, checks memory, cache and database
	roomSessionInstance, err := r.manager.GetRoomSessionById(roomId, roomPassword, r.RoomRepo)
	if err != nil {
		return nil, err
	}
	// Store it in memory if it doesn't exist
	r.manager.StoreRoomSession(roomSessionInstance)

	// TODO: Save into the database new or existing roomSession information

	return roomSessionInstance, nil
}

// Subscribe user to any events in a specific room
func (r *RoomSessionService) JoinUserToRoomSession(roomSession *models.RoomSession, user *models.User, conn *websocket.Conn) error {
	// Channel to handle errors that can ocurr during concurrency
	errChannel := make(chan error)

	// Store user session in memory
	userSession := r.manager.StoreUserSession(roomSession, user, conn)

	// Creates a WebSocket Read Loop for Session
	go userSession.CreateSessionLoop()

	// Deferring user session removal in case this func throws an error
	defer func() {
		err := r.manager.RemoveUserFromRoomSession(roomSession, user) // When session ends abruptly we want to remove the connection from memory
		if err != nil {
			errChannel <- err
		}
		close(errChannel)
	}()

	if err := <-errChannel; err != nil {
		return err
	}

	// Store user session information into the database
	go func() {
		if err := r.manager.SaveUserSession(roomSession, userSession, r.RoomRepo); err != nil {
			errChannel <- err
		}
		close(errChannel)
	}()

	// TODO: Send join message to broker
	/* go func() {
		msg := &eventbroker.WebSocketMessage{
			Type:      eventbroker.MessageTypeUserConnected,
			RoomId:    roomSession.RoomId,
			NodeId:    r.manager.NodeId,
			UserId:    userSession.GetSessionUser(),
			Payload:   "Greetings to all in room",
			Timestamp: time.Now(),
		}

		log.Println("Publishing to broker...")
		if err := r.manager.PublishMessage(userSession, msg); err != nil {
			errChannel <- err
			close(errChannel)
		}

		log.Println("Published!")

	}()
	*/
	go r.manager.SubscribeToSessionEvents(userSession)

	// TODO: Create a heartbeat

	return nil
}
