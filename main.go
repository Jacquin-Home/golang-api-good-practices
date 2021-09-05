package main

import (
	"hotel-jp/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/room", api.CreateRoom).Methods("POST")
	r.HandleFunc("/room", api.GetRoom).Methods("GET")
	r.HandleFunc("/", api.Version).Methods("GET")
	r.HandleFunc("/room/{id}", api.DeleteRoom).Methods("DELETE")
	r.HandleFunc("/room/{id}", api.PatchRoom).Methods("PATCH")
	log.Println("HTTP server is up and running...")
	http.ListenAndServe(":7000", r)
}
