package domain

import "errors"

var NotFoundError error = errors.New("not found")
var AlreadyExistsError error = errors.New("already exists")
