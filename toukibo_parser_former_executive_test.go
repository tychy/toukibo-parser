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

// 退任日（または事由）を更正した役員。元の退任は下線で抹消され、更正の記載は「◯日更正」で終わり登記日を伴わない。
// 更正後の退任が履歴に入り、現任には入らないこと。
func TestParseExecutiveHistoryCorrectedResignation(t *testing.T) {
	type ended struct{ name, position, resignedAt, reason string }
	tests := []struct {
		sample string
		want   []ended
	}{
		// 令和8年4月1日辞任を令和8年3月31日辞任に更正
		{"sample1781", []ended{{"後藤貞男", "取締役", "令和8年3月31日", "辞任"}}},
		// 令和6年3月31日辞任を令和6年4月1日辞任に更正（取締役・代表取締役の両方）
		{"sample1782", []ended{
			{"近藤雄一郎", "取締役", "令和6年4月1日", "辞任"},
			{"川嵜靖之", "取締役", "令和6年4月1日", "辞任"},
			{"近藤雄一郎", "代表取締役", "令和6年4月1日", "辞任"},
			{"川嵜靖之", "代表取締役", "令和6年4月1日", "辞任"},
		}},
		// 代表取締役の退任を辞任に更正（取締役としては辞任の登記が抹消されて現任）
		{"sample1783", []ended{{"中島克宣", "代表取締役", "令和5年3月7日", "辞任"}}},
	}

	for _, tt := range tests {
		t.Run(tt.sample, func(t *testing.T) {
			h, err := ParseByPDFPath("testdata/pdf/" + tt.sample + ".pdf")
			if err != nil {
				t.Fatal(err)
			}
			history := h.GetHoujinExecutiveHistory()
			for _, w := range tt.want {
				found := false
				for _, v := range history {
					if v.Name == w.name && v.Position == w.position {
						if v.IsValid {
							continue
						}
						if v.ResignedAt != w.resignedAt || v.ResignedReason != w.reason {
							t.Fatalf("%s %s: got (%q, %q), want (%q, %q)", w.name, w.position, v.ResignedAt, v.ResignedReason, w.resignedAt, w.reason)
						}
						found = true
					}
				}
				if !found {
					t.Fatalf("%s %s が履歴に無い: %+v", w.name, w.position, history)
				}
			}

			execs, err := h.GetHoujinExecutives()
			if err != nil {
				t.Fatal(err)
			}
			for _, w := range tt.want {
				for _, v := range execs {
					if v.Name == w.name && v.Position == w.position {
						t.Fatalf("%s %s は現任に含まれない: %+v", w.name, w.position, v)
					}
				}
			}
		})
	}
}
