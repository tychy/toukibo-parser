package toukibo_test

import (
	"testing"

	"github.com/tychy/toukibo-parser/toukibo"
)

func TestYenToNumber(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{
			input:    "1234円",
			expected: 1234,
		},
		{
			input:    "1234万円",
			expected: 12340000,
		},
		{
			input:    "1234億円",
			expected: 123400000000,
		},
		{
			input:    "1234兆円",
			expected: 1234000000000000,
		},
		{
			input:    "1234兆円1234億5678万9012円",
			expected: 1234123456789012,
		},
		{
			input:    "１２３４円",
			expected: 1234,
		},
		{
			input:    "１２３４万円",
			expected: 12340000,
		},
		{
			input:    "１２３４億円",
			expected: 123400000000,
		},
		{
			input:    "１２３４兆円",
			expected: 1234000000000000,
		},
		{
			input:    "１２３４兆円１２３４億５６７８万９０１２円",
			expected: 1234123456789012,
		},
		{
			input:    "金1円",
			expected: 1,
		},
		{
			input:    "",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		actual := toukibo.YenToNumber(tc.input)
		if actual != tc.expected {
			t.Errorf("YenToNumber(%s) = %d; want %d", tc.input, actual, tc.expected)
		}
	}

}
