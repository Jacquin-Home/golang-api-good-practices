package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
)

type Room struct {
	Id           uuid.UUID `json:"id"`
	Availability string    `json:"availability"`
}

type SuccessResponse struct {
	SuccessResponse string `json:"success"`
}

type Problem struct {
	Description string `json:"description"`
}

var room Room

func CreateRoom(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var receivedRoom Room
	err = json.Unmarshal(body, &receivedRoom)
	if err != nil {
		log.Println(err)
	}
	roomId := uuid.New()
	room = Room{
		Id:           roomId,
		Availability: receivedRoom.Availability,
	}
	var success SuccessResponse
	success = SuccessResponse{
		SuccessResponse: "Room created!",
	}
	jsonData, err := json.Marshal(success)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if room.Id.ID() != 0 {
		jsonData, err := json.Marshal(room)
		if err != nil {
			log.Println(err)
		}
		w.Write(jsonData)
		return
	}
	var problem = Problem{Description: "Room not found!"}
	jsonErr, err := json.Marshal(problem)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(404)
	w.Write(jsonErr)

}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == room.Id.String() {
		// reset Room
		room = Room{}
		w.WriteHeader(204)
		return
	}
	var problem Problem
	problem = Problem{
		Description: "Unable to delete, room not found!",
	}

	jsonData, err := json.Marshal(problem)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	w.Write(jsonData)
}

func PatchRoom(w http.ResponseWriter, r *http.Request) {

	var patchRoom Room

	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(resp, &patchRoom)
	if err != nil {
		log.Println(err)
	}

	params := mux.Vars(r)
	id := params["id"]
	uId, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
	}
	patchRoom.Id = uId

	w.Header().Set("Content-Type", "application/json")
	if patchRoom.Id == room.Id && patchRoom.Availability == room.Availability {
		w.WriteHeader(204)
		return
	} else if patchRoom.Id != room.Id {
		var problem = Problem{
			Description: "Room not found!",
		}
		errData, err := json.Marshal(problem)
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(404)
		w.Write(errData)
		return
	}
	room.Availability = patchRoom.Availability

	var success = SuccessResponse{
		SuccessResponse: "Entity updated!",
	}
	sData, err := json.Marshal(success)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(200)
	w.Write(sData)
}
