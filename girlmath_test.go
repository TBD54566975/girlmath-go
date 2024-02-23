package girlmath_test

import (
	"math/big"
	"testing"

	"github.com/tbd54566975/girlmath"
)

func TestCalculatePayoutAmountSubunitsWithPayoutCurrencySpotPrice(t *testing.T) {
	var tests = []struct {
		payinCurrency                   string
		payoutCurrency                  string
		payinAmountSubunits             int64
		payoutCurrencyUnitsPerPayinUnit string
		expectedPayoutAmountSubunits    *big.Int
	}{
		{"USD", "BTC", 10000, "0.0000325291", big.NewInt(325291)},
		{"BTC", "USD", 325291, "30,741.70", big.NewInt(9999)},
		{"BTC", "KES", 2500, "4560829.07955", big.NewInt(11402)},
		{"USD", "BTC", 10000, "0.00134825401", big.NewInt(13482540)},
		{"USD", "BTC", 10000, "0.00003252984", big.NewInt(325298)},
	}

	for _, tt := range tests {
		t.Run(tt.payinCurrency+" to "+tt.payoutCurrency, func(t *testing.T) {
			result, err := girlmath.CalculatePayoutAmountSubunitsWithPayoutCurrencySpotPrice(tt.payinCurrency, tt.payoutCurrency, tt.payinAmountSubunits, tt.payoutCurrencyUnitsPerPayinUnit)
			if err != nil {
				t.Errorf("calculatePayoutAmountSubunitsWithPayoutCurrencySpotPrice() error = %v", err)
				return
			}
			if result.Cmp(tt.expectedPayoutAmountSubunits) != 0 {
				t.Errorf("Expected %v, got %v", tt.expectedPayoutAmountSubunits, result)
			}
		})
	}
}

func TestCalculatePayoutAmountSubunitsWithPayinCurrencySpotPrice(t *testing.T) {
	tests := []struct {
		payinCurrency                   string
		payoutCurrency                  string
		payinAmountSubunits             int64
		payinCurrencyUnitsPerPayoutUnit string
		expectedPayoutAmountSubunits    *big.Int
	}{
		{"USD", "BTC", 10000, "30,741.70", big.NewInt(325291)},
		{"BTC", "USD", 325291, "0.0000325291", big.NewInt(10000)},
		{"BTC", "KES", 2500, "0.000000219258381", big.NewInt(11402)},
	}

	for _, tt := range tests {
		t.Run(tt.payinCurrency+" to "+tt.payoutCurrency, func(t *testing.T) {
			result, err := girlmath.CalculatePayoutAmountSubunitsWithPayinCurrencySpotPrice(tt.payinCurrency, tt.payoutCurrency, tt.payinAmountSubunits, tt.payinCurrencyUnitsPerPayoutUnit)
			if err != nil {
				t.Errorf("calculatePayoutAmountSubunitsWithPayinCurrencySpotPrice() error = %v", err)
				return
			}
			if result.Cmp(tt.expectedPayoutAmountSubunits) != 0 {
				t.Errorf("Expected %v, got %v", tt.expectedPayoutAmountSubunits, result)
			}
		})
	}
}

func TestConvertSubunitsToUnits(t *testing.T) {
	tests := []struct {
		amountSubunits int64
		currencyCode   string
		expected       string
	}{
		{20000, "USD", "200.00"},
		{12311, "USD", "123.11"},
		{100000000, "BTC", "1.00000000"},
		{123456789, "BTC", "1.23456789"},
		{506379, "BTC", "0.00506379"},
	}

	for _, tt := range tests {
		t.Run(tt.currencyCode, func(t *testing.T) {
			result, err := girlmath.ConvertSubunitsToUnits(tt.amountSubunits, tt.currencyCode)
			if err != nil {
				t.Errorf("convertSubunitsToUnits() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvertUnitsToSubunits(t *testing.T) {
	tests := []struct {
		amountUnits  string
		currencyCode string
		expected     string
	}{
		{"1.23", "USD", "123"},
		{"1.2", "USD", "120"},
		{"1", "USD", "100"},
		{"1.23456789", "BTC", "123456789"},
		{"1.234", "BTC", "123400000"},
		{"1", "BTC", "100000000"},
		{"0.00506379", "BTC", "506379"},
		{"10.00506379", "BTC", "1000506379"},
		{"0.0050637", "BTC", "506370"},
		{".0050637", "BTC", "506370"},
		{"0.25", "USDC", "250000"},
	}

	for _, tt := range tests {
		t.Run(tt.currencyCode, func(t *testing.T) {
			result, err := girlmath.ConvertUnitsToSubunits(tt.amountUnits, tt.currencyCode)
			if err != nil {
				t.Errorf("convertUnitsToSubunits() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
