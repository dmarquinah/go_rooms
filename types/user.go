package types

import (
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	UserId     int       `json:"user_id" bson:"user_id"`
	Email      string    `json:"email" bson:"email"`
	Password   string    `json:"password,omitempty" bson:"password"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UserHandle string    `json:"user_handle" bson:"user_handle"`
}

func BodyToUser(body []byte) (*User, error) {
	if len(body) == 0 {
		return nil, errors.New("empty request body")
	}
	var user User
	err := json.Unmarshal(body, &user)
	if err != nil {
		return nil, errors.New("error parsing body")
	}

	return &user, nil
}
