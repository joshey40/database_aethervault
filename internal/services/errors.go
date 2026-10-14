package services

import "errors"

var ErrInvalidRequest = errors.New("Invalid Request")
var ErrUserAlreadyExists = errors.New("The user already exists")
var ErrInternalError = errors.New("Internal Error")
var ErrInvalidLoginData = errors.New("Username or password wrong")
