package services

import (
	"database/sql"
	"io"
	"net/http"
	"strconv"

	"github.com/dmarquinah/go_rooms/crypto"
	"github.com/dmarquinah/go_rooms/types"
	"github.com/dmarquinah/go_rooms/utils"
)

func LoginUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		types.WriteErrorResponse(w, "Error reading the body of request.", http.StatusBadRequest)
	}
	defer r.Body.Close()

	user, err := types.BodyToUser(body)
	if err != nil {
		types.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user == nil {
		types.WriteErrorResponse(w, "Internal Error on Body.", http.StatusBadRequest)
		return
	}

	// Get actual user data
	userRecord, err := findUserByEmail(user.Email, database)

	if err != nil && err != sql.ErrNoRows {
		types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userRecord == nil {
		types.WriteErrorResponse(w, "User email not found.", http.StatusNotFound)
		return
	}

	// Compare the hashed password with the database
	if !crypto.VerifyPassword(user.Password, userRecord.Password) {
		types.WriteErrorResponse(w, "Email/Password incorrect", http.StatusUnauthorized)
	}

	tokenString := crypto.GenerateJWT(strconv.Itoa(userRecord.UserId), utils.USER_ROLE)

	types.WriteSuccessResponse(w, types.GetSuccessMessage(r), *tokenString)
}

func RegisterUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		types.WriteErrorResponse(w, "Error reading the body of request.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := types.BodyToUser(body)

	if err != nil {
		types.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user == nil {
		types.WriteErrorResponse(w, "Internal Error on Body.", http.StatusBadRequest)
		return
	}

	// Get actual user data
	userRecord, err := findUserByEmail(user.Email, database)

	if err != nil && err != sql.ErrNoRows {
		types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	if userRecord != nil {
		types.WriteErrorResponse(w, "User with provided email already exists.", http.StatusBadRequest)
	}

	inserted, err := insertUser(user, database)

	if err != nil {
		types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	types.WriteSuccessResponse(w, types.GetSuccessMessage(r), inserted)

}

func LoginHost(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		types.WriteErrorResponse(w, "Error reading the body of request.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	host, err := types.BodyToHost(body)
	if err != nil {
		types.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if host == nil {
		types.WriteErrorResponse(w, "Internal Error on Body.", http.StatusBadRequest)
		return
	}

	// Get actual user data
	hostRecord, err := findHostByUsername(host.HostUsername, database)

	if err != nil && err != sql.ErrNoRows {
		types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if hostRecord == nil {
		types.WriteErrorResponse(w, "Host username not found.", http.StatusBadRequest)
		return
	}

	if !hostRecord.IsVerified {
		types.WriteErrorResponse(w, "Host is not verified yet.", http.StatusBadRequest)
		return
	}

	// Compare the hashed password with the database
	if !crypto.VerifyPassword(host.Password, hostRecord.Password) {
		types.WriteErrorResponse(w, "Username/Password incorrect", http.StatusUnauthorized)
		return
	}

	tokenString := crypto.GenerateJWT(strconv.Itoa(hostRecord.HostId), utils.HOST_ROLE)

	types.WriteSuccessResponse(w, types.GetSuccessMessage(r), *tokenString)
}

func RegisterHost(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		types.WriteErrorResponse(w, "Error reading the body of request.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	host, err := types.BodyToHost(body)

	if err != nil {
		types.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if host == nil {
		types.WriteErrorResponse(w, "Internal Error on Body.", http.StatusBadRequest)
		return
	}

	// Get actual user data
	hostRecord, err := findHostByUsername(host.HostUsername, database)

	if err != nil && err != sql.ErrNoRows {
		types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if hostRecord != nil {
		types.WriteErrorResponse(w, "Host with provided username already exists.", http.StatusBadRequest)
		return
	}

	inserted, err := insertHost(host, database)

	if err != nil {
		types.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	types.WriteSuccessResponse(w, types.GetSuccessMessage(r), inserted)
}
