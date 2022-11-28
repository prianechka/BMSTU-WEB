package thingController

import (
	"src/db/thingRepo"
	"src/objects"
)

type ThingController struct {
	Repo thingRepo.ThingRepo
}

func (tc *ThingController) AddThing(markNumber int, thingType string) error {
	things, err := tc.Repo.GetThings()
	if err == nil {
		for _, thing := range things {
			if thing.GetMarkNumber() == markNumber {
				err = ThingAlreadyExistErr
				break
			}
		}
		if err == nil {
			thingDTO := objects.NewThingDTO(markNumber, thingType)
			err = tc.Repo.AddThing(thingDTO)
		}
	}
	return err
}

func (tc *ThingController) GetThings() ([]objects.Thing, error) {
	return tc.Repo.GetThings()
}

func (tc *ThingController) GetFreeThings() ([]objects.Thing, error) {
	resultFreeThings := make([]objects.Thing, objects.Empty)
	things, err := tc.Repo.GetThings()
	if err == nil {
		for _, thing := range things {
			if thing.GetOwnerID() == objects.None {
				resultFreeThings = append(resultFreeThings, thing)
			}
		}
	}
	return resultFreeThings, err
}

func (tc *ThingController) GetThing(id int) (objects.Thing, error) {
	thing, err := tc.Repo.GetThing(id)
	if err == nil && thing.GetID() == objects.None {
		err = ThingNotFoundErr
	}
	return thing, err
}

func (tc *ThingController) DeleteThing(id int) error {
	thing, err := tc.Repo.GetThing(id)
	if err == nil {
		if thing.GetID() != objects.None {
			err = tc.Repo.DeleteThing(id)
		} else {
			err = ThingNotFoundErr
		}
	}
	return err
}

func (tc *ThingController) GetThingRoom(thingID int) (int, error) {
	var result = objects.None
	tmpThing, err := tc.Repo.GetThing(thingID)
	if err == nil {
		if tmpThing.GetID() == objects.None {
			err = ThingNotFoundErr
		} else {
			result = tmpThing.GetRoomID()
		}
	}
	return result, err
}

func (tc *ThingController) TransferThing(thingID, srcRoomID, dstRoomID int) error {
	tmpThing, err := tc.Repo.GetThing(thingID)
	if err == nil {
		if tmpThing.GetID() == objects.None {
			err = ThingNotFoundErr
		} else if tmpThing.GetRoomID() != srcRoomID {
			err = BadSrcRoomErr
		} else if dstRoomID == objects.None {
			err = BadDstRoomErr
		} else {
			err = tc.Repo.TransferThingRoom(thingID, srcRoomID, dstRoomID)
		}
	}
	return err
}

func (tc *ThingController) GetCurrentOwner(thingID int) (int, error) {
	var result = objects.None
	tmpThing, err := tc.Repo.GetThing(thingID)
	if err == nil {
		if tmpThing.GetID() == objects.None {
			err = ThingNotFoundErr
		} else {
			result = tmpThing.GetOwnerID()
		}
	}
	return result, err
}

func (tc *ThingController) GetThingIDByMarkNumber(markNumber int) (int, error) {
	id, err := tc.Repo.GetThingIDByMarkNumber(markNumber)
	if id == objects.None {
		err = ThingNotFoundErr
	}
	return id, err
}
