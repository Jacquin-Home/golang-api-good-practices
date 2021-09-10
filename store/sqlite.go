package store

import (
	"database/sql"
	"fmt"
	"hotel-jp/api"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

type RoomService struct {
	Db *DB
}

func NewDB() (*DB, error) {
	//err := os.Remove("./store/local.db")
	//if err != nil {
	//	return nil, err
	//}

	db, err := sql.Open("sqlite3", "./store/local.db")
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func CreateRoomTable(db *DB) {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS rooms (
			id           INTEGER PRIMARY KEY,
			availability TEXT,
			created      TEXT,
			modified     TEXT
		)
	`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}

func (db *DB) GetAllRooms() ([]api.Room, error) {
	fmt.Println("hi123")
	rows, err := db.Query(`SELECT * FROM rooms;`)
	if err != nil {
		return nil, err
	}
	x, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range x {
		fmt.Println(item)
	}

	var id string
	var availability string
	var created string
	var modified string

	var rooms []api.Room

	for rows.Next() {
		err := rows.Scan(&id, &availability, &created, &modified)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(id)
		fmt.Println(availability)
		fmt.Println(created)
		fmt.Println(modified)
		newRoom := api.Room{
			Availability: availability,
		}
		rooms = append(rooms, newRoom)
	}
	defer rows.Close()

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
