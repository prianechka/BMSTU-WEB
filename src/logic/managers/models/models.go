package models

import "src/objects"

type StudentFullInfo struct {
	Student objects.Student
	Room    objects.Room
}

type ThingFullInfo struct {
	Thing objects.Thing
	Room  objects.Room
}
