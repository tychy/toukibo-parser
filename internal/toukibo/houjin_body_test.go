package toukibo

import "testing"

func TestDissolvedHoujinWithoutLiquidator(t *testing.T) {
	h := &HoujinBody{
		HoujinDissolvedAt: "令和3年12月15日",
		HoujinExecutive: []HoujinExecutiveValueArray{{
			{Name: "山田太郎", Position: "監査役", IsValid: true},
		}},
	}

	representatives, err := h.GetHoujinRepresentatives()
	if err != nil {
		t.Fatal(err)
	}
	if len(representatives) != 0 {
		t.Fatalf("representative count: want 0, got %d", len(representatives))
	}
}

func TestContinuedHoujinUsesCurrentRepresentative(t *testing.T) {
	h := &HoujinBody{
		HoujinDissolvedAt: "令和1年12月11日",
		HoujinContinuedAt: "令和2年7月1日",
		HoujinExecutive: []HoujinExecutiveValueArray{{
			{Name: "山田太郎", Position: "代表取締役", IsValid: true},
		}},
	}

	representatives, err := h.GetHoujinRepresentatives()
	if err != nil {
		t.Fatal(err)
	}
	if len(representatives) != 1 || representatives[0].Name != "山田太郎" {
		t.Fatalf("representatives: got %+v", representatives)
	}
}
