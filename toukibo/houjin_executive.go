package toukibo

import (
	"fmt"
	"strings"
)

type HoujinExecutiveValue struct {
	one   string
	two   string
	three []string

	Name       string
	Position   string
	Address    string
	IsValid    bool
	RegisterAt string
	ResignedAt string
}

func (h *HoujinExecutiveValue) String() string {
	return fmt.Sprintf("name: %s, position: %s, address: %s, isValid: %v, registerAt: %s, resignedAt: %s",
		h.Name, h.Position, h.Address, h.IsValid, h.RegisterAt, h.ResignedAt)
}

type HoujinExecutiveValueArray []HoujinExecutiveValue

func (hva HoujinExecutiveValueArray) String() string {
	var b strings.Builder
	for _, hv := range hva {
		b.WriteString("{")
		b.WriteString(hv.String())
		b.WriteString("},")
	}
	return b.String()
}
