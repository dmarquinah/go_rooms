package types

import "encoding/json"

type Host struct {
	HostId       int    `json:"host_id" bson:"host_id"`
	HostUsername string `json:"host_username" bson:"host_username"`
	HostPassword string `json:"host_password" bson:"host_password"`
}

func BodyToHost(body []byte) *Host {
	if len(body) == 0 {
		return nil
	}

	var host Host
	err := json.Unmarshal(body, &host)
	if err != nil {
		return nil
	}

	return &host
}
