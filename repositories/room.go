package repositories

import (
	"database/sql"

	"github.com/dmarquinah/go_rooms/database/db"
	"github.com/dmarquinah/go_rooms/types"
)

type RoomRepository struct {
	dbConn *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{dbConn: db}
}

func (r *RoomRepository) FindRoomById(id string) (*types.Room, error) {
	var room types.Room
	var AssignedHost sql.NullInt16

	row := r.dbConn.QueryRow(db.GET_ROOM_ID_STATEMENT, id)
	if err := row.Scan(&room.RoomId, &room.UserOwner, &AssignedHost, &room.RoomCode, &room.StartDate, &room.EndDate); err != nil {
		return nil, err
	}

	if AssignedHost.Valid {
		room.AssignedHost = int(AssignedHost.Int16)
	}

	return &room, nil
}
