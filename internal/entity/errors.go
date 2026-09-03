package entity

import "errors"

var (
	ErrInvalidEntityId = errors.New("invalid entity id")
	ErrEntityNotFound  = errors.New("entity not found")
)
