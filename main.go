package main

import (
	"hotel-jp/api"
	"hotel-jp/store"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	db, err := store.NewDB()
	if err != nil {
		panic(err)
	}

	//store.CreateRoomTable(db)

	service := api.Service{
		Store: db,
	}
	//service.SayHi2()

	r := mux.NewRouter()
	r.HandleFunc("/room", api.CreateRoom).Methods("POST")
	r.HandleFunc("/room", api.GetRoom).Methods("GET")
	r.HandleFunc("/", api.Version).Methods("GET")
	r.HandleFunc("/room/{id}", api.DeleteRoom).Methods("DELETE")
	r.HandleFunc("/room/{id}", api.PatchRoom).Methods("PATCH")

	r.HandleFunc("/hi", service.SayHi).Methods("POST")

	log.Println("HTTP server is up and running...")

	err = http.ListenAndServe(":7000", r)
	if err != nil {
		panic(err)
	}
}
