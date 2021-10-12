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
	r.HandleFunc("/rooms", roomService.CreateRoom).Methods("POST")
	r.HandleFunc("/rooms", roomService.GetRooms).Methods("GET")
	r.HandleFunc("/", api.Version).Methods("GET")
	r.HandleFunc("/rooms/{id}", roomService.DeleteRoom).Methods("DELETE")
	r.HandleFunc("/rooms/{id}", roomService.PatchRoom).Methods("PATCH")

	r.HandleFunc("/hi", roomService.SayHi).Methods("POST")

	log.Println("HTTP server is up and running...")

	//err = http.ListenAndServe("localhost:7000", r)
	err = http.ListenAndServe("localhost:7000", r)
	if err != nil {
		panic(err)
	}
}
