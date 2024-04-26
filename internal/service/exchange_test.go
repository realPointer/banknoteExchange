package service_test

import (
	"testing"

	"github.com/realPointer/banknoteExchange/internal/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCalculateExchanges(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		amount         int
		banknotes      []int
		expectedResult [][]int
		expectedError  error
	}{
		{
			name:      "Valid exchanges",
			amount:    400,
			banknotes: []int{5000, 2000, 1000, 500, 200, 100, 50},
			expectedResult: [][]int{
				{200, 200},
				{200, 100, 100},
				{200, 100, 50, 50},
				{200, 50, 50, 50, 50},
				{100, 100, 100, 100},
				{100, 100, 100, 50, 50},
				{100, 100, 50, 50, 50, 50},
				{100, 50, 50, 50, 50, 50, 50},
				{50, 50, 50, 50, 50, 50, 50, 50},
			},
			expectedError: nil,
		},
		{
			name:           "No exchanges",
			amount:         7,
			banknotes:      []int{5},
			expectedResult: nil,
			expectedError:  service.ErrNoExchanges,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := &zerolog.Logger{}
			exchangeService := service.NewExchangeService(logger)

			result, err := exchangeService.CalculateExchanges(tc.amount, tc.banknotes)

			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
