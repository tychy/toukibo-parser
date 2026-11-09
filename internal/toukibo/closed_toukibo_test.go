package toukibo

import "testing"

// 閉鎖された登記簿の「登記記録に関する事項」の欄。閉鎖会社は「解　散」の欄を持たず、
// 解散が登記記録の欄にだけ「…に合併し解散」と書かれる。行末は全角スペースで埋められている。
const closedToukiRecordPart = " ┃登記記録に関する│　令和７年１０月１日東京都台東区元浅草三丁目７番１号から本店移転　　　　　┃" +
	" ┃事項　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　　令和　７年１０月２１日登記┃" +
	" ┃　　　　　　　　├─────────────────────────────────────┨" +
	" ┃　　　　　　　　│　令和８年８月１日東京都港区麻布台一丁目３番１号サンプル株式会社に合併し　┃" +
	" ┃　　　　　　　　│　解散　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　┃" +
	" ┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　令和　８年　８月　３日登記┃" +
	" ┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　令和　８年　８月　３日閉鎖┃ "

func TestConsumeHoujinClosedRecord(t *testing.T) {
	h := &HoujinBody{}
	h.ConsumeHoujinClosedRecord(closedToukiRecordPart)

	if h.HoujinClosedAt != "令和8年8月3日" {
		t.Fatalf("closed at: want 令和8年8月3日, got %q", h.HoujinClosedAt)
	}
	// 直前の「本店移転」「登記」の日付ではなく、合併し解散の直前の日付を採る
	if h.HoujinDissolvedAt != "令和8年8月1日" {
		t.Fatalf("dissolved at: want 令和8年8月1日, got %q", h.HoujinDissolvedAt)
	}
}

func TestConsumeHoujinClosedRecordKeepsDissolvedAt(t *testing.T) {
	h := &HoujinBody{HoujinDissolvedAt: "令和8年7月1日"}
	h.ConsumeHoujinClosedRecord(closedToukiRecordPart)

	if h.HoujinDissolvedAt != "令和8年7月1日" {
		t.Fatalf("dissolved at was overwritten: %q", h.HoujinDissolvedAt)
	}
}

// 管轄変更で閉鎖された登記簿（解散の記載が無い）は解散日を埋めない。
func TestConsumeHoujinClosedRecordWithoutDissolution(t *testing.T) {
	s := " ┃登記記録に関する│　令和７年１０月１日東京都台東区元浅草三丁目７番１号から本店移転　　　　　┃" +
		" ┃事項　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　　令和　７年１０月２１日登記┃" +
		" ┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　令和　７年１０月２１日閉鎖┃ "

	h := &HoujinBody{}
	h.ConsumeHoujinClosedRecord(s)

	if h.HoujinClosedAt != "令和7年10月21日" {
		t.Fatalf("closed at: want 令和7年10月21日, got %q", h.HoujinClosedAt)
	}
	if h.HoujinDissolvedAt != "" {
		t.Fatalf("dissolved at should be empty, got %q", h.HoujinDissolvedAt)
	}
}
