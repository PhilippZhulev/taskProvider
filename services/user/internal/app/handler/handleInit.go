package handler

import (
	"errors"

	"gitlab.com/taskProvider/services/user/internal/app/helpers"
)


var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

// Init ...
type Init struct {
	hesh helpers.Hesh
}