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
		{
			// sample1082
			input:    "発行済株式の総数普通株式　　　　　５１３８株Ａ種優先株式　　　１２５０株Ｂ種優先株式　　　３２５０株",
			expected: 5138,
		},
	}

	for _, tc := range testCases {
		if got := toukibo.StockToNumber(tc.input); got != tc.expected {
			t.Errorf("StockToNumber(%s) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}
