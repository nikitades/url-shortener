package usecases

import (
	"errors"
)

var BadRequestError error = errors.New("bad request")
var NotFoundError error = errors.New("not found")
