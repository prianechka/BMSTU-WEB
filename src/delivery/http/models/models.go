package models

import "src/objects"

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

type SettleInRoomRequestMessage struct {
	RoomID int `json:"roomID"`
}

type ResponseWithJWTMessage struct {
	Token string `json:"token"`
}
