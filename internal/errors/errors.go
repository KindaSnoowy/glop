package customErrors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("expired token")
