package toukibo_parser

import "testing"

// sample1023 は令和6年5月31日に辞任した取締役（代表取締役を兼ねる）を含む実サンプル。
func TestParseFormerExecutives(t *testing.T) {
	h, err := ParseByPDFPath("testdata/pdf/sample1023.pdf")
	if err != nil {
		t.Fatal(err)
	}

	former := h.GetHoujinFormerExecutives()
	if len(former) != 2 {
		t.Fatalf("former executive count: want 2, got %d (%+v)", len(former), former)
	}
	for i, v := range former {
		if v.Name != "見上耕" {
			t.Fatalf("former executive[%d].Name = %q, want 見上耕", i, v.Name)
		}
		if v.ResignedAt != "令和6年5月31日" {
			t.Fatalf("former executive[%d].ResignedAt = %q, want 令和6年5月31日", i, v.ResignedAt)
		}
		if v.ResignedReason != "辞任" {
			t.Fatalf("former executive[%d].ResignedReason = %q, want 辞任", i, v.ResignedReason)
		}
		if v.IsValid {
			t.Fatalf("former executive[%d] must not be valid", i)
		}
	}

	// 現任の役員に退任済みの行が混ざらないこと
	execs, err := h.GetHoujinExecutives()
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range execs {
		if v.ResignedAt != "" || v.ResignedReason != "" {
			t.Fatalf("executive[%d] must not have resigned fields: %+v", i, v)
		}
	}
}

// sample1023 の役員欄は、現任6行と、令和6年5月31日に辞任した取締役（代表取締役を兼ねる）2行で構成される。
func TestParseExecutiveHistory(t *testing.T) {
	h, err := ParseByPDFPath("testdata/pdf/sample1023.pdf")
	if err != nil {
		t.Fatal(err)
	}

	history := h.GetHoujinExecutiveHistory()
	execs, err := h.GetHoujinExecutives()
	if err != nil {
		t.Fatal(err)
	}
	former := h.GetHoujinFormerExecutives()
	if len(history) != len(execs)+len(former) {
		t.Fatalf("history = current + former: history=%d current=%d former=%d", len(history), len(execs), len(former))
	}
	if len(former) != 2 {
		t.Fatalf("former count: want 2, got %d (%+v)", len(former), former)
	}

	current := 0
	for i, v := range history {
		if v.IsValid {
			current++
			if v.ResignedAt != "" || v.ResignedReason != "" {
				t.Fatalf("history[%d] must not have resigned fields: %+v", i, v)
			}
			continue
		}
		if v.Name != "見上耕" || v.ResignedReason != "辞任" || v.ResignedAt != "令和6年5月31日" {
			t.Fatalf("history[%d] = %+v, want 見上耕 辞任 令和6年5月31日", i, v)
		}
	}
	if current != len(execs) {
		t.Fatalf("current rows in history: want %d, got %d", len(execs), current)
	}
}
