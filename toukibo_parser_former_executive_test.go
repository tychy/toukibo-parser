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

// 登記が抹消され「◯日抹消により回復」で役員の地位が戻った実サンプル。
// 回復した行を「抹消」による退任として扱わないこと。
func TestParseExecutiveRestoredByErasure(t *testing.T) {
	tests := []struct {
		sample, name, position string
	}{
		// 辞任の登記が抹消され、回復が最後の記載（現任）
		{"sample1784", "工藤一進", "取締役"},
		// 辞任の登記が抹消されて回復し、その後に重任
		{"sample1785", "伊東政子", "監査役"},
		// 住所移転の登記が抹消されて回復（現任）
		{"sample1786", "中井富士男", "取締役"},
	}

	for _, tt := range tests {
		t.Run(tt.sample, func(t *testing.T) {
			h, err := ParseByPDFPath("testdata/pdf/" + tt.sample + ".pdf")
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range h.GetHoujinExecutiveHistory() {
				if v.Name == tt.name && v.Position == tt.position && !v.IsValid {
					t.Fatalf("%s %s に終了行がある: %+v", tt.name, tt.position, v)
				}
			}
			execs, err := h.GetHoujinExecutives()
			if err != nil {
				t.Fatal(err)
			}
			current := 0
			for _, v := range execs {
				if v.Name == tt.name && v.Position == tt.position {
					current++
				}
			}
			if current != 1 {
				t.Fatalf("%s %s の現任行: want 1, got %d (%+v)", tt.name, tt.position, current, execs)
			}
		})
	}
}

// 監査役会の定めの廃止に伴い「◯日監査役会の定めの廃止により変更」と書き直された監査役の実サンプル。
// 変更の行を「廃止」による退任として扱わないこと。
func TestParseExecutiveChangedByAbolition(t *testing.T) {
	tests := []struct {
		sample, name string
	}{
		// 変更が最後の記載（現任）
		{"sample1534", "門松優介"},
		{"sample1209", "立木達"},
		{"sample1788", "馬場文明"},
		{"sample1788", "植村公彦"},
		// 変更の後に重任
		{"sample1787", "入江正邦"},
	}

	for _, tt := range tests {
		t.Run(tt.sample+"/"+tt.name, func(t *testing.T) {
			h, err := ParseByPDFPath("testdata/pdf/" + tt.sample + ".pdf")
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range h.GetHoujinExecutiveHistory() {
				if v.Name == tt.name && !v.IsValid {
					t.Fatalf("%s に終了行がある: %+v", tt.name, v)
				}
			}
			execs, err := h.GetHoujinExecutives()
			if err != nil {
				t.Fatal(err)
			}
			current := 0
			for _, v := range execs {
				if v.Name == tt.name && v.Position == "監査役" {
					current++
				}
			}
			if current != 1 {
				t.Fatalf("%s の現任行: want 1, got %d (%+v)", tt.name, current, execs)
			}
		})
	}
}

// 右の列で折り返され、日付と事由が別の行にある退任の実サンプル。
func TestParseExecutiveHistoryWrappedResignation(t *testing.T) {
	tests := []struct {
		sample, name, position, resignedAt, reason string
	}{
		// 「令和７年３月３１日職務」「執行者辞任」
		{"sample1789", "片岡宣秀", "職務執行者", "令和7年3月31日", "辞任"},
		// 「令和６年６月３０日取締」「役が１名となったため抹消」
		{"sample1790", "道面陽子", "代表取締役", "令和6年6月30日", "抹消"},
	}
	for _, tt := range tests {
		t.Run(tt.sample, func(t *testing.T) {
			h, err := ParseByPDFPath("testdata/pdf/" + tt.sample + ".pdf")
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, v := range h.GetHoujinExecutiveHistory() {
				if v.Name != tt.name || v.Position != tt.position || v.IsValid {
					continue
				}
				if v.ResignedAt != tt.resignedAt || v.ResignedReason != tt.reason {
					t.Fatalf("%s %s: got (%q, %q), want (%q, %q)", tt.name, tt.position, v.ResignedAt, v.ResignedReason, tt.resignedAt, tt.reason)
				}
				found = true
			}
			if !found {
				t.Fatalf("%s %s の終了行が無い", tt.name, tt.position)
			}
		})
	}
}

// 役員の任期の始まり（StartedAt）と登記日（RegisterAt）。
func TestParseExecutiveStartedAt(t *testing.T) {
	tests := []struct {
		sample, name, position         string
		current                        bool
		startedAt, registerAt, endedAt string
	}{
		// 重任の任期を辞任で終えた行。下線付き（取り消し）の辞任の登記日では上書きしない
		{"sample1781", "後藤貞男", "取締役", false, "令和7年5月22日", "令和7年6月9日", "令和8年3月31日"},
		// 辞任の登記が抹消され回復した現任の行。就任から引き継ぐ
		{"sample1784", "工藤一進", "取締役", true, "令和1年11月30日", "令和1年12月23日", ""},
		// 重任の登記が抹消され回復したあとに退任。抹消された重任（平成27年）ではなく、その前の重任から引き継ぐ
		{"sample1791", "横田国寿", "取締役", false, "平成17年6月25日", "令和7年8月6日", "平成19年6月25日"},
		{"sample1791", "横田国寿", "取締役", true, "令和7年8月6日", "令和7年8月6日", ""},
		// 「令和６年１１月１４日職務」「執行者就任」と折り返された就任
		{"sample1789", "平岡佳明", "職務執行者", true, "令和6年11月14日", "令和6年11月14日", ""},
	}
	for _, tt := range tests {
		t.Run(tt.sample+"/"+tt.name, func(t *testing.T) {
			h, err := ParseByPDFPath("testdata/pdf/" + tt.sample + ".pdf")
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, v := range h.GetHoujinExecutiveHistory() {
				if v.Name != tt.name || v.Position != tt.position || v.IsValid != tt.current {
					continue
				}
				found = true
				if v.StartedAt != tt.startedAt || v.RegisterAt != tt.registerAt || v.ResignedAt != tt.endedAt {
					t.Fatalf("got (started %q, registered %q, ended %q), want (%q, %q, %q)",
						v.StartedAt, v.RegisterAt, v.ResignedAt, tt.startedAt, tt.registerAt, tt.endedAt)
				}
			}
			if !found {
				t.Fatalf("%s %s (current=%v) が履歴に無い", tt.name, tt.position, tt.current)
			}
		})
	}
}
