package main

import (
	"golang-api-good-practices/api"
	"golang-api-good-practices/store"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	db, err := store.NewDB()
	if err != nil {
		panic(err)
	}
	// todo: is it needed?
	// https://github.com/jackc/pgx/issues/685
	defer db.Close()

	roomService := api.RoomService{
		Store: db,
	}

	r := mux.NewRouter()
	r.HandleFunc("/rooms", roomService.CreateRoom).Methods(http.MethodPost)
	r.HandleFunc("/rooms", roomService.GetRooms).Methods(http.MethodGet)
	r.HandleFunc("/", api.Version).Methods(http.MethodGet)
	r.HandleFunc("/rooms/{id}", roomService.DeleteRoom).Methods(http.MethodDelete)
	r.HandleFunc("/rooms/{id}", roomService.PatchRoom).Methods(http.MethodPatch)

	log.Println("HTTP server is up and running...")

	//err = http.ListenAndServe("localhost:7000", r)
	err = http.ListenAndServe("localhost:7000", r)
	if err != nil {
		panic(err)
	}
}
