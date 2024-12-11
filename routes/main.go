package routes

import (
	"database/sql"
	"net/http"

	eventbroker "github.com/dmarquinah/go_rooms/pkg/event_broker"
	"github.com/dmarquinah/go_rooms/repositories"
	"github.com/dmarquinah/go_rooms/services"
	"github.com/redis/go-redis/v9"
)

func BuildRouter(database *sql.DB, redisClient *redis.Client, pubSubClient *redis.Client) *http.ServeMux {
	router := http.NewServeMux()

	// Basic API Routes
	createAuthRoutes(router, database)
	createUserRoutes(router, database)
	createHostRoutes(router, database)
	createRoomRoutes(router, database)

	// WS Routes
	createRoomSessionRoutes(router, database, redisClient, pubSubClient)

	return router
}

func createRoomSessionRoutes(router *http.ServeMux, database *sql.DB, redisClient *redis.Client, pubSubClient *redis.Client) {

	// Creating the Room Session Manager
	roomSessionManager := eventbroker.NewRoomSessionManager(redisClient, pubSubClient, "kroom")
	roomRepository := repositories.NewRoomRepository(database)
	wsService := services.NewRoomSessionService(roomSessionManager, roomRepository)

	// Create the WebSocket handler
	wsHandler := NewRoomSessionHandler(wsService)

	wsHandler.registerRoomSessionRoutes(router)
}
