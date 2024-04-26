package service

import (
	"errors"

	"github.com/rs/zerolog"
)

type ExchangeService struct {
	l *zerolog.Logger
}

func NewExchangeService(l *zerolog.Logger) *ExchangeService {
	return &ExchangeService{
		l: l,
	}
}

func (s *ExchangeService) CalculateExchanges(amount int, banknotes []int) ([][]int, error) {
	s.l.Debug().Msgf("CalculateExchanges called with amount: %d, banknotes: %v", amount, banknotes)

	var exchanges [][]int

	var backtrack func(banknoteIndex, remainingAmount int, currentCombination []int)
	backtrack = func(banknoteIndex, remainingAmount int, currentCombination []int) {
		s.l.Trace().Msgf("backtrack called with banknoteIndex: %d, remainingAmount: %d, currentCombination: %v",
			banknoteIndex, remainingAmount, currentCombination)

		if remainingAmount == 0 {
			exchanges = append(exchanges, append([]int{}, currentCombination...))
			s.l.Debug().Msgf("Found a valid exchange: %v", currentCombination)

			return
		}

		for i := banknoteIndex; i < len(banknotes); i++ {
			currentBanknote := banknotes[i]
			if currentBanknote <= remainingAmount {
				currentCombination = append(currentCombination, currentBanknote)
				s.l.Trace().Msgf("Trying banknote: %d, updated currentCombination: %v", currentBanknote, currentCombination)
				backtrack(i, remainingAmount-currentBanknote, currentCombination)
				currentCombination = currentCombination[:len(currentCombination)-1]
			}
		}
	}

	backtrack(0, amount, []int{})

	if len(exchanges) == 0 {
		s.l.Debug().Msg("No exchanges found")

		return nil, ErrNoExchanges
	}

	s.l.Debug().Msgf("Found %d possible exchanges. Exchanges: %v", len(exchanges), exchanges)

	return exchanges, nil
}

var ErrNoExchanges = errors.New("no exchanges found")
