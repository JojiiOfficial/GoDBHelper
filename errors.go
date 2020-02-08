package godbhelper

import "errors"

var (
	//ErrDBNotSupported error if database is not supported
	ErrDBNotSupported = errors.New("Database not supported")
)
