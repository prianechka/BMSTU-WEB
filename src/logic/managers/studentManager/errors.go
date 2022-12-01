package studentManager

import "errors"

var (
	ThingHasOwnerErr     = errors.New("thing has owner")
	StudentIsNotOwnerErr = errors.New("thing has owner")
)
