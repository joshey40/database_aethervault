package services

import "errors"

var ErrInvalidRequest = errors.New("Invalid Request")
var ErrUserAlreadyExists = errors.New("The user already exists")
var ErrInternal = errors.New("Internal Error")
var ErrInvalidLoginData = errors.New("Username or password wrong")
var ErrDeviceNotFound = errors.New("Device was not found")
var ErrUserNotFound = errors.New("User was not found")
