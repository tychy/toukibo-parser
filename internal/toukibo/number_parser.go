package toukibo

import (
	"regexp"
	"strconv"
	"strings"
)

// UnitMultiplier は日本語の数値単位と倍率を表す
type UnitMultiplier struct {
	Unit       rune
	Multiplier int
}

// JapaneseNumberUnits は共通の日本語数値単位
var JapaneseNumberUnits = []UnitMultiplier{
	{'兆', 1000000000000},
	{'億', 100000000},
	{'万', 10000},
}

// ParseJapaneseNumber は日本語の数値表記をパースする共通関数
func ParseJapaneseNumber(s string, units []UnitMultiplier) (int, error) {
	// 全角数字を半角に変換
	s = ZenkakuToHankaku(s)
	
	total := 0
	remaining := s
	
	// 各単位ごとに処理
	for _, unit := range units {
		pattern := `(\d+)` + string(unit.Unit)
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(remaining)
		
		if len(matches) > 1 {
			value, err := strconv.Atoi(matches[1])
			if err != nil {
				return 0, err
			}
			total += value * unit.Multiplier
			// 処理済み部分を削除
			remaining = strings.Replace(remaining, matches[0], "", 1)
		}
	}
	
	// 残りの数字（単位なし）を処理
	remaining = strings.TrimSpace(remaining)
	if remaining != "" {
		// 数字以外の文字を削除
		regex := regexp.MustCompile(`\d+`)
		matches := regex.FindString(remaining)
		if matches != "" {
			value, err := strconv.Atoi(matches)
			if err != nil {
				return 0, err
			}
			total += value
		}
	}
	
	return total, nil
}

// ParseJapaneseNumberWithSuffix は特定の接尾辞を持つ日本語数値をパースする
func ParseJapaneseNumberWithSuffix(s string, units []UnitMultiplier, suffix string) (int, error) {
	// 接尾辞を削除
	s = strings.TrimSuffix(s, suffix)
	s = strings.TrimSpace(s)
	
	// プレフィックス（金など）を削除
	if strings.HasPrefix(s, "金") {
		s = strings.TrimPrefix(s, "金")
	}
	
	return ParseJapaneseNumber(s, units)
}