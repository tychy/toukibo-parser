package pdf

import (
	"strings"
	"testing"
)

func TestGetPlainTextMarksUnderlinedText(t *testing.T) {
	r, err := Open("../../testdata/pdf/sample39.pdf")
	if err != nil {
		t.Fatal(err)
	}
	page := r.Page(1)
	if got := len(page.underlines()); got == 0 {
		t.Fatal("expected underlines on sample39")
	}

	plainText, err := page.GetPlainText(nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(plainText, "\u2063") {
		t.Fatal("plain text must not contain deletion markers")
	}

	markedText, err := page.getPlainText(nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(markedText, "\u2063") {
		t.Fatal("underlined text marker was not emitted")
	}
}

func TestIsUnderlinedText(t *testing.T) {
	underlines := []underlineSegment{{minX: 100, maxX: 110, y: 318.2}}
	if !isUnderlinedText(Point{X: 100, Y: 319.28}, Point{X: 110, Y: 319.28}, 9.8, underlines) {
		t.Fatal("text immediately above an underline must be detected")
	}
	if isUnderlinedText(Point{X: 100, Y: 325}, Point{X: 110, Y: 325}, 9.8, underlines) {
		t.Fatal("text far from an underline must not be detected")
	}
	if isUnderlinedText(Point{X: 200, Y: 319.28}, Point{X: 210, Y: 319.28}, 9.8, underlines) {
		t.Fatal("text beside an underline must not be detected")
	}
}
