// Package customerrors -> erros customizados do projeto
package customerrors

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)
