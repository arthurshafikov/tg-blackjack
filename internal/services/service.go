package services

import (
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type Logger interface {
	Error(err error)
}

type Services struct {
}

type Deps struct {
	Repository *repository.Repository
	Logger
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
