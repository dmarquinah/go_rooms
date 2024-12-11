package models

import "sync"

type RoomSession struct {
	RoomId     string
	PwdCode    string
	SessionMap map[string]*UserSession // Mapping active user sessions
	Mutex      sync.Mutex              // Protects Users map
}

func NewRoomSession(roomId string) *RoomSession {
	return &RoomSession{
		RoomId:     roomId,
		SessionMap: make(map[string]*UserSession),
	}
}
