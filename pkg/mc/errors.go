package mc

import "errors"

var ErrNotFound = errors.New("not found")
var ErrValidation = errors.New("validation failed")
var ErrNotImplemented = errors.New("not implemented")
var ErrGlobusAPI = errors.New("globus api error")
