package toukibo

import "testing"

func TestGetHoujinFormerExecutives(t *testing.T) {
	body := &HoujinBody{
		HoujinExecutive: []HoujinExecutiveValueArray{
			{
				// 現任
				{Name: "現任太郎", Position: "代表取締役", RegisterAt: "令和3年4月1日", IsValid: true},
				// 重任で置き換わった過去の行（退任日なし）
				{Name: "重任花子", Position: "取締役", RegisterAt: "令和2年4月1日", IsValid: false},
				{Name: "重任花子", Position: "取締役", RegisterAt: "令和3年4月1日", IsValid: true},
			},
			{
				// 辞任した役員
				{Name: "退任次郎", Position: "取締役", RegisterAt: "令和3年5月12日", ResignedAt: "令和3年2月28日", ResignedReason: "辞任", IsValid: false},
				// 死亡した役員
				{Name: "死亡三郎", Position: "監査役", RegisterAt: "令和2年6月1日", ResignedAt: "令和4年7月3日", ResignedReason: "死亡", IsValid: false},
			},
		},
	}

	got := body.GetHoujinFormerExecutives()
	want := []struct {
		name, position, registerAt, resignedAt, resignedReason string
	}{
		{"退任次郎", "取締役", "令和3年5月12日", "令和3年2月28日", "辞任"},
		{"死亡三郎", "監査役", "令和2年6月1日", "令和4年7月3日", "死亡"},
	}

	if len(got) != len(want) {
		t.Fatalf("former executive count: want %d, got %d (%+v)", len(want), len(got), got)
	}
	for i, w := range want {
		g := got[i]
		if g.Name != w.name || g.Position != w.position || g.RegisterAt != w.registerAt ||
			g.ResignedAt != w.resignedAt || g.ResignedReason != w.resignedReason || g.IsValid {
			t.Fatalf("former executive[%d] = %+v, want %+v", i, g, w)
		}
	}
}

func TestGetHoujinFormerExecutivesEmpty(t *testing.T) {
	body := &HoujinBody{}
	if got := body.GetHoujinFormerExecutives(); len(got) != 0 {
		t.Fatalf("want empty, got %+v", got)
	}

	// 現任だけの法人では空になる
	body = &HoujinBody{
		HoujinExecutive: []HoujinExecutiveValueArray{
			{{Name: "現任太郎", Position: "代表取締役", RegisterAt: "令和3年4月1日", IsValid: true}},
		},
	}
	if got := body.GetHoujinFormerExecutives(); len(got) != 0 {
		t.Fatalf("want empty, got %+v", got)
	}
}
