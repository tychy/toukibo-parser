package toukibo_test

import (
	"testing"

	"github.com/tychy/toukibo-parser/toukibo"
)

func TestStockToNumber(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{
			input:    "発行済株式の総数5万株",
			expected: 50000,
		},
	}

	for _, tc := range testCases {
		if got := toukibo.StockToNumber(tc.input); got != tc.expected {
			t.Errorf("StockToNumber(%s) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}
