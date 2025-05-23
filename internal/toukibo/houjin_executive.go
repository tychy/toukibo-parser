package toukibo

import (
	"fmt"
	"strings"
)

type HoujinExecutiveValue struct {
	Name       string `yaml:"Name"`
	Position   string `yaml:"Position"`
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
		b.WriteString("  - Name: " + hv.Name + "\n")
		b.WriteString("    Position: " + hv.Position + "\n")
		if DebugOn {
			b.WriteString("    IsValid: " + fmt.Sprintf("%v", hv.IsValid) + "\n")
			b.WriteString("    RegisterAt: " + hv.RegisterAt + "\n")
			b.WriteString("    ResignedAt: " + hv.ResignedAt + "\n")
		}
	}
	return b.String()
}
