package model

import "errors"

var (
	ErrNotFound      = errors.New("bookmark not found")
	ErrURLRequired   = errors.New("url is required")
	ErrTitleRequired = errors.New("title is required")
)
