package models

type ShortResponseMessage struct {
	Comment string `json:"comment"`
}

type AuthRequestMessage struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AddNewThingRequestMessage struct {
	MarkNumber string `json:"markNumber"`
	ThingType  string `json:"thingType"`
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
