package types

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dmarquinah/go_rooms/crypto"
	"github.com/dmarquinah/go_rooms/db"
	"github.com/dmarquinah/go_rooms/utils"
)

type Room struct {
	RoomId       int       `json:"room_id" bson:"room_id"`
	UserOwner    int       `json:"user_owner,omitempty" bson:"user_owner"`
	AssignedHost int       `json:"assigned_host,omitempty" bson:"assigned_host"`
	RoomCode     string    `json:"room_code,omitempty" bson:"room_code"`
	StartDate    time.Time `json:"start_date,omitempty" bson:"start_date"`
	EndDate      time.Time `json:"end_date,omitempty" bson:"end_date"`
}

func BodyToRoom(body []byte) (*Room, error) {
	if len(body) == 0 {
		return nil, errors.New("empty request body")
	}
	var room Room
	err := json.Unmarshal(body, &room)
	if err != nil {
		return nil, errors.New("error parsing body")
	}

	return &room, nil
}

func GetRoomById(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	id := r.PathValue("id")

	if id == "" {
		WriteErrorResponse(w, "Invalid Room ID", http.StatusInternalServerError)
		return
	}

	roomRecord, err := findRoomById(id, database)

	if err != nil && err != sql.ErrNoRows {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if roomRecord == nil {
		WriteErrorResponse(w, "Room not found.", http.StatusNotFound)
		return
	}

	WriteSuccessResponse(w, GetSuccessMessage(r), roomRecord)
}

func findRoomById(id string, database *sql.DB) (*Room, error) {
	var room Room
	var AssignedHost sql.NullInt16

	row := database.QueryRow(db.GET_ROOM_ID_STATEMENT, id)
	if err := row.Scan(&room.RoomId, &room.UserOwner, &AssignedHost, &room.RoomCode, &room.StartDate, &room.EndDate); err != nil {
		return nil, err
	}

	if AssignedHost.Valid {
		room.AssignedHost = int(AssignedHost.Int16)
	}

	return &room, nil
}

func CreateRoom(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	var id string
	var role string
	var ok bool

	if role, ok = r.Context().Value(utils.RoleKey).(string); !ok {
		WriteErrorResponse(w, "Internal Error on context.", http.StatusInternalServerError)
		return
	}

	if role == utils.HOST_ROLE {
		WriteErrorResponse(w, "Host can't create rooms.", http.StatusUnauthorized)
		return
	}

	if id, ok = r.Context().Value(utils.IdKey).(string); !ok {
		WriteErrorResponse(w, "Internal Error on context.", http.StatusBadRequest)
		return
	}
	conv_id, err := strconv.Atoi(id)

	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, "Error reading the body of request.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	room, err := BodyToRoom(body)

	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if room == nil {
		WriteErrorResponse(w, "Internal Error on Body.", http.StatusInternalServerError)
		return
	}

	// Setting up user Id into model
	room.UserOwner = conv_id

	// A user can't create multiple rooms on the same date
	roomRecord, err := findRoomByUserDate(id, room.StartDate, database)
	if err != nil && err != sql.ErrNoRows {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(roomRecord)

	if roomRecord != nil && roomRecord.RoomId != 0 {
		WriteErrorResponse(w, "Room already created on selected date.", http.StatusBadRequest)
		return
	}

	inserted, err := insertRoom(room, database)

	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteSuccessResponse(w, GetSuccessMessage(r), inserted)

}

func findRoomByUserDate(userId string, startDate time.Time, database *sql.DB) (*Room, error) {
	var room Room

	row := database.QueryRow(db.GET_ROOM_USER_DATE_STATEMENT, userId, startDate)
	if err := row.Scan(&room.RoomId); err != nil {
		return nil, err
	}

	return &room, nil

}

func insertRoom(room *Room, database *sql.DB) (*Room, error) {
	// Creating Room private code
	generatedCode, err := crypto.GenerateRandomCode(6)
	if err != nil {
		return nil, err
	}

	insertResult, err := database.Exec(db.INSERT_ROOM_STATEMENT, room.UserOwner, generatedCode, room.StartDate, room.EndDate)

	if err != nil {
		return nil, err
	}

	id, err := insertResult.LastInsertId()

	if err != nil {
		return nil, err
	}

	return &Room{RoomId: int(id)}, nil

}
