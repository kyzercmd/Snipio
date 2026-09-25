package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid email or password")
	ErrDuplicateEmail     = errors.New("models: email already registered")
)
