package toukibo_test

import (
	"testing"

	"github.com/tychy/toukibo-parser/internal/toukibo"
)

func TestGetStockNumber(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
		remain   string
	}{
		{
			input:    "5万株",
			expected: 50000,
			remain:   "",
		},
		{
			input:    "1万9013株各種の株式の数普通株式",
			expected: 19013,
			remain:   "各種の株式の数普通株式",
		},
	}

	for _, tc := range testCases {
		got, remain := toukibo.GetStockNumber(tc.input)
		if got != tc.expected {
			t.Errorf("GetStockNumber(%s) = %d; want %d", tc.input, got, tc.expected)
		}
		if remain != tc.remain {
			t.Errorf("GetStockNumber(%s) = %s; want %s", tc.input, remain, tc.remain)
		}
	}
}

func TestGetHoujinStockWithNumberedAndDividendPreferredStock(t *testing.T) {
	stock := toukibo.GetHoujinStock("発行済株式の総数600193株各種の株式の数第1回B種優先株式60万株甲種優先配当株式188株丙種優先配当株式5株")

	if stock.Total != 600193 {
		t.Fatalf("total = %d; want 600193", stock.Total)
	}
	want := []toukibo.HoujinPreferredStock{
		{Type: "第1回B種優先株式", Amount: 600000},
		{Type: "甲種優先配当株式", Amount: 188},
		{Type: "丙種優先配当株式", Amount: 5},
	}
	if len(stock.Preferred) != len(want) {
		t.Fatalf("preferred stock count = %d; want %d: %+v", len(stock.Preferred), len(want), stock.Preferred)
	}
	for i := range want {
		if stock.Preferred[i] != want[i] {
			t.Errorf("preferred stock[%d] = %+v; want %+v", i, stock.Preferred[i], want[i])
		}
	}
}

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
			input:    "発行済株式の総数１万９０１３株各種の株式の数普通株式　　　１万５２５株Ａ種優先株式　　３４３４株Ａ１種優先株式　１１００株Ａ２種優先株式　２１０３株Ａ３種優先株式　１８５１株",
			expected: 19013,
		},
		{
			// sample1082
			input:    "発行済株式の総数普通株式　　　　　５１３８株Ａ種優先株式　　　１２５０株Ｂ種優先株式　　　３２５０株",
			expected: 9638,
		},
	}

	for _, tc := range testCases {
		if got := toukibo.StockToNumber(tc.input); got != tc.expected {
			t.Errorf("StockToNumber(%s) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}
