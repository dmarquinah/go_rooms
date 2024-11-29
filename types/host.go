package types

import (
	"encoding/json"
	"errors"
	"time"
)

type Host struct {
	HostId       int       `json:"host_id" bson:"host_id"`
	HostUsername string    `json:"host_username" bson:"host_username"`
	Password     string    `json:"host_password,omitempty" bson:"host_password"`
	IsVerified   bool      `json:"is_verified" bson:"is_verified"`
	Description  string    `json:"description" bson:"description"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

func BodyToHost(body []byte) (*Host, error) {
	if len(body) == 0 {
		return nil, errors.New("empty request body")
	}

	var host Host
	err := json.Unmarshal(body, &host)
	if err != nil {
		return nil, errors.New("error parsing body")
	}

	return &host, nil
}
