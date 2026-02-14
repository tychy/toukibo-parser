package toukibo

import (
	"strings"
	"testing"
)

func TestParseBodyMainParsesHoujinPurposeWithHistory(t *testing.T) {
	pair := " ┃目　的　　　　　│　１．旧目的　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　┃ " +
		revert1 +
		" ┃　　　　　　　　│　１．新目的　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　┃ " +
		"┃　　　　　　　　│　　　　　　　　　令和　３年　６月２５日変更　　令和　３年　７月１２日登記┃ "

	h, err := ParseBody([]string{pair})
	if err != nil {
		t.Fatalf("ParseBody returned error: %v", err)
	}

	if len(h.HoujinPurpose) != 2 {
		t.Fatalf("expected 2 purpose history entries, got %d", len(h.HoujinPurpose))
	}

	if h.HoujinPurpose[0].IsValid {
		t.Fatalf("expected first purpose entry to be invalid")
	}
	if h.HoujinPurpose[0].Value != "１．旧目的" {
		t.Fatalf("unexpected first purpose value: %q", h.HoujinPurpose[0].Value)
	}

	if !h.HoujinPurpose[1].IsValid {
		t.Fatalf("expected latest purpose entry to be valid")
	}
	if h.HoujinPurpose[1].Value != "１．新目的" {
		t.Fatalf("unexpected latest purpose value: %q", h.HoujinPurpose[1].Value)
	}
	if h.HoujinPurpose[1].RegisterAt != "令和３年７月１２日" {
		t.Fatalf("unexpected registerAt: %q", h.HoujinPurpose[1].RegisterAt)
	}
}

func TestParseBodyMainParsesHoujinPurposeTou(t *testing.T) {
	pair := " ┃目的等　　　　　│　事業　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　┃ " +
		"┃　　　　　　　　│　⑴　テスト事業　　　　　　　　　　　　　　　　　　　　　　　　　　　　　┃ "

	h, err := ParseBody([]string{pair})
	if err != nil {
		t.Fatalf("ParseBody returned error: %v", err)
	}

	if len(h.HoujinPurpose) != 1 {
		t.Fatalf("expected 1 purpose entry, got %d", len(h.HoujinPurpose))
	}
	if !h.HoujinPurpose[0].IsValid {
		t.Fatalf("expected purpose entry to be valid")
	}
	if !strings.Contains(h.HoujinPurpose[0].Value, "事業") || !strings.Contains(h.HoujinPurpose[0].Value, "テスト事業") {
		t.Fatalf("unexpected purpose value: %q", h.HoujinPurpose[0].Value)
	}
}

func TestGetHoujinPurposeReturnsLatestValidValue(t *testing.T) {
	h := &Houjin{body: &HoujinBody{HoujinPurpose: HoujinValueArray{
		{Value: "１．旧目的", IsValid: false},
		{Value: "１．新目的", IsValid: true},
	}}}

	if got := h.GetHoujinPurpose(); got != "１．新目的" {
		t.Fatalf("GetHoujinPurpose() = %q, want %q", got, "１．新目的")
	}

	history := h.GetHoujinPurposeHistory()
	if len(history) != 2 {
		t.Fatalf("GetHoujinPurposeHistory length = %d, want 2", len(history))
	}
}
