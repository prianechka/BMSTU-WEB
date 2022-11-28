package studentController

import "errors"

var (
	BadAccIDErr             = errors.New("bad accID")
	BadParamsErr            = errors.New("bad params")
	StudentAlreadyInBaseErr = errors.New("student is already in base")
	StudentNotFoundErr      = errors.New("student not found")
	StudentAlreadyLiveErr   = errors.New("student is already living ")
	StudentNotLivingErr     = errors.New("student doesn't live")
)
