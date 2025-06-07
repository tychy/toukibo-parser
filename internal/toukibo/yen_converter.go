package toukibo

// 1万円 -> 10000
// １億円 -> 100000000
// １兆円 -> 1000000000000
// 30億2000万1000円 -> 3020001000

// 上限は兆円まで
func YenToNumber(yen string) int {
	// 新しい共通関数を使用
	result, err := ParseJapaneseNumberWithSuffix(yen, JapaneseNumberUnits, "円")
	if err != nil {
		// エラーの場合は0を返す（既存の動作を維持）
		return 0
	}
	return result
}
