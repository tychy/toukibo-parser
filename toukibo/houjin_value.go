package toukibo

import (
	"fmt"
	"strings"
)

type HoujinValue struct {
	Value      string
	IsValid    bool
	RegisterAt string
}

func (h *HoujinValue) String() string {
	return fmt.Sprintf("value: %s, isValid: %v, registerAt: %s",
		h.Value, h.IsValid, h.RegisterAt)
}

type HoujinValueArray []HoujinValue

func (hva HoujinValueArray) String() string {
	var b strings.Builder
	for _, hv := range hva {
		b.WriteString("{")
		b.WriteString(hv.String())
		b.WriteString("},")
	}
	return b.String()
}
