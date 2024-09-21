package toukibo

import (
	"strconv"
	"strings"
)

type HoujinPreferredStock struct {
	Type   string `yaml:"Type"`
	Amount int    `yaml:"Amount"`
}

type HoujinStock struct {
	Total     int
	Preferred []HoujinPreferredStock
}

func (s HoujinStock) Sum() int {
	sum := 0
	for _, p := range s.Preferred {
		sum += p.Amount
	}
	return sum
}

func (s HoujinStock) String() string {
	var b strings.Builder
	for _, p := range s.Preferred {
		b.WriteString("  - Type: " + p.Type + "\n")
		b.WriteString("    Amount: " + strconv.Itoa(p.Amount) + "\n")
	}
	return b.String()
}
