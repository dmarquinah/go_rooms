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

type User struct {
	UserId     int       `json:"user_id" bson:"user_id"`
	Email      string    `json:"email" bson:"email"`
	Password   string    `json:"password" bson:"password"`
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

func GetUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if id, ok := r.Context().Value(utils.IdKey).(string); ok {
		userRecord, err := getUserFromId(id, database)
		if err != nil {
			WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)

		}

		if userRecord == nil {
			WriteErrorResponse(w, "Error retrieving user data.", http.StatusInternalServerError)
		}

		WriteSuccessResponse(w, GetSuccessMessage(r), userRecord)

	}
}

func getUserFromId(id string, database *sql.DB) (*User, error) {
	rows, err := database.Query(db.GET_LOGGED_USER_STATEMENT, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterating rows
	rows.Next()
	var user User
	err = rows.Scan(&user.UserId, &user.Email, &user.CreatedAt, &user.UserHandle)
	if err != nil {
		return nil, err
	}

	// Check errors on rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func findUserByEmail(email string, database *sql.DB) (*User, error) {
	var user User
	row := database.QueryRow(db.GET_USER_LOGIN_STATEMENT, email)
	if err := row.Scan(&user.UserId, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func insertUser(user *User, database *sql.DB) (*User, error) {
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

	return &User{UserId: int(id)}, nil

}
