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
	// Channel to handle errors that concurrent calls might have
	errChannel := make(chan error)

	// Store user session in memory
	userSession := r.manager.StoreUserSession(roomSession, user, conn)

	// Creates a WebSocket Read Loop for Session
	go userSession.CreateSessionLoop(roomSession)

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

	// TODO: Create a pingpong

	return nil
}
