package toukibo

import "testing"

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
}
