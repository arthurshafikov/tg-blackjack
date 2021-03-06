package core

import "errors"

var (
	ErrNotFound    = errors.New("404 Not Found")
	ErrServerError = errors.New("500 Server Error")

	ErrChatRegistered = errors.New("chat already registered")
	ErrActiveGame     = errors.New("chat has active game")
	ErrNoActiveGame   = errors.New("chat has no active game")
	ErrDeckEmpty      = errors.New("deck is empty")
	ErrAlreadyStopped = errors.New("you are already stopped drawing")
	ErrBusted         = errors.New("more than 21")
	ErrCantDraw       = errors.New("player can't draw more cards")
)
