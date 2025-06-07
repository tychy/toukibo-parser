package toukibo

import (
	"fmt"
	"regexp"
	"strings"
)

// DateExtractor は日付抽出のための共通関数
type DateExtractor struct {
	suffixes []string
	pattern  *regexp.Regexp
}

// NewDateExtractor は新しいDateExtractorを作成する
func NewDateExtractor(suffixes []string) *DateExtractor {
	pattern := fmt.Sprintf(`([%s]+)　*(%s)`, ZenkakuStringPattern, strings.Join(suffixes, "|"))
	return &DateExtractor{
		suffixes: suffixes,
		pattern:  regexp.MustCompile(pattern),
	}
}

// Extract は与えられたテキストから日付を抽出する
func (de *DateExtractor) Extract(text string) (string, bool) {
	matches := de.pattern.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[1], true
	}
	return "", false
}

// 共通の日付抽出関数
var (
	// 登記日抽出
	RegisterDateExtractor = NewDateExtractor([]string{"登記"})
	// 退任日抽出
	ResignDateExtractor = NewDateExtractor([]string{"退任", "辞任"})
	// 就任日抽出
	AppointDateExtractor = NewDateExtractor([]string{"就任"})
	// 解散日抽出
	DissolvedDateExtractor = NewDateExtractor([]string{"解散"})
	// 破産日抽出
	BankruptDateExtractor = NewDateExtractor([]string{"破産手続開始"})
	// 継続日抽出
	ContinuedDateExtractor = NewDateExtractor([]string{"継続"})
)

// ExtractDateWithSuffix は指定されたサフィックスで日付を抽出する汎用関数
func ExtractDateWithSuffix(text string, suffixes []string) (string, bool) {
	extractor := NewDateExtractor(suffixes)
	return extractor.Extract(text)
}