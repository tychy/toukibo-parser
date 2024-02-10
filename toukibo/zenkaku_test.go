package toukibo_test

import (
	"testing"

	"github.com/tychy/toukibo_parser/toukibo"
)

func TestZenkakuToHankaku(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "１２３４円",
			expected: "1234円",
		},
	}

	for _, tc := range testCases {
		actual := toukibo.ZenkakuToHankaku(tc.input)
		if actual != tc.expected {
			t.Errorf("ZenkakuToHankaku(%s) = %s; want %s", tc.input, actual, tc.expected)
		}
	}
}
