package service

import (
	"github.com/rs/zerolog"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Exchange interface {
	CalculateExchanges(amount int, banknotes []int) ([][]int, error)
}

type Services struct {
	Exchange
}

func NewServices(l *zerolog.Logger) *Services {
	return &Services{
		Exchange: NewExchangeService(l),
	}
}
