package toukibo

import "testing"

func TestConsumeHoujinDissolvedAtUsesLatestDate(t *testing.T) {
	h := &HoujinBody{}
	s := "┃解　散　　　　│　令和3年1月1日株主総会の決議により解散┃" +
		"┃解　散　　　　│　令和5年12月15日株主総会の決議により解散┃"
	if !h.ConsumeHoujinDissolvedAt(s) {
		t.Fatal("dissolution was not consumed")
	}
	if h.HoujinDissolvedAt != "令和5年12月15日" {
		t.Fatalf("latest dissolution: want 令和5年12月15日, got %q", h.HoujinDissolvedAt)
	}
}

func TestConsumeHoujinContinuedAtUsesLatestDate(t *testing.T) {
	h := &HoujinBody{}
	s := "┃会社継続　　　│　令和4年1月22日会社継続┃" +
		"┃会社継続　　　│　令和6年2月3日会社継続┃"
	if !h.ConsumeHoujinContinuedAt(s) {
		t.Fatal("continuation was not consumed")
	}
	if h.HoujinContinuedAt != "令和6年2月3日" {
		t.Fatalf("latest continuation: want 令和6年2月3日, got %q", h.HoujinContinuedAt)
	}
}

func TestGetHoujinExecutiveValueDeletedUnderline(t *testing.T) {
	s := "┃　　　　　　　　│" + deletedTextMarker + "　取締役　　　　　山　田　太　郎　　　　　　　　　　　　　　　　　　　　　┃" +
		revert2 +
		"┃　　　　　　　　│　取締役　　　　　鈴　木　花　子　　　　　　　│令和　３年　３月３１日就任┃"

	executives, err := GetHoujinExecutiveValue(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(executives) != 2 {
		t.Fatalf("executive count: want 2, got %d", len(executives))
	}
	if executives[0].IsValid {
		t.Fatal("underlined executive must be invalid")
	}
	if executives[0].Name != "山田太郎" || executives[0].Position != "取締役" {
		t.Fatalf("deleted executive: got %+v", executives[0])
	}
	if !executives[1].IsValid {
		t.Fatal("current executive must remain valid")
	}
	if executives[1].Name != "鈴木花子" || executives[1].Position != "取締役" {
		t.Fatalf("current executive: got %+v", executives[1])
	}
}

func TestGetHoujinExecutiveValueQualificationChange(t *testing.T) {
	s := "┃　　　　　　　　│　社員　　　　　　山　田　太　郎　　　　　　　　　　　　　　　　　　　　　┃" +
		revert2 +
		"┃　　　　　　　　│　社員　　　　　　鈴　木　花　子　　　　　　　│令和　３年　３月３１日資格┃" +
		"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│変更　　　　　　　　　　　┃" +
		"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨" +
		"┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　│令和　３年　３月３１日登記┃"

	executives, err := GetHoujinExecutiveValue(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(executives) != 2 {
		t.Fatalf("executive count: want 2, got %d", len(executives))
	}

	if executives[0].Name != "山田太郎" || executives[0].Position != "社員" {
		t.Fatalf("previous executive: got %+v", executives[0])
	}
	if executives[0].IsValid {
		t.Fatal("previous executive must be invalid after qualification change")
	}

	if executives[1].Name != "鈴木花子" || executives[1].Position != "社員" {
		t.Fatalf("current executive: got %+v", executives[1])
	}
	if !executives[1].IsValid {
		t.Fatal("current executive must remain valid after qualification change")
	}
	if executives[1].RegisterAt != "令和3年3月31日" {
		t.Fatalf("current executive registerAt: want 令和3年3月31日, got %q", executives[1].RegisterAt)
	}
}
