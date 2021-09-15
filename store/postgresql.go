package store

import (
	"context"
	"fmt"
	"golang-api-good-practices/api"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

type RoomService struct {
	Db *DB
}

func NewDB() (*DB, error) {
	//dbPool, err := pgxpool.Connect(context.Background(), "postgresql://localuser:the-secret@localhost:5432/hotel")
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
	// todo: is it needed?
	// https://github.com/jackc/pgx/issues/685
	//defer dbPool.Close()

	return &DB{dbPool}, nil
}

func (db *DB) GetAllRooms() ([]api.Room, error) {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer conn.Release()

	var rooms []api.Room
	rows, err := conn.Query(context.Background(), `
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

func (db *DB) SaveRoom(room api.Room) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	stmt := `INSERT INTO gapp.rooms (
				availability, created, modified
			 ) VALUES ($1, $2, $3)
	`
	_, err = conn.Exec(context.Background(), stmt, room.Availability, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *DB) DeleteRoom(id string) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		log.Println()
		return err
	}

	stmt := `DELETE FROM gapp.rooms 
				   WHERE id = $1
	`
	_, err = conn.Exec(context.Background(), stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
