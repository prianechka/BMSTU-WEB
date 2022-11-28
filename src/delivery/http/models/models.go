package models

type ShortResponseMessage struct {
	Comment string `json:"comment"`
}

type AuthRequestMessage struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseWithJWTMessage struct {
	Token string `json:"token"`
}

type GetAllStudentsRequest struct {
	Token string `json:"token"`
}
