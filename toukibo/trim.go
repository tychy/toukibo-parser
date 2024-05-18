package toukibo

import (
	"fmt"
	"regexp"
	"strings"
)

// trim all space
func trimAllSpace(s string) string {
	return strings.ReplaceAll(s, "　", "")
}

func trimLeadingTrailingSpace(s string) string {
	// 先頭のスペースを削除
	s = strings.TrimLeft(s, "　")
	s = strings.TrimLeft(s, " ")
	// 最後のスペースを削除
	s = strings.TrimRight(s, "　")
	s = strings.TrimRight(s, " ")
	return s
}

func trimPattern(s, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(s, "")
}

func trimChangeAndRegisterAt(s string) (string, string, string) {
	// trim ┃　　　　　　　　│　　　　　　　　　平成３０年　７月３１日変更　　平成３０年　８月２７日登記┃
	pattern := fmt.Sprintf("┃　*│　*([%s]+)変更　*([%s]+)登記　*┃", ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimPattern(s, pattern), trimAllSpace(matches[1]), trimAllSpace(matches[2])
	}
	return s, "", ""
}

func trimRegisterAt(s string) (string, string) {
	// trim ┃事項　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　　平成２０年　７月２５日登記┃
	pattern := fmt.Sprintf("┃[%s]+　*│　*([%s]+)(登記|移記)　*┃", ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimPattern(s, pattern), trimAllSpace(matches[1])
	}
	return s, ""
}
