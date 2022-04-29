package core

import "errors"

var (
	ErrNotFound    = errors.New("404 Not Found")
	ErrServerError = errors.New("500 Server Error")

	ErrActiveGame = errors.New("chat has active game")
	ErrDeckEmpty  = errors.New("deck is empty")
)
