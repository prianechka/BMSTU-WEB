package models

import "src/objects"

type StudentFullInfo struct {
	Student objects.Student `json:"student"`
	Room    objects.Room    `json:"room"`
}

type ThingFullInfo struct {
	Thing objects.Thing `json:"thing"`
	Room  objects.Room  `json:"room"`
}
