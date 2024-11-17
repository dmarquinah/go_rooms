package types

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dmarquinah/go_rooms/crypto"
	"github.com/dmarquinah/go_rooms/db"
	"github.com/dmarquinah/go_rooms/utils"
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

func GetHost(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if id, ok := r.Context().Value(utils.IdKey).(string); ok {
		hostRecord, err := findHostById(id, database)
		if err != nil {
			WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if hostRecord == nil {
			WriteErrorResponse(w, "Error retrieving host data.", http.StatusInternalServerError)
			return
		}

		if !hostRecord.IsVerified {
			WriteErrorResponse(w, "Unverified host.", http.StatusForbidden)
			return
		}

		WriteSuccessResponse(w, GetSuccessMessage(r), hostRecord)
		return
	}
}

func findHostById(id string, database *sql.DB) (*Host, error) {
	var host Host
	var hostDescription sql.NullString
	row := database.QueryRow(db.GET_LOGGED_HOST_STATEMENT, id)
	if err := row.Scan(&host.HostId, &host.HostUsername, &host.IsVerified, &host.CreatedAt, &hostDescription); err != nil {
		return nil, err
	}

	if hostDescription.Valid {
		host.Description = hostDescription.String
	}

	return &host, nil
}

func findHostByUsername(username string, database *sql.DB) (*Host, error) {
	var host Host
	row := database.QueryRow(db.GET_HOST_LOGIN_STATEMENT, username)
	if err := row.Scan(&host.HostId, &host.HostUsername, &host.Password, &host.IsVerified, &host.CreatedAt); err != nil {
		return nil, err
	}

	return &host, nil
}

func insertHost(host *Host, database *sql.DB) (*Host, error) {
	password, err := crypto.HashPassword(host.Password)

	if err != nil {
		return nil, err
	}

	insertResult, err := database.Exec(db.INSERT_HOST_STATEMENT, host.HostUsername, password)

	if err != nil {
		return nil, err
	}

	id, err := insertResult.LastInsertId()

	if err != nil {
		return nil, err
	}

	return &Host{HostId: int(id)}, nil

}
