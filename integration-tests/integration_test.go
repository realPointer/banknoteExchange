package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	v1 "github.com/realPointer/banknoteExchange/internal/controller/http/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get("http://banknote_exchange:8080/health")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestExchangeBanknotes(t *testing.T) {
	serverURL := "http://banknote_exchange:8080"

	requestBody := v1.ExchangeRequest{
		Amount:    400,
		Banknotes: []int{5000, 2000, 1000, 500, 200, 100, 50},
	}
	requestJSON, err := json.Marshal(requestBody)
	require.NoError(t, err)

	resp, err := http.Post(serverURL+"/v1/exchange", "application/json", bytes.NewBuffer(requestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var exchangeResponse v1.ExchangeResponse
	err = json.NewDecoder(resp.Body).Decode(&exchangeResponse)
	require.NoError(t, err)

	expectedExchanges := [][]int{
		{200, 200},
		{200, 100, 100},
		{200, 100, 50, 50},
		{200, 50, 50, 50, 50},
		{100, 100, 100, 100},
		{100, 100, 100, 50, 50},
		{100, 100, 50, 50, 50, 50},
		{100, 50, 50, 50, 50, 50, 50},
		{50, 50, 50, 50, 50, 50, 50, 50},
	}

	assert.ElementsMatch(t, expectedExchanges, exchangeResponse.Exchanges)
}
