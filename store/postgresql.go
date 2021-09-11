package store

import (
	"context"
	"golang-api-good-practices/api"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

type RoomService struct {
	Db *DB
}

func NewDB() (*DB, error) {
	dbPool, err := pgxpool.Connect(context.Background(), "postgresql://localuser:the-secret@localhost:5432/hotel")
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
	err := db.Ping(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := db.Acquire(context.Background())
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer conn.Release()

	var rooms []api.Room
	//err = conn.QueryRow(context.Background(), `
	//	SELECT id, availability
	//	  FROM gapp.rooms;
	//`).Scan(&room.ID, &room.Availability)
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(room.ID)
	//fmt.Println(room.Availability)
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

	//var availability string
	//var created string
	//var modified string

	//var rooms []api.Room

	//for rows.Next() {
	//	err := rows.Scan(&id, &availability, &created, &modified)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	fmt.Println(id)
	//	fmt.Println(availability)
	//	fmt.Println(created)
	//	fmt.Println(modified)
	//	newRoom := api.Room{
	//		Availability: availability,
	//	}
	//	rooms = append(rooms, newRoom)
	//}
	//defer rows.Close()

	return rooms, nil
}

//func (db *DB) GetAllRooms(ctx context.Context) ([]api.Room, error) {
//	tx, err := db.BeginTx(ctx, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	stmt, err := tx.Prepare(`SELECT * FROM rooms`)
//	if err != nil {
//		return nil, err
//	}
//
//	results, err := stmt.Query()
//	if err != nil {
//		tx.Rollback()
//		return nil, err
//	}
//	defer results.Close()
//	for results.Next() {
//		var id int
//		fmt.Println(results.Scan(&id))
//	}
//
//	rooms := make([]api.Room, 0)
//	err = results.Scan(&rooms)
//	if err != nil {
//		return nil, err
//	}
//
//	return rooms, err
//}
