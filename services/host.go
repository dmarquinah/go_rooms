package services

import (
	"database/sql"
	"net/http"

	"github.com/dmarquinah/go_rooms/crypto"
	"github.com/dmarquinah/go_rooms/database/db"
	"github.com/dmarquinah/go_rooms/types"
	"github.com/dmarquinah/go_rooms/utils"
)

func GetHost(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if id, ok := r.Context().Value(utils.IdKey).(string); ok {
		hostRecord, err := findHostById(id, database)
		if err != nil {
			types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if hostRecord == nil {
			types.WriteErrorResponse(w, "Error retrieving host data.", http.StatusInternalServerError)
			return
		}

		if !hostRecord.IsVerified {
			types.WriteErrorResponse(w, "Unverified host.", http.StatusForbidden)
			return
		}

		types.WriteSuccessResponse(w, types.GetSuccessMessage(r), hostRecord)
		return
	}
}

func findHostById(id string, database *sql.DB) (*types.Host, error) {
	var host types.Host
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

func findHostByUsername(username string, database *sql.DB) (*types.Host, error) {
	var host types.Host
	row := database.QueryRow(db.GET_HOST_LOGIN_STATEMENT, username)
	if err := row.Scan(&host.HostId, &host.HostUsername, &host.Password, &host.IsVerified, &host.CreatedAt); err != nil {
		return nil, err
	}

	return &host, nil
}

func insertHost(host *types.Host, database *sql.DB) (*types.Host, error) {
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

	return &types.Host{HostId: int(id)}, nil

}
