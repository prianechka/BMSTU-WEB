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

type StudentHistoryResponseMessage struct {
	RoomID int `json:"room-id"`
}

type ThingOwnerHistoryResponseMessage struct {
	OwnerStudentNumber string `json:"owner-student-number"`
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
	RoomID int `json:"roomID"`
}

type StudentThingsActsRequestMessage struct {
	StudentNumber string `json:"stud-number"`
	Status        string `json:"status"`
}

type ResponseWithJWTMessage struct {
	Token string `json:"token"`
}

type ThingFullInfoResponse struct {
	Thing objects.ThingResponseDTO `json:"thing"`
}

type StudentFullInfoResponse struct {
	Student objects.StudentResponseDTO `json:"student"`
}

func CreateThingFullInfoResponse(things []models.ThingFullInfo) []ThingFullInfoResponse {
	result := make([]ThingFullInfoResponse, objects.Empty)
	for _, thing := range things {
		result = append(result, ThingFullInfoResponse{
			Thing: objects.CreateThingResponse(thing.Thing),
		})
	}
	return result
}

func CreateThingInfoResponse(thing models.ThingFullInfo) ThingFullInfoResponse {
	return ThingFullInfoResponse{Thing: objects.CreateThingResponse(thing.Thing)}
}

func CreateStudentFullInfoResponse(student models.StudentFullInfo) StudentFullInfoResponse {
	return StudentFullInfoResponse{
		Student: objects.CreateStudentResponseSingle(student.Student),
	}
}
