package thingController

import (
	"database/sql"
	"src/db/thingRepo"
	"src/objects"
	appErrors "src/utils/error"
)

type ThingController struct {
	Repo thingRepo.ThingRepo
}

func (tc *ThingController) AddThing(markNumber int, thingType string) error {
	things, err := tc.Repo.GetThings(objects.Null, objects.Null)
	if err == nil {
		for _, thing := range things {
			if thing.GetMarkNumber() == markNumber {
				err = appErrors.ThingAlreadyExistErr
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

func (tc *ThingController) GetThings(page, size int) ([]objects.Thing, error) {
	return tc.Repo.GetThings(page, size)
}

func (tc *ThingController) GetFreeThings(page, size int) ([]objects.Thing, error) {
	resultFreeThings := make([]objects.Thing, objects.Empty)
	things, err := tc.Repo.GetThings(page, size)
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
	if err == sql.ErrNoRows {
		err = appErrors.ThingNotFoundErr
	}
	return thing, err
}

func (tc *ThingController) DeleteThing(id int) error {
	thing, err := tc.Repo.GetThing(id)
	if err == nil {
		if thing.GetID() != objects.None {
			err = tc.Repo.DeleteThing(id)
		}
	} else if err == sql.ErrNoRows {
		err = appErrors.ThingNotFoundErr
	}
	return err
}

func (tc *ThingController) GetThingRoom(thingID int) (int, error) {
	var result = objects.None
	tmpThing, err := tc.Repo.GetThing(thingID)
	if err == nil {
		result = tmpThing.GetRoomID()
	} else if err == sql.ErrNoRows {
		err = appErrors.ThingNotFoundErr
	}
	return result, err
}

func (tc *ThingController) TransferThing(thingID, srcRoomID, dstRoomID int) error {
	tmpThing, err := tc.Repo.GetThing(thingID)
	if err == nil {
		if tmpThing.GetRoomID() != srcRoomID {
			err = appErrors.BadSrcRoomErr
		} else if dstRoomID == objects.None || dstRoomID == srcRoomID {
			err = appErrors.BadDstRoomErr
		} else {
			err = tc.Repo.TransferThingRoom(thingID, srcRoomID, dstRoomID)
		}
	} else if err == sql.ErrNoRows {
		err = appErrors.ThingNotFoundErr
	}
	return err
}

func (tc *ThingController) GetCurrentOwner(thingID int) (int, error) {
	var result = objects.None
	tmpThing, err := tc.Repo.GetThing(thingID)
	if err == nil {
		result = tmpThing.GetOwnerID()
	} else if err == sql.ErrNoRows {
		err = appErrors.ThingNotFoundErr
	}
	return result, err
}

func (tc *ThingController) GetThingIDByMarkNumber(markNumber int) (int, error) {
	id, err := tc.Repo.GetThingIDByMarkNumber(markNumber)
	if err == sql.ErrNoRows {
		err = appErrors.ThingNotFoundErr
	}
	return id, err
}
