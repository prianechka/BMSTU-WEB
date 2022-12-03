package appErrors

import "errors"

var (
	RoomNotFoundErr         = errors.New("room not found")
	BadAccIDErr             = errors.New("bad accID")
	BadStudentParamsErr     = errors.New("bad params")
	StudentAlreadyInBaseErr = errors.New("student is already in base")
	StudentNotFoundErr      = errors.New("student not found")
	StudentAlreadyLiveErr   = errors.New("student is already living ")
	StudentNotLivingErr     = errors.New("student doesn't live")
	UserNotFoundErr         = errors.New("user not found")
	LoginOccupedErr         = errors.New("login is already exist")
	BadUserParamsErr        = errors.New("bad params")
	BadThingParamsErr       = errors.New("bad params")
	ThingAlreadyExistErr    = errors.New("already existed")
	ThingNotFoundErr        = errors.New("thing not found")
	BadSrcRoomErr           = errors.New("bad src room")
	BadDstRoomErr           = errors.New("bad dst room")
	PasswordNotEqualErr     = errors.New("passwords aren't equal ")
	ThingHasOwnerErr        = errors.New("thing has owner")
	StudentIsNotOwnerErr    = errors.New("thing has owner")
)
