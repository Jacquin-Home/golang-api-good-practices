package store

import (
	"context"
	"fmt"
	"golang-api-good-practices/api"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

type RoomService struct {
	Db *DB
}

func NewDB() (*DB, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPasswd := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbDatabase := os.Getenv("POSTGRES_DB")
	dbPool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPasswd, dbHost, dbPort, dbDatabase))
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &DB{dbPool}, nil
}

// GetAllRooms list all available rooms
func (db *DB) GetAllRooms(ctx context.Context) ([]api.Room, error) {

	var rooms []api.Room
	rows, err := db.Query(ctx, `
		SELECT id, availability
		  FROM gapp.rooms;
	`)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var room api.Room
		err := rows.Scan(&room.ID, &room.Availability)
		if err != nil {
			log.Println(err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// SaveRoom received a Room object and store it on database
func (db *DB) SaveRoom(ctx context.Context, room api.Room) error {

	stmt := `INSERT INTO gapp.rooms (
				availability, created, modified
			 ) VALUES ($1, NOW(), NOW())
	`
	_, err := db.Exec(ctx, stmt, room.Availability)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteRoom deletes a room of the specified room id
func (db *DB) DeleteRoom(ctx context.Context, id string) error {

	stmt := `DELETE FROM gapp.rooms 
				   WHERE id = $1
	`
	_, err := db.Exec(ctx, stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
