package types

import (
	"encoding/json"
	"errors"
	"time"
)

type Room struct {
	RoomId       int       `json:"room_id" bson:"room_id"`
	UserOwner    int       `json:"user_owner,omitempty" bson:"user_owner"`
	AssignedHost int       `json:"assigned_host,omitempty" bson:"assigned_host"`
	RoomCode     string    `json:"room_code,omitempty" bson:"room_code"`
	StartDate    time.Time `json:"start_date,omitempty" bson:"start_date"`
	EndDate      time.Time `json:"end_date,omitempty" bson:"end_date"`
}

func BodyToRoom(body []byte) (*Room, error) {
	if len(body) == 0 {
		return nil, errors.New("empty request body")
	}
	var room Room
	err := json.Unmarshal(body, &room)
	if err != nil {
		return nil, errors.New("error parsing body")
	}

	return &room, nil
}
