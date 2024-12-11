package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dmarquinah/go_rooms/middlewares"
	"github.com/dmarquinah/go_rooms/models"
	socketmanager "github.com/dmarquinah/go_rooms/pkg/socket_manager"
	"github.com/dmarquinah/go_rooms/services"
	"github.com/dmarquinah/go_rooms/utils"
	"github.com/gorilla/websocket"
)

type RoomSessionHandler struct {
	upgrader           websocket.Upgrader
	roomSessionService *services.RoomSessionService
}

func NewRoomSessionHandler(
	roomSessionService *services.RoomSessionService,
) *RoomSessionHandler {
	return &RoomSessionHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Implement your origin validation
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		roomSessionService: roomSessionService,
	}
}

func (h *RoomSessionHandler) registerRoomSessionRoutes(router *http.ServeMux) {
	router.Handle("/room/{id}/ws", middlewares.JWTmiddlewareWS(h.handleJoinRoomSession()))
}

func (h *RoomSessionHandler) handleJoinRoomSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userId string
		var userRole string
		var ok bool

		conn, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil {
			// Note: At this point we can't write to response writer
			// as the connection upgrade has already been attempted
			return
		}
		defer conn.Close()

		conn.SetCloseHandler(func(code int, text string) error {
			log.Printf("Connection closed with code %d: %s", code, text)
			message := websocket.FormatCloseMessage(code, "")
			return conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		})

		roomId := r.PathValue("id")
		roomPassword := r.URL.Query().Get("room_password")

		if roomId == "" || roomPassword == "" {
			socketmanager.WriteWSError(conn, "Room ID/Password are required", socketmanager.WSErrBadRequest)
			return
		}

		if _, err := strconv.Atoi(roomId); err != nil {
			socketmanager.WriteWSError(conn, "Invalid Room", socketmanager.WSErrBadRequest)
			return
		}

		if userId, ok = r.Context().Value(utils.IdKey).(string); !ok {
			socketmanager.WriteWSError(conn, "Error retrieving data", socketmanager.WSErrInternalError)
			return
		}

		if userRole, ok = r.Context().Value(utils.RoleKey).(string); !ok {
			socketmanager.WriteWSError(conn, "Error retrieving data", socketmanager.WSErrInternalError)
			return
		}

		// The roomSessionInstance will hold the in-memory information about the room in session
		roomSessionInstance, err := h.roomSessionService.GetOrCreateRoomSession(r.Context(), roomId, roomPassword)
		if err != nil {
			socketmanager.WriteWSError(conn, fmt.Sprintf("Error connecting to room session: %v", err), socketmanager.WSErrInternalError)
			return
		}

		// Creates a basic User model for creating a connection
		userInstance := models.NewUserInstance(userId, userRole)

		// Join into a room to be able to receive and send messages to/from that room
		err = h.roomSessionService.JoinUserToRoomSession(roomSessionInstance, userInstance, conn)
		if err != nil {
			socketmanager.WriteWSError(conn, "Error stablishing session", socketmanager.WSErrInternalError)
			return
		}
	}
}
