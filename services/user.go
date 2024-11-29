package services

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dmarquinah/go_rooms/crypto"
	"github.com/dmarquinah/go_rooms/database/db"
	"github.com/dmarquinah/go_rooms/types"
	"github.com/dmarquinah/go_rooms/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if id, ok := r.Context().Value(utils.IdKey).(string); ok {
		userRecord, err := findUserById(id, database)
		if err != nil {
			types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userRecord == nil {
			types.WriteErrorResponse(w, "Error retrieving user data.", http.StatusInternalServerError)
			return
		}

		types.WriteSuccessResponse(w, types.GetSuccessMessage(r), userRecord)
		return
	}
}

func findUserById(id string, database *sql.DB) (*types.User, error) {
	var user types.User
	var userHandle sql.NullString

	row := database.QueryRow(db.GET_LOGGED_USER_STATEMENT, id)
	if err := row.Scan(&user.UserId, &user.Email, &user.CreatedAt, &userHandle); err != nil {
		return nil, err
	}

	if userHandle.Valid {
		user.UserHandle = userHandle.String
	}

	return &user, nil
}

func findUserByEmail(email string, database *sql.DB) (*types.User, error) {
	var user types.User
	row := database.QueryRow(db.GET_USER_LOGIN_STATEMENT, email)
	if err := row.Scan(&user.UserId, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func insertUser(user *types.User, database *sql.DB) (*types.User, error) {
	password, err := crypto.HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	insertResult, err := database.Exec(db.INSERT_USER_STATEMENT, user.Email, password, time.Now())

	if err != nil {
		return nil, err
	}

	id, err := insertResult.LastInsertId()

	if err != nil {
		return nil, err
	}

	return &types.User{UserId: int(id)}, nil

}
