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

// this is temporary, and will be replaced when there is a DB
var room Room

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
}
type RoomService struct {
	Store Store
}

func (rs RoomService) SayHi2() {
	rooms, err := rs.Store.GetAllRooms()
	if err != nil {
		log.Println(err)
	}
	for _, item := range rooms {
		fmt.Println(item)
	}
}

func (rs RoomService) SayHi(w http.ResponseWriter, r *http.Request) {
	rooms, err := rs.Store.GetAllRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte("hi"))
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
	room = Room{
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

	//if id != room.ID.String() {
	//	var problem = Problem{
	//		Description: "Unable to delete, room not found!",
	//	}
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusNotFound)
	//	err := json.NewEncoder(w).Encode(problem)
	//	if err != nil {
	//		log.Println(err)
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	return
	//}
	//// reset Room
	//room = Room{}
	w.WriteHeader(http.StatusNoContent)
}

func PatchRoom(w http.ResponseWriter, r *http.Request) {

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

	w.Header().Set("Content-Type", "application/json")
	if patchRoom.ID == room.ID && patchRoom.Availability == room.Availability {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if patchRoom.ID != room.ID {
		problem := Problem{
			Description: "Room not found!",
		}
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(problem)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	room.Availability = patchRoom.Availability

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
