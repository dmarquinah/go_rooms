package types

import (
	"database/sql"
	"io"
	"net/http"
	"strconv"

	"github.com/dmarquinah/go_rooms/crypto"
)

func LoginUser(w http.ResponseWriter, r *http.Request, database *sql.DB) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the body of request.", http.StatusBadRequest)
	}
	defer r.Body.Close()

	user, err := BodyToUser(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if user == nil {
		http.Error(w, "Error wrapping the body to user.", http.StatusBadRequest)
		return false
	}

	// Get actual user data
	userRecord, err := findUserByEmail(user.Email, database)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if userRecord == nil {
		http.Error(w, "User email not found.", http.StatusBadRequest)
		return false
	}

	// Compare the hashed password with the database
	if !crypto.VerifyPassword(user.Password, userRecord.Password) {
		http.Error(w, "Email/Password incorrect", http.StatusUnauthorized)
		return false
	}

	tokenString := crypto.GenerateJWT(userRecord.UserId)

	resData := ResponseData{
		Message: *tokenString,
	}

	resJSON := GetResponseDataJson(resData)

	if resJSON == nil {
		http.Error(w, "Error parsing the response data to JSON. ", http.StatusInternalServerError)
		return false
	}

	w.Write(*resJSON)
	w.WriteHeader(http.StatusOK)

	return true
}

func RegisterUser(w http.ResponseWriter, r *http.Request, database *sql.DB) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the body of request.", http.StatusBadRequest)
	}
	defer r.Body.Close()

	user, err := BodyToUser(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if user == nil {
		http.Error(w, "Error wrapping the body to user.", http.StatusBadRequest)
		return false
	}

	// Get actual user data
	userRecord, err := findUserByEmail(user.Email, database)

	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if userRecord != nil {
		http.Error(w, "User with provided email already exists.", http.StatusBadRequest)
		return false
	}

	inserted, err := createUser(user, database)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	resData := ResponseData{
		Message: strconv.Itoa(inserted.UserId),
	}

	resJSON := GetResponseDataJson(resData)

	if resJSON == nil {
		http.Error(w, "Error parsing the response data to JSON. ", http.StatusInternalServerError)
		return false
	}

	res, err := w.Write(*resJSON)
	w.WriteHeader(http.StatusOK)

	if err != nil {
		http.Error(w, "Error writing response. ", http.StatusInternalServerError)
		return false
	}

	return res != 0
}
