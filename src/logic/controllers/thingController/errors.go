package thingController

import "errors"

var (
	ThingAlreadyExistErr = errors.New("already existed")
	ThingNotFoundErr     = errors.New("thing not found")
	BadSrcRoomErr        = errors.New("bad src room")
	BadDstRoomErr        = errors.New("bad dst room")
)
