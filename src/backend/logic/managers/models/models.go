package models

import "src/objects"

type StudentFullInfo struct {
	Student objects.Student `json:"student"`
}

type ThingFullInfo struct {
	Thing objects.Thing `json:"thing"`
}
