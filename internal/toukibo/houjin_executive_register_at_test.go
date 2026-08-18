package toukibo

import "testing"

func TestGetRegisterAt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "zenkaku spaces", input: "令和　３年　３月３１日登記", want: "令和3年3月31日"},
		{name: "no spaces", input: "令和３年３月３１日登記", want: "令和3年3月31日"},
		{name: "halfwidth numbers", input: "令和3年3月31日登記", want: "令和3年3月31日"},
		{name: "gannen", input: "令和　元年　５月　１日登記", want: "令和元年5月1日"},
		{name: "reiwa 1", input: "令和　１年　９月１０日登記", want: "令和1年9月10日"},
		{name: "heisei", input: "平成３１年　４月２４日登記", want: "平成31年4月24日"},
		{name: "showa", input: "昭和５５年　９月２２日登記", want: "昭和55年9月22日"},
		{name: "taisho", input: "大正　１４年　３月　１日登記", want: "大正14年3月1日"},
		{name: "meiji", input: "明治４５年　７月３０日登記", want: "明治45年7月30日"},
		{name: "appointment is not register", input: "令和　３年　３月３１日就任", wantErr: true},
		{name: "reappointment is not register", input: "令和　２年　６月２５日重任", wantErr: true},
		{name: "resignation is not register", input: "平成２９年　５月２６日辞任", wantErr: true},
		{name: "retirement is not register", input: "令和　３年　６月２１日退任", wantErr: true},
		{name: "death is not register", input: "令和　４年　７月　３日死亡", wantErr: true},
		{name: "iki is not register", input: "平成１２年　２月２８日移記", wantErr: true},
		{name: "correction is not register", input: "令和　３年　６月１０日更正", wantErr: true},
		{name: "outside auditor label is not a date", input: "監査役の登記", wantErr: true},
		{name: "wrapped day fragment is not a date", input: "月１２日登記", wantErr: true},
		{name: "appointment register is not a register date", input: "令和　３年　４月　１日就任の登記", wantErr: true},
		{name: "empty", input: "", wantErr: true},
		{name: "blank", input: "　　　　　　　　　　　　　", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRegisterAt(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("getRegisterAt(%q) = %q, want error", tt.input, got)
				}
				if got != "" {
					t.Fatalf("getRegisterAt(%q) returned %q with error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("getRegisterAt(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("getRegisterAt(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractDates(t *testing.T) {
	tests := []struct {
		name         string
		three        []string
		wantRegister string
		wantResigned string
	}{
		{
			name:         "register only",
			three:        []string{"令和　３年　４月　１日登記"},
			wantRegister: "令和3年4月1日",
		},
		{
			name:  "appointment is ignored",
			three: []string{"平成１７年　６月　１日就任"},
		},
		{
			name:         "appointment then register uses register",
			three:        []string{"平成３１年　３月２８日就任", "平成３１年　４月２４日登記"},
			wantRegister: "平成31年4月24日",
		},
		{
			name:         "reappointment then register uses register",
			three:        []string{"令和　２年　６月２５日重任", "令和　２年　６月３０日登記"},
			wantRegister: "令和2年6月30日",
		},
		{
			name:         "last register wins",
			three:        []string{"令和　２年　６月３０日登記", "令和　３年　６月２９日登記"},
			wantRegister: "令和3年6月29日",
		},
		{
			name:         "resign and register",
			three:        []string{"平成２９年　５月２６日辞任", "平成２９年　５月２６日登記"},
			wantRegister: "平成29年5月26日",
			wantResigned: "平成29年5月26日",
		},
		{
			name:         "retire and register",
			three:        []string{"令和　３年　６月２１日退任", "令和　３年　６月２９日登記"},
			wantRegister: "令和3年6月29日",
			wantResigned: "令和3年6月21日",
		},
		{
			name:         "death and register",
			three:        []string{"令和　４年　７月　３日死亡", "令和　４年　７月１１日登記"},
			wantRegister: "令和4年7月11日",
			wantResigned: "令和4年7月3日",
		},
		{
			name:         "dismiss and register",
			three:        []string{"令和　２年　１月１０日解任", "令和　２年　１月１５日登記"},
			wantRegister: "令和2年1月15日",
			wantResigned: "令和2年1月10日",
		},
		{
			name:         "leave and register",
			three:        []string{"平成３０年　４月　１日退社", "平成３０年　４月　５日登記"},
			wantRegister: "平成30年4月5日",
			wantResigned: "平成30年4月1日",
		},
		{
			name:         "erase and register",
			three:        []string{"令和　１年　８月　１日抹消", "令和　１年　８月　２日登記"},
			wantRegister: "令和1年8月2日",
			wantResigned: "令和1年8月1日",
		},
		{
			name:         "abolish and register",
			three:        []string{"令和　５年　３月　１日廃止", "令和　５年　３月　２日登記"},
			wantRegister: "令和5年3月2日",
			wantResigned: "令和5年3月1日",
		},
		{
			name:         "responsibility change on one line",
			three:        []string{"令和　３年　５月１９日責任変更", "令和　３年　５月２０日登記"},
			wantRegister: "令和3年5月20日",
			wantResigned: "令和3年5月19日",
		},
		{
			name:         "responsibility change split across lines is not a resigned date",
			three:        []string{"令和　３年　５月１９日責任", "変更", "令和　３年　５月２０日登記"},
			wantRegister: "令和3年5月20日",
		},
		{
			name:         "outside auditor wraps across lines",
			three:        []string{"令和　４年１０月１３日社外", "監査役の登記"},
			wantRegister: "令和4年10月13日",
		},
		{
			name:  "wrapped day fragment alone is ignored",
			three: []string{"月１２日登記"},
		},
		{
			name:  "iki is ignored",
			three: []string{"平成１２年　２月２８日移記"},
		},
		{
			name:  "empty columns",
			three: []string{"", "　　　　　　　　　　　　　", "－－－－－－－－－－－－－"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerAt, resignedAt := extractDates(tt.three)
			if registerAt != tt.wantRegister || resignedAt != tt.wantResigned {
				t.Fatalf("extractDates(%q) = (%q, %q), want (%q, %q)",
					tt.three, registerAt, resignedAt, tt.wantRegister, tt.wantResigned)
			}
		})
	}
}

func TestApplyDatesToExecutives(t *testing.T) {
	t.Run("applies the same dates to every executive in the block", func(t *testing.T) {
		evs := []HoujinExecutiveValue{
			{Name: "山田太郎", Position: "取締役", IsValid: true},
			{Name: "鈴木花子", Position: "取締役", IsValid: true},
		}
		applyDatesToExecutives(evs, []string{"令和　３年　３月３１日就任", "令和　３年　４月　１日登記"})

		for i, ev := range evs {
			if ev.RegisterAt != "令和3年4月1日" {
				t.Fatalf("evs[%d].RegisterAt = %q, want 令和３年４月１日", i, ev.RegisterAt)
			}
			if ev.ResignedAt != "" {
				t.Fatalf("evs[%d].ResignedAt = %q, want empty", i, ev.ResignedAt)
			}
			if !ev.IsValid {
				t.Fatalf("evs[%d] must remain valid", i)
			}
		}
	})

	t.Run("resignation invalidates every executive in the block", func(t *testing.T) {
		evs := []HoujinExecutiveValue{
			{Name: "山田太郎", Position: "代表取締役", IsValid: true},
		}
		applyDatesToExecutives(evs, []string{"平成２９年　５月２６日辞任", "平成２９年　５月２６日登記"})

		if evs[0].RegisterAt != "平成29年5月26日" {
			t.Fatalf("RegisterAt = %q, want 平成２９年５月２６日", evs[0].RegisterAt)
		}
		if evs[0].ResignedAt != "平成29年5月26日" {
			t.Fatalf("ResignedAt = %q, want 平成２９年５月２６日", evs[0].ResignedAt)
		}
		if evs[0].IsValid {
			t.Fatal("resigned executive must be invalid")
		}
	})

	t.Run("empty date column leaves fields unchanged", func(t *testing.T) {
		evs := []HoujinExecutiveValue{
			{Name: "山田太郎", Position: "取締役", IsValid: true},
		}
		applyDatesToExecutives(evs, []string{"平成１７年　６月　１日就任", "　　　　　　　　　　　　　"})

		if evs[0].RegisterAt != "" || evs[0].ResignedAt != "" || !evs[0].IsValid {
			t.Fatalf("unchanged executive: got %+v", evs[0])
		}
	})

	t.Run("empty executive list is a no-op", func(t *testing.T) {
		applyDatesToExecutives(nil, []string{"令和　３年　４月　１日登記"})
	})
}

type wantExecutive struct {
	Name       string
	Position   string
	RegisterAt string
	ResignedAt string
	IsValid    bool
}

func assertExecutives(t *testing.T, got HoujinExecutiveValueArray, want []wantExecutive) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("executive count: want %d, got %d (%+v)", len(want), len(got), got)
	}
	for i := range want {
		g, w := got[i], want[i]
		if g.Name != w.Name || g.Position != w.Position || g.RegisterAt != w.RegisterAt || g.ResignedAt != w.ResignedAt || g.IsValid != w.IsValid {
			t.Fatalf("executive[%d]: got {Name:%q Position:%q RegisterAt:%q ResignedAt:%q IsValid:%v}, want %+v",
				i, g.Name, g.Position, g.RegisterAt, g.ResignedAt, g.IsValid, w)
		}
	}
}

func TestGetHoujinExecutiveValueRegisterAt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantExecutive
	}{
		{
			name: "no right column date stays empty",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　　　　　　　　　　　　　　　┃" +
				revert2 +
				"┃　　　　　　　　│　代表取締役　　　山　田　太　郎　　　　　　　　　　　　　　　　　　　　　┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", IsValid: true},
				{Name: "山田太郎", Position: "代表取締役", IsValid: true},
			},
		},
		{
			name:  "appointment without register stays empty",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│平成１７年　６月　１日就任┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", IsValid: true},
			},
		},
		{
			name: "appointment then empty register column stays empty",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│平成１７年　６月　１日就任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│　　　　　　　　　　　　　┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", IsValid: true},
			},
		},
		{
			name: "appointment and register uses register not appointment",
			input: "┃　　　　　　　　│　取締役　　　　　鈴　木　花　子　　　　　　　│平成３１年　３月２８日就任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│平成３１年　４月２４日登記┃",
			want: []wantExecutive{
				{Name: "鈴木花子", Position: "取締役", RegisterAt: "平成31年4月24日", IsValid: true},
			},
		},
		{
			name: "reappointment and register uses register not reappointment",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　２年　６月２５日重任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　２年　６月３０日登記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", RegisterAt: "令和2年6月30日", IsValid: true},
			},
		},
		{
			name: "resignation and register invalidates the executive",
			input: "┃　　　　　　　　│　代表取締役　　　山　田　太　郎　　　　　　　│平成１７年　６月　１日就任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│平成２９年　５月２６日辞任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│平成２９年　５月２６日登記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "代表取締役", RegisterAt: "平成29年5月26日", ResignedAt: "平成29年5月26日"},
			},
		},
		{
			name: "gannen register date",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　元年　５月　１日就任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　元年　５月１０日登記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", RegisterAt: "令和元年5月10日", IsValid: true},
			},
		},
		{
			name: "reiwa 1 register date",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　１年　９月　１日就任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　１年　９月１０日登記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", RegisterAt: "令和1年9月10日", IsValid: true},
			},
		},
		{
			name: "same register date is applied to every name in the block",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　３年　４月　１日登記┃" +
				"┃　　　　　　　　│　取締役　　　　　鈴　木　花　子　　　　　　　　　　　　　　　　　　　　　┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", RegisterAt: "令和3年4月1日", IsValid: true},
				{Name: "鈴木花子", Position: "取締役", RegisterAt: "令和3年4月1日", IsValid: true},
			},
		},
		{
			name: "audit-and-supervisory committee wraps and still takes register date",
			input: "┃　　　　　　　　│　取締役・監査等　佐　藤　り　か　　　　　　　│令和　３年　６月１８日就任┃" +
				"┃　　　　　　　　│　委員　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　３年　６月２１日登記┃",
			want: []wantExecutive{
				{Name: "佐藤りか", Position: "取締役・監査等委員", RegisterAt: "令和3年6月21日", IsValid: true},
			},
		},
		{
			name:  "iki is not treated as register date",
			input: "┃　　　　　　　　│　理事長　　　　　山　田　太　郎　　　　　　　│平成１２年　２月２８日移記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "理事長", IsValid: true},
			},
		},
		{
			name: "underlined executive can still carry a register date",
			input: "┃　　　　　　　　│" + deletedTextMarker + "　取締役　　　　　山　田　太　郎　　　　　　　│令和　３年　４月　１日登記┃" +
				revert2 +
				"┃　　　　　　　　│　取締役　　　　　鈴　木　花　子　　　　　　　│令和　４年　４月　１日登記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", RegisterAt: "令和3年4月1日"},
				{Name: "鈴木花子", Position: "取締役", RegisterAt: "令和4年4月1日", IsValid: true},
			},
		},
		{
			name: "last register in the same block wins",
			input: "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　２年　６月３０日登記┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　３年　６月２９日登記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "取締役", RegisterAt: "令和3年6月29日", IsValid: true},
			},
		},
		{
			name: "outside auditor register wraps across the right column",
			input: "┃　　　　　　　　│　監査役　　　　　山　田　太　郎　　　　　　　│令和　４年１０月１３日社外┃" +
				"┃　　　　　　　　│　（社外監査役）　　　　　　　　　　　　　　　│監査役の登記　　　　　　　┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "監査役", RegisterAt: "令和4年10月13日", IsValid: true},
			},
		},
		{
			name:  "auditor-of-register label is not treated as a date",
			input: "┃　　　　　　　　│　監査役　　　　　山　田　太　郎　　　　　　　│監査役の登記　　　　　　　┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "監査役", IsValid: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHoujinExecutiveValue(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assertExecutives(t, got, tt.want)
		})
	}
}

func TestGetHoujinExecutiveValueRegisterAtDateOnlyAppliesToPrevious(t *testing.T) {
	t.Run("resignation register on a following date-only row", func(t *testing.T) {
		s := "┃　　　　　　　　│　代表取締役　　　山　田　太　郎　　　　　　　│平成１７年　６月　１日就任┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│　　　　　　　　　　　　　┃" +
			revert4 +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│平成２９年　５月２６日辞任┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│平成２９年　５月２６日登記┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "代表取締役", RegisterAt: "平成29年5月26日", ResignedAt: "平成29年5月26日"},
		})
	})

	t.Run("retirement register on a following date-only row", func(t *testing.T) {
		s := "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　２年　６月２５日重任┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　２年　６月３０日登記┃" +
			revert4 +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　３年　６月２１日退任┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　３年　６月２９日登記┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "取締役", RegisterAt: "令和3年6月29日", ResignedAt: "令和3年6月21日"},
		})
	})

	t.Run("death register on a following date-only row", func(t *testing.T) {
		s := "┃　　　　　　　　│　無限責任社員　　山　田　太　郎　　　　　　　　　　　　　　　　　　　　　┃" +
			revert5 +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　４年　７月　３日死亡┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　４年　７月１１日登記┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "無限責任社員", RegisterAt: "令和4年7月11日", ResignedAt: "令和4年7月3日"},
		})
	})

	t.Run("date-only row without register does not rewrite previous", func(t *testing.T) {
		s := "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　２年　６月３０日登記┃" +
			revert4 +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　３年　６月２１日退任┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "取締役", RegisterAt: "令和2年6月30日", IsValid: true},
		})
	})
}

func TestGetHoujinExecutiveValueInheritsRegisterAtAcrossCorrection(t *testing.T) {
	t.Run("address change then correction keeps previous register", func(t *testing.T) {
		s := "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　２年　３月１２日登記┃" +
			revert2 +
			"┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　６年　１月　１日変更┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　６年　１月　９日修正┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "取締役", RegisterAt: "令和2年3月12日"},
			{Name: "山田太郎", Position: "取締役", RegisterAt: "令和2年3月12日", IsValid: true},
		})
	})

	t.Run("correction without register keeps previous register", func(t *testing.T) {
		s := "┃　　　　　　　　│　代表取締役　　　山　田　太　郎　　　　　　　│令和　５年　６月　６日登記┃" +
			revert2 +
			"┃　　　　　　　　│　代表取締役　　　山　田　太　郎　　　　　　　│令和　６年　６月１０日更正┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "代表取締役", RegisterAt: "令和5年6月6日"},
			{Name: "山田太郎", Position: "代表取締役", RegisterAt: "令和5年6月6日", IsValid: true},
		})
	})

	t.Run("does not inherit across different positions", func(t *testing.T) {
		s := "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　　　　　　　　　　　　　　　┃" +
			revert2 +
			"┃　　　　　　　　│　代表取締役　　　山　田　太　郎　　　　　　　│令和　３年　４月　１日登記┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "取締役", IsValid: true},
			{Name: "山田太郎", Position: "代表取締役", RegisterAt: "令和3年4月1日", IsValid: true},
		})
	})
}
