package toukibo

import "strings"

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
