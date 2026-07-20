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
	if got := len(page.underlineYs()); got == 0 {
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
	underlineYs := []float64{318.2}
	if !isUnderlinedText(319.28, 9.8, underlineYs) {
		t.Fatal("text immediately above an underline must be detected")
	}
	if isUnderlinedText(325, 9.8, underlineYs) {
		t.Fatal("text far from an underline must not be detected")
	}
}
