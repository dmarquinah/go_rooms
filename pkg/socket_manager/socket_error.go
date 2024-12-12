package socketmanager

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// ErrorResponse maintains the same structure as your HTTP errors
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// WSError represents custom error codes for WebSocket
const (
	WSErrBadRequest    = 4000 // Similar to HTTP 400
	WSErrUnauthorized  = 4001 // Similar to HTTP 401
	WSErrForbidden     = 4003 // Similar to HTTP 403
	WSErrNotFound      = 4004 // Similar to HTTP 404
	WSErrInternalError = 4500 // Similar to HTTP 500
)

// WriteWSError sends an error message through the WebSocket connection
func WriteWSError(conn *websocket.Conn, message string, code int) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}

	// Marshal the error response
	payload, err := json.Marshal(response)
	if err != nil {
		// If marshaling fails, try to send a basic error
		basicError := []byte(`{"status":"error","message":"Internal server error","code":4500}`)
		closeErr := conn.WriteMessage(websocket.TextMessage, basicError)
		if closeErr != nil {
			return
		}
		return
	}

	// Write the error message as a text message
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return
	}

	// Optional: Close the connection with an appropriate code
	// You might want to only close on certain error types
	if code >= 4000 {
		conn.Close()
	}

}
