package models

import (
	"src/logic/managers/models"
	"src/objects"
)

type ShortResponseMessage struct {
	Comment string `json:"comment"`
}

type StudentResponseMessage struct {
	Students []objects.Student `json:"students"`
}

type AuthRequestMessage struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AddNewThingRequestMessage struct {
	MarkNumber int    `json:"markNumber"`
	ThingType  string `json:"thingType"`
}

type TransferThingRequestMessage struct {
	NewRoomID int `json:"room-id"`
}

type AddNewStudentRequestMessage struct {
	Login         string `json:"login"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Group         string `json:"group"`
	StudentNumber string `json:"studentNumber"`
}

type ChangeStudentGroupRequestMessage struct {
	NewGroup string `json:"newGroup"`
}

type StudentLiveActsRequestMessage struct {
	StudentNumber string `json:"student-number"`
	RoomID        int    `json:"roomID"`
}

type StudentThingsActsRequestMessage struct {
	StudentNumber string `json:"student-number"`
	MarkNumber    int    `json:"mark-number"`
	Status        string `json:"status"`
}

type ResponseWithJWTMessage struct {
	Token string `json:"token"`
}

type ThingFullInfoResponse struct {
	Thing objects.ThingResponseDTO `json:"thing"`
	Room  objects.RoomResponseDTO  `json:"room"`
}

type StudentFullInfoResponse struct {
	Student objects.StudentResponseDTO `json:"student"`
	Room    objects.RoomResponseDTO    `json:"room"`
}

func CreateThingFullInfoResponse(things []models.ThingFullInfo) []ThingFullInfoResponse {
	result := make([]ThingFullInfoResponse, objects.Empty)
	for _, thing := range things {
		result = append(result, ThingFullInfoResponse{
			Thing: objects.CreateThingResponse(thing.Thing),
			Room:  objects.CreateRoomResponseDTO(thing.Room),
		})
	}
	return result
}

func CreateStudentFullInfoResponse(student models.StudentFullInfo) StudentFullInfoResponse {
	return StudentFullInfoResponse{
		Student: objects.CreateStudentResponseSingle(student.Student),
		Room:    objects.CreateRoomResponseDTO(student.Room),
	}
}
