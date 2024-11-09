package types

import (
	"database/sql"
	"io"
	"net/http"
	"strconv"

	"github.com/dmarquinah/go_rooms/crypto"
)

func LoginUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, "Error reading the body of request.", http.StatusBadRequest)
	}
	defer r.Body.Close()

	user, err := BodyToUser(body)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user == nil {
		WriteErrorResponse(w, "Error wrapping the body to user.", http.StatusBadRequest)
		return
	}

	// Get actual user data
	userRecord, err := findUserByEmail(user.Email, database)

	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userRecord == nil {
		WriteErrorResponse(w, "User email not found.", http.StatusBadRequest)
		return
	}

	// Compare the hashed password with the database
	if !crypto.VerifyPassword(user.Password, userRecord.Password) {
		WriteErrorResponse(w, "Email/Password incorrect", http.StatusUnauthorized)
	}

	tokenString := crypto.GenerateJWT(userRecord.UserId)

	WriteSuccessResponse(w, GetSuccessMessage(r), *tokenString)
}

func RegisterUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, "Error reading the body of request.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := BodyToUser(body)

	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user == nil {
		WriteErrorResponse(w, "Error wrapping the body to user.", http.StatusBadRequest)
		return
	}

	// Get actual user data
	userRecord, err := findUserByEmail(user.Email, database)

	if err != nil && err != sql.ErrNoRows {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	if userRecord != nil {
		WriteErrorResponse(w, "User with provided email already exists.", http.StatusBadRequest)
	}

	inserted, err := insertUser(user, database)

	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	WriteSuccessResponse(w, GetSuccessMessage(r), strconv.Itoa(inserted.UserId))

}