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
	pgUser := os.Getenv("PG_USERNAME")
	pgPass := os.Getenv("PG_PASSWORD")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDb := os.Getenv("PG_DATABASE")
	psqlConStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", pgUser, pgPass, pgHost, pgPort, pgDb)
	dbPool, err := pgxpool.Connect(context.Background(), psqlConStr)
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

func (db *DB) PatchRoom(room api.Room) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	stmt := `UPDATE gapp.rooms 
				SET availability
			  WHERE id=$1
	`
	_, err = conn.Exec(context.Background(), stmt, room.ID)
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
