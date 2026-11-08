package toukibo

import (
	"strings"
	"testing"
)

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
			registerAt, resigned, _ := extractDates(tt.three)
			if registerAt != tt.wantRegister || resigned.At != tt.wantResigned {
				t.Fatalf("extractDates(%q) = (%q, %q), want (%q, %q)",
					tt.three, registerAt, resigned.At, tt.wantRegister, tt.wantResigned)
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
		{
			name: "outside auditor label without a name is not another person",
			input: "┃　　　　　　　　│　監査役　　　　　山　田　太　郎　　　　　　　│令和　５年　６月２３日重任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　社外監査役　　　　　　　　　　　　　　　　　│令和　５年　６月２９日登記┃",
			want: []wantExecutive{
				{Name: "山田太郎", Position: "監査役", RegisterAt: "令和5年6月29日", IsValid: true},
			},
		},
		{
			name: "outside director label without a name is not another person",
			input: "┃　　　　　　　　│　取締役　　　　　鈴　木　花　子　　　　　　　│令和　７年　６月２７日就任┃" +
				"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
				"┃　　　　　　　　│　社外取締役　　　　　　　　　　　　　　　　　│令和　７年　７月　２日登記┃",
			want: []wantExecutive{
				{Name: "鈴木花子", Position: "取締役", RegisterAt: "令和7年7月2日", IsValid: true},
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

	t.Run("corrected resignation date applies to previous", func(t *testing.T) {
		// 令和８年４月１日辞任で登記した後、退任日を３月３１日に更正した記載。
		// 元の辞任は下線（抹消）付きで、更正の記載には「登記」が無い。
		deleted := func(s string) string {
			var b strings.Builder
			for _, r := range s {
				b.WriteString(deletedTextMarker)
				b.WriteRune(r)
			}
			return b.String()
		}
		s := "┃　　　　　　　　│　取締役　　　　　山　田　太　郎　　　　　　　│令和　７年　５月２２日重任┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　７年　６月　９日登記┃" +
			revert4 +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│" + deleted("令和　８年　４月　１日辞任") + "┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　８年　４月　８日登記┃" +
			revert4 +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　８年　３月３１日辞任┃" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
			"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　８年　４月１７日更正┃"

		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			// 下線付き（取り消し）の辞任の登記日（令和８年４月８日）では上書きしない
			{Name: "山田太郎", Position: "取締役", RegisterAt: "令和7年6月9日", ResignedAt: "令和8年3月31日"},
		})
		if got[0].ResignedReason != "辞任" {
			t.Fatalf("ResignedReason = %q, want 辞任", got[0].ResignedReason)
		}
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

func TestGetResignedAtReason(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantAt     string
		wantReason string
	}{
		{name: "resign", input: "令和　３年　６月２２日辞任", wantAt: "令和3年6月22日", wantReason: "辞任"},
		{name: "retire", input: "令和　３年　６月２２日退任", wantAt: "令和3年6月22日", wantReason: "退任"},
		{name: "death", input: "令和　４年　７月　３日死亡", wantAt: "令和4年7月3日", wantReason: "死亡"},
		{name: "dismiss", input: "令和　２年　１月１０日解任", wantAt: "令和2年1月10日", wantReason: "解任"},
		{name: "erase", input: "令和　１年　８月　１日抹消", wantAt: "令和1年8月1日", wantReason: "抹消"},
		{name: "abolish", input: "令和　５年　３月　１日廃止", wantAt: "令和5年3月1日", wantReason: "廃止"},
		{name: "leave", input: "平成３０年　４月　１日退社", wantAt: "平成30年4月1日", wantReason: "退社"},
		{name: "responsibility change", input: "令和　３年　５月１９日責任変更", wantAt: "令和3年5月19日", wantReason: "責任変更"},
		// 「職務執行者辞任」のように日付を伴わない記載は、理由だけ返して退任日は空にする
		{name: "no date before reason", input: "職務執行者辞任", wantReason: "辞任"},
		{name: "date fragment before reason", input: "令和　３年　６月辞任", wantReason: "辞任"},
		{name: "definition clause before abolish", input: "監査役の監査の範囲を会計に関するものに限定する旨の定款の定めの廃止", wantReason: "廃止"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := getResignedAt(tt.input)
			if !ok {
				t.Fatalf("getResignedAt(%q) = not found, want (%q, %q)", tt.input, tt.wantAt, tt.wantReason)
			}
			if got.At != tt.wantAt || got.Reason != tt.wantReason {
				t.Fatalf("getResignedAt(%q) = (%q, %q), want (%q, %q)", tt.input, got.At, got.Reason, tt.wantAt, tt.wantReason)
			}
		})
	}
}

// 辞任の登記が抹消され、「◯日抹消により回復」で役員の地位が戻った記載。
// 回復は退任ではないので、回復の行に退任を付けない。
func TestGetHoujinExecutiveValueRestoredByErasure(t *testing.T) {
	const (
		blank = "　　　　　　　　　　　　　　　　　　　　　　　"
		sep   = "┃　　　　　　　　│" + blank + "├－－－－－－－－－－－－－┨"
	)
	row := func(mid, right string) string {
		return "┃　　　　　　　　│" + mid + "│" + right + "┃"
	}
	deleted := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			b.WriteString(deletedTextMarker)
			b.WriteRune(r)
		}
		return b.String()
	}
	name := "　取締役　　　　　山　田　太　郎　　　　　　　"
	erased := row(name, "令和　１年１１月３０日就任") + sep +
		row(blank, "令和　１年１２月２３日登記") + revert4 +
		row(blank, deleted("令和　５年　５月　１日辞任")) + sep +
		row(blank, "令和　５年　８月１７日登記") + revert4 +
		row(blank, "辞任の登記　　　　　　　　") + sep +
		row(blank, "令和　５年　９月１１日抹消") + revert4 +
		row(name, "　　　　　　　　　　　　　") + sep +
		row(blank, "令和　５年　９月１１日抹消") +
		row(blank, "により回復　　　　　　　　")

	t.Run("restored row is current", func(t *testing.T) {
		got, err := GetHoujinExecutiveValue(erased)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Fatalf("executive count: want 2, got %d (%+v)", len(got), got)
		}
		last := got[len(got)-1]
		if !last.IsValid || last.ResignedAt != "" || last.ResignedReason != "" {
			t.Fatalf("restored row must be current without resignation: %+v", last)
		}
		if got[0].IsValid || got[0].ResignedReason != "" {
			t.Fatalf("erased row must be replaced without resignation: %+v", got[0])
		}
	})

	t.Run("reappointment after restoration replaces restored row", func(t *testing.T) {
		s := erased + revert4 +
			row(name, "令和　７年　５月１５日重任") + sep +
			row(blank, "令和　７年　５月２６日登記")
		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 3 {
			t.Fatalf("executive count: want 3, got %d (%+v)", len(got), got)
		}
		for i, v := range got {
			if v.ResignedAt != "" || v.ResignedReason != "" {
				t.Fatalf("executive[%d] must not have resignation: %+v", i, v)
			}
			if v.IsValid != (i == 2) {
				t.Fatalf("executive[%d].IsValid = %v", i, v.IsValid)
			}
		}
		if got[2].RegisterAt != "令和7年5月26日" {
			t.Fatalf("RegisterAt = %q, want 令和7年5月26日", got[2].RegisterAt)
		}
	})
}

// 監査役会の定めの廃止に伴い、社外監査役の注記が外れて登記が書き直された記載。
// 「廃止により変更」は退任ではないので、変更の行が現任になる。
func TestGetHoujinExecutiveValueChangedByAbolition(t *testing.T) {
	const (
		blank = "　　　　　　　　　　　　　　　　　　　　　　　"
		sep   = "┃　　　　　　　　│" + blank + "├－－－－－－－－－－－－－┨"
	)
	row := func(mid, right string) string {
		return "┃　　　　　　　　│" + mid + "│" + right + "┃"
	}
	name := "　監査役　　　　　山　田　太　郎　　　　　　　"
	outside := "　（社外監査役）　　　　　　　　　　　　　　　"
	changed := row(name, "令和　５年　６月１９日就任") + sep +
		row(outside, "令和　５年　６月１９日登記") + revert4 +
		row(name, "令和　６年　３月２１日監査") +
		row(blank, "役会の定めの廃止により変更") + sep +
		row(blank, "令和　６年　３月２９日登記")

	t.Run("changed row is current", func(t *testing.T) {
		got, err := GetHoujinExecutiveValue(changed)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "監査役", RegisterAt: "令和5年6月19日"},
			{Name: "山田太郎", Position: "監査役", RegisterAt: "令和6年3月29日", IsValid: true},
		})
		for i, v := range got {
			if v.ResignedReason != "" {
				t.Fatalf("executive[%d] must not have resignation: %+v", i, v)
			}
		}
	})

	t.Run("reappointment after change replaces changed row", func(t *testing.T) {
		s := changed + revert4 +
			row(name, "令和　７年　６月２０日重任") + sep +
			row(blank, "令和　７年　６月２７日登記")
		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		assertExecutives(t, got, []wantExecutive{
			{Name: "山田太郎", Position: "監査役", RegisterAt: "令和5年6月19日"},
			{Name: "山田太郎", Position: "監査役", RegisterAt: "令和6年3月29日"},
			{Name: "山田太郎", Position: "監査役", RegisterAt: "令和7年6月27日", IsValid: true},
		})
	})
}

// 右の列で折り返され、日付と事由が別の行にある退任の記載。
func TestExtractDatesWrappedResignation(t *testing.T) {
	tests := []struct {
		name       string
		three      []string
		wantAt     string
		wantReason string
	}{
		{
			name:       "executor resigned",
			three:      []string{"令和　７年　３月３１日職務", "執行者辞任　　　　　　　　", "", "令和　７年　４月　３日登記"},
			wantAt:     "令和7年3月31日",
			wantReason: "辞任",
		},
		{
			name:       "representative erased when one director remains",
			three:      []string{"令和　６年　６月３０日取締", "役が１名となったため抹消　", "", "令和　６年　７月２２日登記"},
			wantAt:     "令和6年6月30日",
			wantReason: "抹消",
		},
		{
			name:       "three lines",
			three:      []string{"令和　６年１０月３０日会社", "を代表しない取締役の不存在", "により抹消　　　　　　　　", "令和　６年１１月　５日登記"},
			wantAt:     "令和6年10月30日",
			wantReason: "抹消",
		},
		{
			// 日付と事由の間に登記がある場合は、別の記載の日付なので採用しない
			name:       "register between date and reason",
			three:      []string{"令和　２年　３月２３日登記", "職務執行者辞任　　　　　　"},
			wantAt:     "",
			wantReason: "辞任",
		},
		{
			// 下線（抹消）の印は文字列の塊の先頭に付く。日付の行だけに下線がある場合も採用しない
			name:       "deleted date line",
			three:      []string{deletedTextMarker + "令和　７年　３月３１日職務", "執行者辞任"},
			wantAt:     "",
			wantReason: "辞任",
		},
		{
			// 事由の行の判定（getResignedAt）は従来のまま。補完で抹消済みの日付を付けないことだけを見る
			name:       "deleted date and reason lines",
			three:      []string{deletedTextMarker + "令和　７年　３月３１日職務", deletedTextMarker + "執行者辞任"},
			wantAt:     "",
			wantReason: "辞任",
		},
		{
			// 日付と事由の間に別の出来事（回復・変更・就任）がある場合は、その出来事の日付なので採用しない
			name:       "restoration between date and reason",
			three:      []string{"令和　５年　９月１１日抹消", "により回復", "職務執行者辞任"},
			wantAt:     "",
			wantReason: "辞任",
		},
		{
			name:       "change between date and reason",
			three:      []string{"令和　６年　３月２１日監査役会", "の定めの廃止により変更", "職務執行者辞任"},
			wantAt:     "",
			wantReason: "辞任",
		},
		{
			name:       "appointment between date and reason",
			three:      []string{"令和　７年　３月３１日就任", "職務執行者辞任"},
			wantAt:     "",
			wantReason: "辞任",
		},
		{
			// 1行に収まっていても、日付と事由の間に語句があれば同じように取る
			name:       "executor resigned in one line",
			three:      []string{"令和　７年　３月３１日職務執行者辞任", "令和　７年　４月　３日登記"},
			wantAt:     "令和7年3月31日",
			wantReason: "辞任",
		},
		{
			// 1行に収まる記載は従来どおり
			name:       "single line",
			three:      []string{"令和　６年　６月２５日退任", "令和　６年　６月２７日登記"},
			wantAt:     "令和6年6月25日",
			wantReason: "退任",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, found := extractDates(tt.three)
			if found != (tt.wantReason != "") {
				t.Fatalf("extractDates(%q): found = %v", tt.three, found)
			}
			if got.At != tt.wantAt || got.Reason != tt.wantReason {
				t.Fatalf("extractDates(%q) = (%q, %q), want (%q, %q)", tt.three, got.At, got.Reason, tt.wantAt, tt.wantReason)
			}
		})
	}
}

func TestStartedDateOf(t *testing.T) {
	tests := []struct {
		name  string
		three []string
		want  string
	}{
		{"appointed", []string{"令和　５年　６月１９日就任", "令和　５年　６月１９日登記"}, "令和5年6月19日"},
		{"reappointed", []string{"令和　５年　６月２７日重任", "令和　５年　６月３０日登記"}, "令和5年6月27日"},
		{"executor wrapped", []string{"令和　６年１１月１４日職務", "執行者就任", "令和　６年１１月１４日登記"}, "令和6年11月14日"},
		{"joined", []string{"令和　６年　４月１６日加入", "令和　６年　４月１７日登記"}, "令和6年4月16日"},
		{"qualification change", []string{"令和　３年　４月　１日資格", "変更", "令和　３年　４月　５日登記"}, "令和3年4月1日"},
		{"representative right granted", []string{"令和　３年　７月　１日代表", "権付与", "令和　３年　７月　１日登記"}, "令和3年7月1日"},
		{"entered", []string{"令和　２年　４月　１日入社", "令和　２年　４月　３日登記"}, "令和2年4月1日"},
		{"appointed by court", []string{"令和　２年　４月　１日東京地方", "裁判所の選任", "令和　２年　４月　３日登記"}, "令和2年4月1日"},
		// 重任などで置き換わった行は下線付きだが、始まりの日付は読む
		{"superseded row", []string{deletedTextMarker + "令和　３年　５月１６日重任", "令和　３年　５月１７日登記"}, "令和3年5月16日"},
		// 住所移転・氏の変更などの行には始まりの日付が無い（前の行から引き継ぐ）
		{"address change", []string{"令和　６年１０月２３日住所移転", "令和　６年１１月　８日登記"}, ""},
		{"no date", []string{"", ""}, ""},
		// 日付と始まりの語の間に別の出来事がある日付は採用しない
		{"other event between", []string{"令和　６年　３月２１日住所移転", "就任"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startedDateOf(tt.three); got != tt.want {
				t.Fatalf("startedDateOf(%q) = %q, want %q", tt.three, got, tt.want)
			}
		})
	}
}

// 始まりの日付の引き継ぎ: 住所移転の行は前の行から引き継ぎ、重任の登記が抹消されて回復した行は
// 抹消された重任ではなく、その前の行から引き継ぐ。
func TestGetHoujinExecutiveValueStartedAt(t *testing.T) {
	const (
		blank = "　　　　　　　　　　　　　　　　　　　　　　　"
		sep   = "┃　　　　　　　　│" + blank + "├－－－－－－－－－－－－－┨"
	)
	row := func(mid, right string) string {
		return "┃　　　　　　　　│" + mid + "│" + right + "┃"
	}
	name := "　取締役　　　　　山　田　太　郎　　　　　　　"

	t.Run("address change inherits", func(t *testing.T) {
		s := row(name, "令和　２年　６月２５日就任") + sep +
			row(blank, "令和　２年　６月３０日登記") + revert4 +
			row(name, "令和　６年１０月２３日住所") +
			row(blank, "移転　　　　　　　　　　　") + sep +
			row(blank, "令和　６年１１月　８日登記")
		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		last := got[len(got)-1]
		if !last.IsValid || last.StartedAt != "令和2年6月25日" || last.RegisterAt != "令和6年11月8日" {
			t.Fatalf("got %+v", last)
		}
	})

	t.Run("name change inherits", func(t *testing.T) {
		s := row("　取締役　　　　　山　田　花　子　　　　　　　", "令和　２年　６月２５日就任") + sep +
			row(blank, "令和　２年　６月３０日登記") + revert4 +
			row("　取締役　　　　　佐　藤　花　子　　　　　　　", "令和　６年　１月１０日山田") +
			row(blank, "花子の氏変更　　　　　　　") + sep +
			row(blank, "令和　６年　１月２０日登記")
		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		last := got[len(got)-1]
		if !last.IsValid || last.Name != "佐藤花子" || last.StartedAt != "令和2年6月25日" {
			t.Fatalf("got %+v", last)
		}
	})

	t.Run("restored after erased reappointment", func(t *testing.T) {
		s := row(name, "平成１７年　６月２５日重任") + sep +
			row(blank, "平成１７年　８月１０日登記") + revert4 +
			row(name, "平成２７年　７月２１日重任") + sep +
			row(blank, "平成２７年　７月２４日登記") + revert4 +
			row(blank, "重任の登記　　　　　　　　") + sep +
			row(blank, "令和　７年　８月　６日抹消") + revert4 +
			row(name, "　　　　　　　　　　　　　") + sep +
			row(blank, "令和　７年　８月　６日抹消") +
			row(blank, "により回復　　　　　　　　") + revert4 +
			row(blank, "平成１９年　６月２５日退任") + sep +
			row(blank, "令和　７年　８月　６日登記")
		got, err := GetHoujinExecutiveValue(s)
		if err != nil {
			t.Fatal(err)
		}
		last := got[len(got)-1]
		if last.IsValid || last.ResignedAt != "平成19年6月25日" || last.StartedAt != "平成17年6月25日" {
			t.Fatalf("got %+v", last)
		}
	})
}

func TestHasDeletedEvent(t *testing.T) {
	tests := []struct {
		name  string
		three []string
		want  bool
	}{
		{"deleted resignation", []string{deletedTextMarker + "令和　５年　５月　１日辞任", "令和　５年　８月１７日登記"}, true},
		// 登記日の行だけに下線が付いた記載は、取り消された出来事として扱わない
		{"deleted register only", []string{deletedTextMarker + "令和　３年　５月１７日登記"}, false},
		{"no deletion", []string{"令和　６年　６月２５日退任", "令和　６年　６月２７日登記"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasDeletedEvent(tt.three); got != tt.want {
				t.Fatalf("hasDeletedEvent(%q) = %v, want %v", tt.three, got, tt.want)
			}
		})
	}
}

func TestInheritStartedAtFromPrevious(t *testing.T) {
	t.Run("does not go back across a resignation", func(t *testing.T) {
		evs := HoujinExecutiveValueArray{
			{Name: "山田太郎", Position: "取締役", StartedAt: "令和2年6月25日"},
			{Name: "山田太郎", Position: "取締役", ResignedAt: "令和4年6月25日", ResignedReason: "辞任"},
			{Name: "山田太郎", Position: "取締役"},
		}
		inheritStartedAtFromPrevious(evs, 2)
		if evs[2].StartedAt != "" {
			t.Fatalf("StartedAt = %q, want empty", evs[2].StartedAt)
		}
	})
	t.Run("skips rows without started date", func(t *testing.T) {
		evs := HoujinExecutiveValueArray{
			{Name: "山田太郎", Position: "取締役", StartedAt: "平成17年6月25日"},
			{Name: "山田太郎", Position: "取締役"},
			{Name: "山田太郎", Position: "取締役"},
		}
		inheritStartedAtFromPrevious(evs, 2)
		if evs[2].StartedAt != "平成17年6月25日" {
			t.Fatalf("StartedAt = %q, want 平成17年6月25日", evs[2].StartedAt)
		}
	})
}
