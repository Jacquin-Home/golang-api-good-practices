package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID `json:"id"`
	Availability string    `json:"availability"`
}

type SuccessResponse struct {
	SuccessResponse string `json:"success"`
}

type Problem struct {
	Description string `json:"description"`
}

type Availability int

const (
	free Availability = iota
	reserved
	inuse
)

var availability = [...]string{"free", "reserved", "inuse"}

func (av Availability) String() string {
	return availability[av]
}

func (r Room) IsValid() bool {
	for _, item := range availability {
		if item == r.Availability {
			return true
		}
	}
	return false
}

type Store interface {
	GetAllRooms() ([]Room, error)
	SaveRoom(room Room) error
	DeleteRoom(id string) error
	PatchRoom(room Room) error
}
type RoomService struct {
	Store
}

func (rs RoomService) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var receivedRoom Room
	err := json.NewDecoder(r.Body).Decode(&receivedRoom)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !receivedRoom.IsValid() {
		var problem = Problem{
			Description: "Invalid availability value",
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(problem)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// create new room
	room := Room{
		Availability: receivedRoom.Availability,
	}

	err = rs.Store.SaveRoom(room)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	success := SuccessResponse{
		SuccessResponse: "Room created!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(success)
}

func (rs RoomService) GetRooms(w http.ResponseWriter, r *http.Request) {

	var rooms []Room
	rooms, err := rs.Store.GetAllRooms()
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(rooms)

	w.Header().Set("Content-Type", "application/json")
	if len(rooms) == 0 {
		var problem = Problem{Description: "Room not found!"}
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(problem)
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = json.NewEncoder(w).Encode(rooms)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rs RoomService) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	err := rs.Store.DeleteRoom(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs RoomService) PatchRoom(w http.ResponseWriter, r *http.Request) {
	var patchRoom Room
	err := json.NewDecoder(r.Body).Decode(&patchRoom)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	uID, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
	}
	patchRoom.ID = uID

	err = rs.Store.PatchRoom(patchRoom)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	success := SuccessResponse{
		SuccessResponse: "Entity updated!",
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(success)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
